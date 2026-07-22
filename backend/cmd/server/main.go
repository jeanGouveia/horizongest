package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/handler"
	"github.com/jeanGouveia/horizongest/backend/internal/infra/database"
	"github.com/jeanGouveia/horizongest/backend/internal/infra/repository"
	"github.com/jeanGouveia/horizongest/backend/internal/middleware"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

func main() {
	// Carrega .env se existir (ignora erro em produção)
	_ = godotenv.Load()

	// --- Banco de dados ---
	db, err := database.Connect(database.DBConfig{DSN: getEnv("DB_DSN", "app.db")})
	if err != nil {
		log.Fatalf("FATAL: falha ao conectar banco: %v", err)
	}

	// --- Migrações automáticas (estrutura das tabelas) ---
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("FATAL: falha ao executar migrações: %v", err)
	}

	// --- Seed (dados iniciais) ---
	if os.Getenv("RUN_SEED") == "true" {
		log.Println("Executando seed...")
		ctx := context.Background()
		platformUserRepo := repository.NewGormPlatformUserRepository(db)

		// Verificar se já existe usuário admin
		existing, err := platformUserRepo.FindByEmail(ctx, "admin@platform.com")
		if err != nil {
			log.Fatalf("Erro ao verificar usuário admin: %v", err)
		}
		if existing == nil {
			// Criar usuário admin padrão
			hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
			if err != nil {
				log.Fatalf("Erro ao gerar hash da senha: %v", err)
			}

			adminUser := &domain.PlatformUser{
				Name:         "Administrador",
				Email:        "admin@platform.com",
				PasswordHash: string(hash),
				Role:         domain.PlatformRoleAdmin,
				Active:       true,
			}

			if err := platformUserRepo.Create(ctx, adminUser); err != nil {
				log.Fatalf("Erro ao criar usuário admin: %v", err)
			}
			log.Println("Usuário admin criado com sucesso: admin@platform.com / admin123")
		} else {
			log.Println("Usuário admin já existe, pulando seed")
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

	// Platform repositories (Sprint 3.2)
	platformUserRepo := repository.NewGormPlatformUserRepository(db)
	platformSessionRepo := repository.NewGormPlatformSessionRepository(db)
	platformAuditRepo := repository.NewGormPlatformAuditRepository(db)
	planRepo := repository.NewGormPlanRepository(db)
	platformBrandRepo := repository.NewGormPlatformBrandRepository(db) // Sprint 3.5

	// JWT secrets (Sprint 3.4 - Security Hardening)
	jwtPlatformSecret := getEnv("JWT_PLATFORM_SECRET", "your-platform-secret-key-change-in-production")
	jwtTenantSecret := getEnv("JWT_TENANT_SECRET", "your-tenant-secret-key-change-in-production")

	// Platform services (Sprint 3.2)
	sessionDuration := 24 * time.Hour
	platformAuthSvc := service.NewPlatformAuthService(platformUserRepo, platformSessionRepo, jwtPlatformSecret, sessionDuration, bcrypt.DefaultCost)
	platformBrandSvc := service.NewPlatformBrandService(platformBrandRepo) // Sprint 3.5

	// Tenant services
	productSvc := service.NewProductService(productRepo)
	categorySvc := service.NewCategoryService(categoryRepo)
	orderSvc := service.NewOrderService(orderRepo, productRepo)
	stockAdjustmentSvc := service.NewStockAdjustmentService(stockAdjustmentRepo, productRepo)
	stockMovementSvc := service.NewStockMovementService(stockMovementRepo, productRepo)
	mediaSvc := service.NewMediaService(mediaRepo)
	companySvc := service.NewCompanyService(companyRepo)
	companySettingsSvc := service.NewCompanySettingsService(companyRepo, userRepo)
	themeSvc := service.NewThemeService(companyRepo, userRepo)
	businessSvc := service.NewBusinessService(companyRepo, userRepo)
	rbacSvc := service.NewRBACService(userRepo)
	userManagementSvc := service.NewUserManagementService(userRepo, companyRepo, rbacSvc)

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

	emailSvc := service.NewEmailService(false, platformBrand.SupportEmail, platformBrand.PlatformName) // Email disabled by default
	platformSvc := service.NewPlatformService(companyRepo, userRepo, platformUserRepo, platformAuditRepo, emailSvc)
	planSvc := service.NewPlanService(planRepo)
	backupSvc := service.NewBackupService(
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "3306"),
		getEnv("DB_USER", "root"),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_NAME", "horizongest"), // Default DB name - should be configured via environment
		getEnv("BACKUP_DIR", "./backups"),
		platformBrand.PlatformName, // Backup filename prefix from platform brand (Sprint 3.6)
	)
	exportSvc := service.NewExportService(companyRepo, userRepo, getEnv("EXPORT_DIR", "./exports"))

	authHandler := handler.NewAuthHandler(authSvc, userRepo)
	productHandler := handler.NewProductHandler(productSvc)
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	orderHandler := handler.NewOrderHandler(orderSvc)
	stockAdjustmentHandler := handler.NewStockAdjustmentHandler(stockAdjustmentSvc)
	stockMovementHandler := handler.NewStockMovementHandler(stockMovementSvc)
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

	authMw := middleware.NewAuthMiddleware(authSvc)
	tenantMw := middleware.NewTenantMiddleware(userRepo)
	roleMw := middleware.NewRoleMiddleware(rbacSvc)                         // Infrastructure for Sprint 7
	platformAuthMw := middleware.NewPlatformAuthMiddleware(platformAuthSvc) // Sprint 3.2
	rateLimiter := middleware.NewRateLimiter(5, 30)                         // 5 req/min per IP, 30 req/hour per user (Sprint 3.4)

	// --- Router ---
	r := chi.NewRouter()

	// Middlewares globais
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(middleware.CORS)            // CORS middleware
	r.Use(middleware.SecurityHeaders) // Sprint 3.4 - Security headers
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))

	// --- Rotas públicas ---
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"ok","service":"horizongest"}`)
	})

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

		r.Get("/api/dashboard", dashboardHandler.GetDashboard)

		r.Get("/api/me", authHandler.Me)
		r.Put("/api/me", authHandler.UpdateProfile)
		r.Post("/api/me/change-password", authHandler.ChangePassword)

		// Empresas (Tenant Engine - Platform 2.0)
		// REMOVED: POST /api/companies - Company creation is platform-only (Sprint 3.2)
		r.Get("/api/companies", companyHandler.ListCompanies)
		r.Get("/api/companies/{id}", companyHandler.GetCompany)
		r.Put("/api/companies/{id}", companyHandler.UpdateCompany)
		r.Delete("/api/companies/{id}", companyHandler.DeleteCompany)

		// Company Settings (Platform 2.0 - Sprint 5)
		r.Get("/api/company/settings", companySettingsHandler.GetSettings)
		r.Put("/api/company/settings", companySettingsHandler.UpdateSettings)

		// User Management (Platform 2.0 - Sprint 7)
		// Apply RBAC: Only Owner and Admin can manage users
		r.Group(func(r chi.Router) {
			r.Use(roleMw.RequireAny(domain.RoleOwner, domain.RoleAdmin))

			r.Get("/api/company/users", userManagementHandler.ListUsers)
			r.Get("/api/company/users/{id}", userManagementHandler.GetUser)
			r.Post("/api/company/users/add", userManagementHandler.AddUser)
			r.Put("/api/company/users/{id}/role", userManagementHandler.ChangeRole)
			r.Put("/api/company/users/{id}/active", userManagementHandler.SetUserActive)
			r.Delete("/api/company/users/{id}", userManagementHandler.RemoveUser)
		})

		// Tema (White Label - Platform 2.0)
		r.Get("/api/theme", themeHandler.GetTheme)
		r.Get("/api/theme/default", themeHandler.GetDefaultTheme)

		// Business Profile (Business Engine - Platform 2.0)
		r.Get("/api/business/profile", businessHandler.GetBusinessProfile)
		r.Get("/api/business/profile/default", businessHandler.GetDefaultBusinessProfile)

		// Produtos
		r.Post("/api/products", productHandler.CreateProduct)
		r.Get("/api/products", productHandler.ListProducts)
		r.Get("/api/products/active", productHandler.ListActiveProducts)
		r.Get("/api/products/{id}", productHandler.GetProduct)
		r.Put("/api/products/{id}", productHandler.UpdateProduct)
		r.Delete("/api/products/{id}", productHandler.DeleteProduct)
		r.Post("/api/products/{id}/duplicate", productHandler.DuplicateProduct)
		r.Post("/api/products/{id}/archive", productHandler.ArchiveProduct)
		r.Put("/api/products/{id}/ingredients", productHandler.SetProductIngredients)
		r.Get("/api/products/{id}/ingredients", productHandler.GetProductIngredients)

		// Ingredientes
		r.Post("/api/ingredients", productHandler.CreateIngredient)
		r.Get("/api/ingredients", productHandler.ListIngredients)
		r.Get("/api/ingredients/{id}", productHandler.GetIngredient)
		r.Put("/api/ingredients/{id}", productHandler.UpdateIngredient)
		r.Delete("/api/ingredients/{id}", productHandler.DeleteIngredient)
		r.Patch("/api/ingredients/{id}/stock", productHandler.UpdateIngredientStock)

		// Categorias
		r.Post("/api/categories", categoryHandler.CreateCategory)
		r.Get("/api/categories", categoryHandler.ListCategories)
		r.Get("/api/categories/{id}", categoryHandler.GetCategory)
		r.Put("/api/categories/{id}", categoryHandler.UpdateCategory)
		r.Delete("/api/categories/{id}", categoryHandler.DeleteCategory)

		// Pedidos
		r.Post("/api/orders", orderHandler.CreateOrder)
		r.Get("/api/orders", orderHandler.ListOrders)
		r.Get("/api/orders/{id}", orderHandler.GetOrder)
		r.Put("/api/orders/{id}", orderHandler.UpdateOrder)
		r.Patch("/api/orders/{id}/status", orderHandler.UpdateOrderStatus)

		// Ajustes de Estoque
		r.Get("/api/stock-adjustments/pending", stockAdjustmentHandler.ListPendingAdjustments)
		r.Post("/api/stock-adjustments/{id}/approve", stockAdjustmentHandler.ApproveAdjustment)
		r.Post("/api/stock-adjustments/{id}/reject", stockAdjustmentHandler.RejectAdjustment)

		// Stock Movements (Sprint 4)
		stockMovementHandler.RegisterRoutes(r)

		// Mídia
		r.Post("/api/media/upload", mediaHandler.UploadMedia)
		r.Get("/api/media/{id}", mediaHandler.GetMedia)
		r.Delete("/api/media/{id}", mediaHandler.DeleteMedia)
		r.Get("/api/media/entity/{entity_type}/{entity_id}", mediaHandler.GetMediaByEntity)
	})

	// --- Rotas públicas para servir arquivos estáticos ---
	r.Get("/uploads/*", mediaHandler.ServeFile)

	// --- Servidor ---
	port := getEnv("PORT", "8080")
	log.Printf("✅ %s backend iniciado em http://localhost:%s", platformBrand.PlatformName, port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("FATAL: servidor encerrado: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
