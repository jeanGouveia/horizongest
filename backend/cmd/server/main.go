package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/crypto/bcrypt"

	"github.com/jeanGouveia/horizongest/backend/internal/config"
	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/handler"
	"github.com/jeanGouveia/horizongest/backend/internal/infra/database"
	"github.com/jeanGouveia/horizongest/backend/internal/infra/health"
	"github.com/jeanGouveia/horizongest/backend/internal/infra/messaging/rabbitmq"
	"github.com/jeanGouveia/horizongest/backend/internal/infra/redis"
	"github.com/jeanGouveia/horizongest/backend/internal/infra/repository"
	"github.com/jeanGouveia/horizongest/backend/internal/infra/shutdown"
	"github.com/jeanGouveia/horizongest/backend/internal/middleware"
	"github.com/jeanGouveia/horizongest/backend/internal/observability"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
	"github.com/jeanGouveia/horizongest/backend/internal/util"
)

func main() {
	// Carrega .env se existir (ignora erro em produção)
	_ = godotenv.Load()

	// FASE A.4: Load and validate configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// FASE A: Initialize structured logger
	env := cfg.Environment
	serviceName := getEnv("SERVICE_NAME", "horizongest-backend")
	util.InitLogger(serviceName, env)

	// FASE A.4.1: Startup validation - validate critical dependencies before starting
	startupValidator := health.NewStartupValidator(cfg)
	ctx := context.Background()
	if err := startupValidator.Validate(ctx); err != nil {
		util.LogFatal("Startup validation failed - aborting", map[string]interface{}{"error": err.Error()})
	}

	// --- Banco de dados ---
	// FASE A.4.1: Use cfg.DatabaseURL instead of individual DB config
	db, err := database.Connect(database.DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		DBName:   getEnv("DB_NAME", "horizongest"),
		User:     getEnv("DB_USER", "horizongest_user"),
		Password: getEnv("DB_PASSWORD", "horizongest_secure_password"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	})
	// Note: Database connection still uses individual env vars for compatibility
	// Future: migrate to use cfg.DatabaseURL directly
	if err != nil {
		util.LogFatal("falha ao conectar banco", map[string]interface{}{"error": err.Error()})
	}

	// --- Migrações automáticas (estrutura das tabelas) ---
	if err := database.RunMigrations(db); err != nil {
		util.LogFatal("falha ao executar migrações", map[string]interface{}{"error": err.Error()})
	}

	// --- Seed (dados iniciais) ---
	if os.Getenv("RUN_SEED") == "true" {
		util.LogInfo("Executando seed", nil)
		ctx := context.Background()
		platformUserRepo := repository.NewGormPlatformUserRepository(db)

		// Verificar se já existe usuário admin
		existing, err := platformUserRepo.FindByEmail(ctx, "admin@platform.com")
		if err != nil {
			util.LogFatal("Erro ao verificar usuário admin", map[string]interface{}{"error": err.Error()})
		}
		if existing == nil {
			// Criar usuário admin padrão
			hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
			if err != nil {
				util.LogFatal("Erro ao gerar hash da senha", map[string]interface{}{"error": err.Error()})
			}

			adminUser := &domain.PlatformUser{
				Name:         "Administrador",
				Email:        "admin@platform.com",
				PasswordHash: string(hash),
				Role:         domain.PlatformRoleAdmin,
			}
			if err := platformUserRepo.Create(ctx, adminUser); err != nil {
				util.LogFatal("Erro ao criar usuário admin", map[string]interface{}{"error": err.Error()})
			}
			util.LogInfo("Usuário admin criado com sucesso", nil)
		} else {
			util.LogInfo("Usuário admin já existe", nil)
		}
	}

	// --- Injeção de Dependência (DI manual, sem framework) ---
	userRepo := repository.NewGormUserRepository(db)
	productRepo := repository.NewGormProductRepository(db)
	categoryRepo := repository.NewGormCategoryRepository(db)
	stockAdjustmentRepo := repository.NewGormStockAdjustmentRepository(db, productRepo)
	orderRepo := repository.NewGormOrderRepository(db, productRepo, stockAdjustmentRepo)
	mediaRepo := repository.NewGormMediaRepository(db)
	dashboardRepo := repository.NewGormDashboardRepository(db)
	companyRepo := repository.NewGormCompanyRepository(db)
	tokenBlacklistRepo := repository.NewGormTokenBlacklistRepository(db)
	passwordResetRepo := repository.NewGormPasswordResetRepository(db)
	stockMovementRepo := repository.NewGormStockMovementRepository(db)
	financeRepo := repository.NewGormFinanceRepository(db)
	purchaseRepo := repository.NewGormPurchaseRepository(db)
	reportRepo := repository.NewGormReportRepository(db)

	// Audit repository (FASE A.4: Audit Logs)
	auditRepo := repository.NewAuditRepository(db)

	// Outbox repository (FASE A.4: Outbox pattern)
	outboxRepo := repository.NewGormOutboxRepository(db)

	// Platform repositories (Sprint 3.2)
	platformUserRepo := repository.NewGormPlatformUserRepository(db)
	platformSessionRepo := repository.NewGormPlatformSessionRepository(db)
	platformAuditRepo := repository.NewGormPlatformAuditRepository(db)
	planRepo := repository.NewGormPlanRepository(db)
	platformBrandRepo := repository.NewGormPlatformBrandRepository(db) // Sprint 3.5
	impersonationAuditRepo := repository.NewGormImpersonationAuditRepository(db)

	// JWT secrets (Sprint 3.4 - Security Hardening)
	// FASE A: Validar secrets com validação de entropia e comprimento mínimo
	jwtPlatformSecret := getEnv("JWT_PLATFORM_SECRET", "")
	jwtTenantSecret := getEnv("JWT_TENANT_SECRET", "")
	isProduction := env == "production"

	// Validate JWT Platform Secret
	if jwtPlatformSecret == "" {
		if isProduction {
			util.LogFatal("JWT_PLATFORM_SECRET não configurado em produção", nil)
		}
		// Generate secure secret for development
		generatedSecret, err := util.GenerateSecureSecretAlphaNum(32)
		if err != nil {
			log.Fatalf("FATAL: failed to generate secure secret: %v", err)
		}
		jwtPlatformSecret = generatedSecret
		log.Printf("WARNING: JWT_PLATFORM_SECRET gerado automaticamente para desenvolvimento: %s", jwtPlatformSecret[:8]+"...")
	} else {
		if err := util.ValidateJWTSecret(jwtPlatformSecret, isProduction); err != nil {
			util.LogFatal("JWT_PLATFORM_SECRET inválido", map[string]interface{}{"error": err.Error()})
		}
	}

	// Validate JWT Tenant Secret
	if jwtTenantSecret == "" {
		if isProduction {
			util.LogFatal("JWT_TENANT_SECRET não configurado em produção", nil)
		}
		// Generate secure secret for development
		generatedSecret, err := util.GenerateSecureSecretAlphaNum(32)
		if err != nil {
			log.Fatalf("FATAL: failed to generate secure secret: %v", err)
		}
		jwtTenantSecret = generatedSecret
		log.Printf("WARNING: JWT_TENANT_SECRET gerado automaticamente para desenvolvimento: %s", jwtTenantSecret[:8]+"...")
	} else {
		if err := util.ValidateJWTSecret(jwtTenantSecret, isProduction); err != nil {
			util.LogFatal("JWT_TENANT_SECRET inválido", map[string]interface{}{"error": err.Error()})
		}
	}

	// Platform services (Sprint 3.2)
	sessionDuration := 24 * time.Hour
	platformAuthSvc := service.NewPlatformAuthService(platformUserRepo, platformSessionRepo, jwtPlatformSecret, sessionDuration, bcrypt.DefaultCost)
	platformBrandSvc := service.NewPlatformBrandService(platformBrandRepo) // Sprint 3.5

	// Tenant services
	productSvc := service.NewProductService(productRepo, db)
	categorySvc := service.NewCategoryService(categoryRepo)
	orderSvc := service.NewOrderService(orderRepo, productRepo, outboxRepo)
	stockAdjustmentSvc := service.NewStockAdjustmentService(stockAdjustmentRepo, productRepo)
	stockMovementSvc := service.NewStockMovementService(stockMovementRepo, productRepo, db)
	mediaSvc := service.NewMediaService(mediaRepo)
	companySvc := service.NewCompanyService(companyRepo)
	companySettingsSvc := service.NewCompanySettingsService(companyRepo, userRepo)
	themeSvc := service.NewThemeService(companyRepo, userRepo)
	businessSvc := service.NewBusinessService(companyRepo, userRepo)
	rbacSvc := service.NewRBACService(userRepo)
	userManagementSvc := service.NewUserManagementService(userRepo, companyRepo, rbacSvc)
	financeSvc := service.NewFinanceService(financeRepo)
	purchaseSvc := service.NewPurchaseService(purchaseRepo, productRepo, stockMovementRepo, db)
	reportSvc := service.NewReportService(reportRepo)

	// Initialize platform brand config (Sprint 3.6)
	if err := platformBrandSvc.Initialize(context.Background()); err != nil {
		log.Printf("WARNING: failed to initialize platform brand config: %v", err)
	}

	// Get platform brand for email and JWT issuer (Sprint 3.6)
	platformBrand, err := platformBrandSvc.Get(context.Background())
	if err != nil {
		log.Printf("WARNING: failed to get platform brand config: %v", err)
		// Fallback to environment variables
		platformBrand = &domain.PlatformBrandConfig{
			SupportEmail: getEnv("SUPPORT_EMAIL", "noreply@localhost"),
			PlatformName: getEnv("PLATFORM_NAME", "HorizonGest"),
		}
	}

	// Initialize auth service with platform brand as JWT issuer (Sprint 3.6)
	authSvc := service.NewAuthService(userRepo, companyRepo, tokenBlacklistRepo, passwordResetRepo, jwtTenantSecret, platformBrand.PlatformName)

	// Impersonation service
	impersonationSvc := service.NewImpersonationService(authSvc, companyRepo, userRepo, impersonationAuditRepo)

	emailSvc := service.NewEmailService(false, platformBrand.SupportEmail, platformBrand.PlatformName) // Email disabled by default
	platformSvc := service.NewPlatformService(db, companyRepo, userRepo, platformUserRepo, platformAuditRepo, emailSvc)
	planSvc := service.NewPlanService(planRepo)
	backupSvc := service.NewBackupService(
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "3306"),
		getEnv("DB_USER", "root"),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_NAME", "horizongest"), // Default DB name - should be configured via environment
		cfg.StoragePath,                  // FASE A.4.1: Use cfg.StoragePath for backup/exports directory
		platformBrand.PlatformName,       // Backup filename prefix from platform brand (Sprint 3.6)
	)
	exportSvc := service.NewExportService(companyRepo, userRepo, cfg.StoragePath) // FASE A.4.1: Use cfg.StoragePath

	// FASE A.4: Initialize metrics service for Prometheus
	_ = observability.NewMetricsService()

	// FASE A.4: Initialize tracing service for OpenTelemetry
	tracingSvc := observability.NewTracingService(serviceName)
	_ = tracingSvc // Tracing service initialized for distributed tracing

	// FASE A.4: Initialize audit service for audit logging
	auditSvc := service.NewAuditService(auditRepo)

	// --- Async Infrastructure (SPRINT 5D.1 - Critical Fix O1-O4) ---

	// Redis Client (for distributed cache and idempotency)
	// FASE A.4.1: Use cfg.RedisURL for Redis configuration
	var redisClient *redis.Client
	if cfg.RedisURL != "" && cfg.RedisURL != "redis://localhost:6379" {
		redisConfig := redis.Config{
			Host:         getEnv("REDIS_HOST", "localhost"),
			Port:         getEnvAsInt("REDIS_PORT", 6379),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvAsInt("REDIS_DB", 0),
			PoolSize:     cfg.RedisPoolSize,
			MinIdleConns: cfg.RedisMinIdleConns,
			MaxIdleConns: 20,
			MaxRetries:   3,
			DialTimeout:  cfg.RedisDialTimeout,
			ReadTimeout:  cfg.RedisReadTimeout,
			WriteTimeout: cfg.RedisWriteTimeout,
			PoolTimeout:  cfg.RedisPoolTimeout,
		}
		var err error
		redisClient, err = redis.NewClient(redisConfig)
		if err != nil {
			log.Printf("WARNING: failed to initialize Redis client: %v (async features disabled)", err)
		} else {
			log.Println("✅ Redis client initialized")
			defer redisClient.Close()
		}
	} else {
		log.Println("WARNING: REDIS_URL not configured, async features disabled")
	}

	// RabbitMQ Publisher (for event publishing)
	// FASE A.4.1: Use cfg.RabbitMQURL for RabbitMQ configuration
	var rabbitmqPublisher *rabbitmq.RabbitMQPublisher
	if cfg.RabbitMQURL != "" && cfg.RabbitMQURL != "amqp://localhost:5672" {
		rabbitmqConfig := rabbitmq.Config{
			URL:                    cfg.RabbitMQURL,
			Exchange:               cfg.RabbitMQExchange,
			ExchangeType:           "topic",
			QueuePrefix:            "horizongest",
			RetryCount:             3,
			PublisherTimeout:       10 * time.Second,
			ReconnectDelay:         cfg.RabbitMQReconnectInterval,
			EnablePublisherConfirm: true,
		}
		var err error
		rabbitmqPublisher, err = rabbitmq.NewRabbitMQPublisher(rabbitmqConfig)
		if err != nil {
			log.Printf("WARNING: failed to initialize RabbitMQ publisher: %v (event publishing disabled)", err)
		} else {
			log.Println("✅ RabbitMQ publisher initialized")
			defer rabbitmqPublisher.Close()
		}
	} else {
		log.Println("WARNING: RABBITMQ_URL not configured, event publishing disabled")
	}

	// EventDispatcher (for outbox pattern)
	var eventDispatcher *service.EventDispatcher
	if rabbitmqPublisher != nil {
		outboxRepo := repository.NewGormOutboxRepository(db)
		dispatcherConfig := service.DispatcherConfig{
			BatchSize:  100,
			Interval:   5 * time.Second,
			RetryCount: 3,
		}
		eventDispatcher = service.NewEventDispatcher(outboxRepo, rabbitmqPublisher, dispatcherConfig)
		ctx := context.Background()
		eventDispatcher.Start(ctx)
		log.Println("✅ EventDispatcher started")
		defer eventDispatcher.Shutdown()
	}

	// Consumers (for email, webhook, etc.)
	// TODO: Initialize consumers when RabbitMQ is available
	// For now, consumers are not started to avoid blocking server startup

	// Middlewares
	authMw := middleware.NewAuthMiddleware(authSvc)
	tenantMw := middleware.NewTenantMiddleware(userRepo)
	roleMw := middleware.NewRoleMiddleware(rbacSvc)                         // Infrastructure for Sprint 7
	platformAuthMw := middleware.NewPlatformAuthMiddleware(platformAuthSvc) // Sprint 3.2
	rateLimiter := middleware.NewRateLimiter(5, 30)                         // 5 req/min per IP, 30 req/hour per user (Sprint 3.4)

	authHandler := handler.NewAuthHandler(authSvc, userRepo, auditSvc)
	productHandler := handler.NewProductHandler(productSvc)
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	orderHandler := handler.NewOrderHandler(orderSvc)
	stockAdjustmentHandler := handler.NewStockAdjustmentHandler(stockAdjustmentSvc)
	stockMovementHandler := handler.NewStockMovementHandler(stockMovementSvc, roleMw)
	mediaHandler := handler.NewMediaHandler(mediaSvc)
	dashboardHandler := handler.NewDashboardHandler(dashboardRepo)
	companyHandler := handler.NewCompanyHandler(companySvc)
	companySettingsHandler := handler.NewCompanySettingsHandler(companySettingsSvc)
	themeHandler := handler.NewThemeHandler(themeSvc)
	businessHandler := handler.NewBusinessHandler(businessSvc)
	userManagementHandler := handler.NewUserManagementHandler(userManagementSvc, userRepo)
	systemHandler := handler.NewSystemHandler()

	// Platform handlers (Sprint 3.2)
	platformAuthHandler := handler.NewPlatformAuthHandler(platformAuthSvc)
	platformCompanyHandler := handler.NewPlatformCompanyHandler(platformSvc)
	platformDashboardHandler := handler.NewPlatformDashboardHandler(platformSvc)
	planHandler := handler.NewPlanHandler(planSvc)
	backupHandler := handler.NewBackupHandler(backupSvc)
	exportHandler := handler.NewExportHandler(exportSvc)
	platformBrandHandler := handler.NewPlatformBrandHandler(platformBrandSvc) // Sprint 3.5
	impersonationHandler := handler.NewImpersonationHandler(impersonationSvc)
	forensicHandler := handler.NewForensicHandler() // Forensic investigation
	financeHandler := handler.NewFinanceHandler(financeSvc)
	purchaseHandler := handler.NewPurchaseHandler(purchaseSvc)
	reportHandler := handler.NewReportHandler(reportSvc)

	// --- Router ---
	r := chi.NewRouter()

	// Middlewares globais
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(middleware.CORS)                  // CORS middleware
	r.Use(middleware.SecurityHeaders)       // Sprint 3.4 - Security headers
	r.Use(middleware.CorrelationMiddleware) // FASE A: CorrelationID and RequestID propagation
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))

	// --- Rotas públicas ---
	// FASE A: Health check endpoints
	healthHandler := handler.NewHealthHandler(cfg)
	r.Route("/health", healthHandler.RegisterRoutes)

	// FASE A: Metrics endpoint for Prometheus
	r.Handle("/metrics", promhttp.Handler())

	// Public platform branding endpoint (Sprint 3.6 - White Label)
	r.Get("/api/public/brand", platformBrandHandler.GetPublicPlatformBrand)

	// Public business types endpoint
	r.Get("/api/public/business-types", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		types := domain.AllBusinessTypes()
		type BusinessTypeResponse struct {
			Value string `json:"value"`
			Label string `json:"label"`
		}
		response := make([]BusinessTypeResponse, len(types))
		for i, bt := range types {
			response[i] = BusinessTypeResponse{
				Value: bt.String(),
				Label: bt.DisplayName(),
			}
		}
		json.NewEncoder(w).Encode(response)
	})

	r.Route("/api/system", func(r chi.Router) {
		systemHandler.RegisterRoutes(r)
	})

	// Platform routes (Sprint 3.2)
	r.Route("/api/platform/auth", func(r chi.Router) {
		r.Use(rateLimiter.RateLimitByIP) // Sprint 3.4 - Rate limiting
		r.Post("/login", platformAuthHandler.Login)
		r.Post("/logout", platformAuthHandler.Logout)
		r.Get("/me", platformAuthHandler.Me)
	})

	r.Route("/api/platform/dashboard", func(r chi.Router) {
		r.Use(platformAuthMw.Auth)
		r.Use(platformAuthMw.RequireAdmin)
		r.Get("/stats", platformDashboardHandler.GetStats)
	})

	r.Route("/api/platform/companies", func(r chi.Router) {
		r.Use(platformAuthMw.Auth)
		r.Use(platformAuthMw.RequireAdmin)
		r.Post("/", platformCompanyHandler.CreateCompany)
		r.Get("/", platformCompanyHandler.ListCompanies)
		r.Get("/{id}", platformCompanyHandler.GetCompany)
		r.Put("/{id}", platformCompanyHandler.UpdateCompany)
		r.Post("/{id}/deactivate", platformCompanyHandler.DeactivateCompany)
		r.Post("/{id}/activate", platformCompanyHandler.ActivateCompany)
		r.Get("/{id}/owner", platformCompanyHandler.GetCompanyOwner)
		r.Post("/{id}/owner/reset-password", platformCompanyHandler.ResetOwnerPassword)
		r.Post("/{id}/login-as", platformCompanyHandler.LoginAsCompany)
		r.Post("/{id}/trial", platformCompanyHandler.SetCompanyTrial)
		r.Post("/{id}/suspend", platformCompanyHandler.SuspendCompany)
		r.Post("/{id}/cancel", platformCompanyHandler.CancelCompany)
		r.Post("/{id}/reactivate", platformCompanyHandler.ReactivateCompany)
	})

	r.Route("/api/platform/users", func(r chi.Router) {
		r.Use(platformAuthMw.Auth)
		r.Use(platformAuthMw.RequireAdmin)
		r.Post("/{id}/block", platformCompanyHandler.BlockUser)
		r.Post("/{id}/unblock", platformCompanyHandler.UnblockUser)
	})

	r.Route("/api/platform/plans", func(r chi.Router) {
		r.Use(platformAuthMw.Auth)
		r.Use(platformAuthMw.RequireAdmin)
		r.Post("/", planHandler.CreatePlan)
		r.Get("/", planHandler.ListPlans)
		r.Get("/active", planHandler.ListActivePlans)
		r.Get("/{id}", planHandler.GetPlan)
		r.Put("/{id}", planHandler.UpdatePlan)
		r.Delete("/{id}", planHandler.DeletePlan)
	})

	r.Route("/api/platform/backup", func(r chi.Router) {
		r.Use(platformAuthMw.Auth)
		r.Use(platformAuthMw.RequireAdmin)
		r.Post("/", backupHandler.CreateBackup)
		r.Get("/", backupHandler.ListBackups)
		r.Delete("/", backupHandler.DeleteBackup)
	})

	r.Route("/api/platform/export", func(r chi.Router) {
		r.Use(platformAuthMw.Auth)
		r.Use(platformAuthMw.RequireAdmin)
		r.Post("/companies", exportHandler.ExportCompanies)
		r.Post("/users", exportHandler.ExportUsers)
	})

	r.Route("/api/platform/brand", func(r chi.Router) {
		r.Use(platformAuthMw.Auth)
		r.Use(platformAuthMw.RequireAdmin)
		r.Get("/", platformBrandHandler.GetPlatformBrand)
		r.Put("/", platformBrandHandler.UpdatePlatformBrand)
	}) // Sprint 3.5

	r.Route("/api/platform/impersonation", func(r chi.Router) {
		r.Use(platformAuthMw.Auth)
		r.Post("/start", impersonationHandler.StartImpersonation)
		r.Post("/end", impersonationHandler.EndImpersonation)
		r.Get("/active", impersonationHandler.GetActiveImpersonation)
		r.Get("/history", impersonationHandler.GetImpersonationHistory)
	})

	// Forensic investigation endpoints (only enabled in development/debug mode)
	if os.Getenv("ENABLE_FORENSIC") == "true" {
		log.Println("WARNING: Forensic endpoints enabled - DO NOT use in production")
		r.Route("/api/forensic", func(r chi.Router) {
			forensicHandler.RegisterRoutes(r)
		})
	}

	r.Route("/api/auth", func(r chi.Router) {
		r.Use(rateLimiter.RateLimitByIP) // Sprint 3.4 - Rate limiting
		r.Post("/login", authHandler.Login)
		r.Post("/logout", authHandler.Logout)
		r.Post("/request-password-reset", authHandler.RequestPasswordReset)
		r.Post("/reset-password", authHandler.ResetPassword)
	})

	// --- Rotas privadas (protegidas pelo AuthMiddleware) ---
	r.Group(func(r chi.Router) {
		r.Use(authMw.Auth)
		r.Use(tenantMw.Tenant)

		// Dashboard e perfil (todos autenticados)
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager, domain.RoleEmployee))
			r.Get("/api/dashboard", dashboardHandler.GetDashboard)
			r.Get("/api/me", authHandler.Me)
			r.Put("/api/me", authHandler.UpdateProfile)
			r.Post("/api/me/change-password", authHandler.ChangePassword)
			r.Get("/api/me/company", companyHandler.GetCurrentCompany)
		})

		// Empresas (Tenant Engine - Platform 2.0)
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin))
			r.Get("/api/companies", companyHandler.ListCompanies)
			r.Get("/api/companies/{id}", companyHandler.GetCompany)
		})
		r.Group(func(r chi.Router) {
			r.Use(roleMw.Require(domain.RoleOwner))
			r.Put("/api/companies/{id}", companyHandler.UpdateCompany)
			r.Delete("/api/companies/{id}", companyHandler.DeleteCompany)
		})

		// Company Settings (Platform 2.0 - Sprint 5)
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin))
			r.Get("/api/company/settings", companySettingsHandler.GetSettings)
			r.Put("/api/company/settings", companySettingsHandler.UpdateSettings)
		})

		// User Management (Platform 2.0 - Sprint 7)
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin))
			r.Get("/api/company/users", userManagementHandler.ListUsers)
			r.Get("/api/company/users/{id}", userManagementHandler.GetUser)
			r.Post("/api/company/users/add", userManagementHandler.AddUser)
			r.Put("/api/company/users/{id}/active", userManagementHandler.SetUserActive)
		})
		r.Group(func(r chi.Router) {
			r.Use(roleMw.Require(domain.RoleOwner))
			r.Put("/api/company/users/{id}/role", userManagementHandler.ChangeRole)
			r.Delete("/api/company/users/{id}", userManagementHandler.RemoveUser)
		})

		// Tema (White Label - Platform 2.0)
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager, domain.RoleEmployee))
			r.Get("/api/theme", themeHandler.GetTheme)
			r.Get("/api/theme/default", themeHandler.GetDefaultTheme)
		})

		// Business Profile (Business Engine - Platform 2.0)
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin))
			r.Get("/api/business/profile", businessHandler.GetBusinessProfile)
		})
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager, domain.RoleEmployee))
			r.Get("/api/business/profile/default", businessHandler.GetDefaultBusinessProfile)
		})

		// Produtos
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager, domain.RoleEmployee))
			r.Get("/api/products", productHandler.ListProducts)
			r.Get("/api/products/active", productHandler.ListActiveProducts)
			r.Get("/api/products/{id}", productHandler.GetProduct)
			r.Get("/api/products/{id}/ingredients", productHandler.GetProductIngredients)
		})
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager))
			r.Post("/api/products", productHandler.CreateProduct)
			r.Put("/api/products/{id}", productHandler.UpdateProduct)
			r.Delete("/api/products/{id}", productHandler.DeleteProduct)
			r.Post("/api/products/{id}/duplicate", productHandler.DuplicateProduct)
			r.Post("/api/products/{id}/archive", productHandler.ArchiveProduct)
			r.Put("/api/products/{id}/ingredients", productHandler.SetProductIngredients)
		})

		// Ingredientes
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager, domain.RoleEmployee))
			r.Get("/api/ingredients", productHandler.ListIngredients)
			r.Get("/api/ingredients/{id}", productHandler.GetIngredient)
		})
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager))
			r.Post("/api/ingredients", productHandler.CreateIngredient)
			r.Put("/api/ingredients/{id}", productHandler.UpdateIngredient)
			r.Delete("/api/ingredients/{id}", productHandler.DeleteIngredient)
			r.Patch("/api/ingredients/{id}/stock", productHandler.UpdateIngredientStock)
		})

		// Categorias
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager, domain.RoleEmployee))
			r.Get("/api/categories", categoryHandler.ListCategories)
			r.Get("/api/categories/{id}", categoryHandler.GetCategory)
		})
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager))
			r.Post("/api/categories", categoryHandler.CreateCategory)
			r.Put("/api/categories/{id}", categoryHandler.UpdateCategory)
			r.Delete("/api/categories/{id}", categoryHandler.DeleteCategory)
		})

		// Pedidos
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager, domain.RoleEmployee))
			r.Get("/api/orders", orderHandler.ListOrders)
			r.Get("/api/orders/{id}", orderHandler.GetOrder)
		})
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager, domain.RoleEmployee))
			r.Post("/api/orders", orderHandler.CreateOrder)
			r.Patch("/api/orders/{id}/status", orderHandler.UpdateOrderStatus)
		})
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager))
			r.Put("/api/orders/{id}", orderHandler.UpdateOrder)
		})

		// Ajustes de Estoque
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager))
			r.Get("/api/stock-adjustments/pending", stockAdjustmentHandler.ListPendingAdjustments)
			r.Post("/api/stock-adjustments/{id}/approve", stockAdjustmentHandler.ApproveAdjustment)
			r.Post("/api/stock-adjustments/{id}/reject", stockAdjustmentHandler.RejectAdjustment)
		})

		// Stock Movements (Sprint 4)
		stockMovementHandler.RegisterRoutes(r)

		// Finance (Sprint 5B.1)
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager))
			financeHandler.RegisterRoutes(r)
		})

		// Purchase (Sprint 5B.1)
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager))
			purchaseHandler.RegisterRoutes(r)
		})

		// Reports (Sprint 5B.1)
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager))
			reportHandler.RegisterRoutes(r)
		})

		// Mídia
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager, domain.RoleEmployee))
			r.Get("/api/media/{id}", mediaHandler.GetMedia)
			r.Get("/api/media/entity/{entity_type}/{entity_id}", mediaHandler.GetMediaByEntity)
		})
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin, domain.RoleManager))
			r.Post("/api/media/upload", mediaHandler.UploadMedia)
			r.Delete("/api/media/{id}", mediaHandler.DeleteMedia)
		})
	})

	// --- Rotas públicas para servir arquivos estáticos ---
	r.Get("/uploads/*", mediaHandler.ServeFile)

	// --- Servidor ---
	port := cfg.ServerPort
	env = cfg.Environment

	// FASE A.4.1: Initialize shutdown manager
	shutdownManager := shutdown.NewShutdownManager(cfg.ShutdownTimeout)

	// FASE A.4.1: Register shutdown functions
	// Register HTTP server shutdown
	shutdownManager.Register("http-server", func(ctx context.Context) error {
		util.LogInfo("Shutting down HTTP server", nil)
		return nil
	})

	// Register RabbitMQ shutdown
	if rabbitmqPublisher != nil {
		shutdownManager.Register("rabbitmq", func(ctx context.Context) error {
			util.LogInfo("Shutting down RabbitMQ publisher", nil)
			rabbitmqPublisher.Close()
			return nil
		})
	}

	// Register Redis shutdown
	if redisClient != nil {
		shutdownManager.Register("redis", func(ctx context.Context) error {
			util.LogInfo("Shutting down Redis client", nil)
			redisClient.Close()
			return nil
		})
	}

	// Register Database shutdown
	shutdownManager.Register("database", func(ctx context.Context) error {
		util.LogInfo("Shutting down database connection", nil)
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	})

	// Register EventDispatcher shutdown
	if eventDispatcher != nil {
		shutdownManager.Register("event-dispatcher", func(ctx context.Context) error {
			util.LogInfo("Shutting down event dispatcher", nil)
			eventDispatcher.Shutdown()
			return nil
		})
	}

	// FASE A.4.1: TLS/HTTPS enforcement in production
	if env == "production" {
		// In production, require HTTPS
		tlsCert := getEnv("TLS_CERT_PATH", "")
		tlsKey := getEnv("TLS_KEY_PATH", "")

		if tlsCert == "" || tlsKey == "" {
			util.LogFatal("TLS_CERT_PATH and TLS_KEY_PATH environment variables are required in production", nil)
		}

		// Configure TLS with secure settings
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			},
			PreferServerCipherSuites: true,
		}

		server := &http.Server{
			Addr:      ":" + port,
			Handler:   r,
			TLSConfig: tlsConfig,
		}

		// Start HTTPS server in goroutine
		serverErr := make(chan error, 1)
		go func() {
			util.LogInfo("Backend iniciado", map[string]interface{}{
				"platform": platformBrand.PlatformName,
				"port":     port,
				"protocol": "HTTPS",
			})
			serverErr <- server.ListenAndServeTLS(tlsCert, tlsKey)
		}()

		// FASE A.4.1: Graceful shutdown with signal handling
		sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		select {
		case err := <-serverErr:
			util.LogFatal("servidor HTTPS encerrado", map[string]interface{}{"error": err.Error()})
		case <-sigCtx.Done():
			util.LogInfo("Shutdown signal received", nil)
			// Shutdown HTTP server
			shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
			defer cancel()
			if err := server.Shutdown(shutdownCtx); err != nil {
				util.LogError("HTTP server shutdown error", map[string]interface{}{"error": err.Error()})
			}
			// Execute all registered shutdown functions
			if err := shutdownManager.Shutdown(); err != nil {
				util.LogError("Shutdown error", map[string]interface{}{"error": err.Error()})
			}
			util.LogInfo("Graceful shutdown completed", nil)
		}
	} else {
		// In development, use HTTP with warning
		util.LogWarn("Backend iniciado em modo HTTP (development)", map[string]interface{}{
			"platform": platformBrand.PlatformName,
			"port":     port,
			"warning":  "HTTP is insecure. Use HTTPS in production.",
		})

		server := &http.Server{
			Addr:    ":" + port,
			Handler: r,
		}

		// Start HTTP server in goroutine
		serverErr := make(chan error, 1)
		go func() {
			serverErr <- server.ListenAndServe()
		}()

		// FASE A.4.1: Graceful shutdown with signal handling
		sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		select {
		case err := <-serverErr:
			util.LogFatal("servidor encerrado", map[string]interface{}{"error": err.Error()})
		case <-sigCtx.Done():
			util.LogInfo("Shutdown signal received", nil)
			// Shutdown HTTP server
			shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
			defer cancel()
			if err := server.Shutdown(shutdownCtx); err != nil {
				util.LogError("HTTP server shutdown error", map[string]interface{}{"error": err.Error()})
			}
			// Execute all registered shutdown functions
			if err := shutdownManager.Shutdown(); err != nil {
				util.LogError("Shutdown error", map[string]interface{}{"error": err.Error()})
			}
			util.LogInfo("Graceful shutdown completed", nil)
		}
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
