package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/jeanGouveia/pratoOnline/backend/internal/handler"
	"github.com/jeanGouveia/pratoOnline/backend/internal/infra/database"
	"github.com/jeanGouveia/pratoOnline/backend/internal/infra/repository"
	"github.com/jeanGouveia/pratoOnline/backend/internal/middleware"
	"github.com/jeanGouveia/pratoOnline/backend/internal/service"
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

	// --- Injeção de Dependência (DI manual, sem framework) ---
	userRepo := repository.NewGormUserRepository(db)
	productRepo := repository.NewGormProductRepository(db)
	categoryRepo := repository.NewGormCategoryRepository(db)
	stockAdjustmentRepo := repository.NewGormStockAdjustmentRepository(db, productRepo)
	orderRepo := repository.NewGormOrderRepository(db, productRepo, stockAdjustmentRepo)
	mediaRepo := repository.NewGormMediaRepository(db)
	dashboardRepo := repository.NewGormDashboardRepository(db)
	companyRepo := repository.NewGormCompanyRepository(db)

	authSvc := service.NewAuthService(userRepo)
	productSvc := service.NewProductService(productRepo)
	categorySvc := service.NewCategoryService(categoryRepo)
	orderSvc := service.NewOrderService(orderRepo, productRepo)
	stockAdjustmentSvc := service.NewStockAdjustmentService(stockAdjustmentRepo, productRepo)
	mediaSvc := service.NewMediaService(mediaRepo)
	companySvc := service.NewCompanyService(companyRepo)
	companySettingsSvc := service.NewCompanySettingsService(companyRepo, userRepo)
	themeSvc := service.NewThemeService(companyRepo, userRepo)
	businessSvc := service.NewBusinessService(companyRepo, userRepo)
	rbacSvc := service.NewRBACService(userRepo)

	authHandler := handler.NewAuthHandler(authSvc, userRepo)
	productHandler := handler.NewProductHandler(productSvc)
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	orderHandler := handler.NewOrderHandler(orderSvc)
	stockAdjustmentHandler := handler.NewStockAdjustmentHandler(stockAdjustmentSvc)
	mediaHandler := handler.NewMediaHandler(mediaSvc)
	dashboardHandler := handler.NewDashboardHandler(dashboardRepo)
	companyHandler := handler.NewCompanyHandler(companySvc)
	companySettingsHandler := handler.NewCompanySettingsHandler(companySettingsSvc)
	themeHandler := handler.NewThemeHandler(themeSvc)
	businessHandler := handler.NewBusinessHandler(businessSvc)
	systemHandler := handler.NewSystemHandler()
	authMw := middleware.NewAuthMiddleware(authSvc)
	tenantMw := middleware.NewTenantMiddleware(userRepo)
	_ = middleware.NewRoleMiddleware(rbacSvc) // Infrastructure for future use (Sprint 6)

	// --- Router ---
	r := chi.NewRouter()

	// Middlewares globais
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))

	// --- Rotas públicas ---
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status":"ok","service":"pratoOnline"}`)
	})

	r.Route("/api/system", func(r chi.Router) {
		systemHandler.RegisterRoutes(r)
	})

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/logout", authHandler.Logout)
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
		r.Post("/api/companies", companyHandler.CreateCompany)
		r.Get("/api/companies", companyHandler.ListCompanies)
		r.Get("/api/companies/{id}", companyHandler.GetCompany)
		r.Put("/api/companies/{id}", companyHandler.UpdateCompany)
		r.Delete("/api/companies/{id}", companyHandler.DeleteCompany)

		// Company Settings (Platform 2.0 - Sprint 5)
		r.Get("/api/company/settings", companySettingsHandler.GetSettings)
		r.Put("/api/company/settings", companySettingsHandler.UpdateSettings)

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
		r.Patch("/api/orders/{id}/status", orderHandler.UpdateOrderStatus)

		// Ajustes de Estoque
		r.Get("/api/stock-adjustments/pending", stockAdjustmentHandler.ListPendingAdjustments)
		r.Post("/api/stock-adjustments/{id}/approve", stockAdjustmentHandler.ApproveAdjustment)
		r.Post("/api/stock-adjustments/{id}/reject", stockAdjustmentHandler.RejectAdjustment)

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
	log.Printf("✅ PratoOnline backend iniciado em http://localhost:%s", port)

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
