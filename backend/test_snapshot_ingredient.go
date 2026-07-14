package main

import (
	"context"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/infra/database"
	"github.com/jeanGouveia/pratoOnline/backend/internal/infra/repository"
	"github.com/jeanGouveia/pratoOnline/backend/internal/service"
)

func main() {
	// Conecta ao banco
	db, err := gorm.Open(sqlite.Open("test_snapshot_ingredient.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Configura pragmas
	sqlDB, _ := db.DB()
	sqlDB.Exec("PRAGMA journal_mode=WAL;")
	sqlDB.Exec("PRAGMA foreign_keys=ON;")

	// Executa migrações
	if err := database.RunMigrations(db); err != nil {
		log.Fatal(err)
	}

	// Cria repositories
	productRepo := repository.NewGormProductRepository(db)
	stockAdjustmentRepo := repository.NewGormStockAdjustmentRepository(db, productRepo)
	orderRepo := repository.NewGormOrderRepository(db, productRepo, stockAdjustmentRepo)

	// Cria services
	stockAdjustmentService := service.NewStockAdjustmentService(stockAdjustmentRepo, productRepo)

	ctx := context.Background()

	// PASSO 1: Criar Ingrediente
	fmt.Println("=== PASSO 1: Criar Ingrediente ===")
	ingredient := &domain.Ingredient{
		Name:          "Queijo Mussarela",
		Unit:          "kg",
		StockQuantity: 10,
		MinStock:      2,
	}
	if err := productRepo.CreateIngredient(ctx, ingredient); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Ingrediente criado: ID=%d, Nome=%s, Unidade=%s\n", ingredient.ID, ingredient.Name, ingredient.Unit)

	// PASSO 2: Criar Produto com Ingrediente
	fmt.Println("\n=== PASSO 2: Criar Produto com Ingrediente ===")
	product := &domain.Product{
		Name:       "Pizza Queijo",
		Price:      50,
		IsComposto: true,
		Active:     true,
	}
	if err := productRepo.CreateProduct(ctx, product); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Produto criado: ID=%d, Nome=%s\n", product.ID, product.Name)

	// Adicionar ingrediente à ficha técnica
	productIngredientsInput := []domain.ProductIngredient{
		{
			ProductID:    product.ID,
			IngredientID: ingredient.ID,
			Quantity:     0.5, // 0.5 kg por pizza
		},
	}
	if err := productRepo.SetProductIngredients(ctx, product.ID, productIngredientsInput); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Ingrediente adicionado à ficha técnica: %s (%.2f %s)\n", ingredient.Name, 0.5, ingredient.Unit)

	// PASSO 3: Criar Pedido
	fmt.Println("\n=== PASSO 3: Criar Pedido ===")
	order := &domain.Order{
		Status:     domain.OrderStatusPending,
		TotalPrice: 100, // 50 * 2
		Notes:      "Pedido de teste",
		Items: []domain.OrderItem{
			{
				ProductID: product.ID,
				Quantity:  2,
				UnitPrice: 50,
			},
		},
	}

	// Pré-carrega ingredientes (com Preload para ter Ingredient populado)
	productIngredientsLoaded, err := productRepo.GetProductIngredients(ctx, product.ID)
	if err != nil {
		log.Fatal(err)
	}
	productIngredientsMap := map[uint][]domain.ProductIngredient{
		product.ID: productIngredientsLoaded,
	}

	if err := orderRepo.CreateOrder(ctx, order, productIngredientsMap); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Pedido criado: ID=%d, Total=%.2f\n", order.ID, order.TotalPrice)

	// PASSO 4: Cancelar Pedido (gera ajustes pendentes)
	fmt.Println("\n=== PASSO 4: Cancelar Pedido ===")
	// Atualizar status para cancelled
	if err := orderRepo.UpdateOrderStatus(ctx, order.ID, domain.OrderStatusCancelled); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Pedido cancelado: ID=%d\n", order.ID)

	// Registrar ajustes pendentes
	if err := stockAdjustmentService.RegisterStockAdjustmentForOrder(
		ctx,
		order.ID,
		domain.OrderStatusCancelled,
		productIngredientsMap,
		order.Items,
	); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Ajustes pendentes registrados para o pedido %d\n", order.ID)

	// PASSO 5: Consultar Ajustes Pendentes
	fmt.Println("\n=== PASSO 5: Consultar Ajustes Pendentes ===")
	adjustments, err := stockAdjustmentService.GetPendingAdjustmentsByOrder(ctx, order.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total de ajustes pendentes: %d\n", len(adjustments))
	for i, adj := range adjustments {
		fmt.Printf("  Ajuste %d:\n", i+1)
		fmt.Printf("    IngredientID: %d\n", adj.IngredientID)
		fmt.Printf("    Snapshot Nome: %s\n", adj.IngredientName)
		fmt.Printf("    Snapshot Unidade: %s\n", adj.IngredientUnit)
		fmt.Printf("    Quantidade: %.4f\n", adj.Quantity)
		fmt.Printf("    Status: %s\n", adj.Status)
	}

	// PASSO 6: Alterar Ingrediente
	fmt.Println("\n=== PASSO 6: Alterar Ingrediente ===")
	ingredient.Name = "Queijo Mussarela Premium"
	ingredient.Unit = "g" // Altera de kg para g
	if err := productRepo.UpdateIngredient(ctx, ingredient); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Ingrediente alterado: ID=%d, Nome=%s, Unidade=%s\n", ingredient.ID, ingredient.Name, ingredient.Unit)

	// PASSO 7: Consultar Ajustes Pendentes Novamente
	fmt.Println("\n=== PASSO 7: Consultar Ajustes Pendentes Novamente ===")
	adjustmentsAfter, err := stockAdjustmentService.GetPendingAdjustmentsByOrder(ctx, order.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total de ajustes pendentes: %d\n", len(adjustmentsAfter))
	for i, adj := range adjustmentsAfter {
		fmt.Printf("  Ajuste %d:\n", i+1)
		fmt.Printf("    IngredientID: %d\n", adj.IngredientID)
		fmt.Printf("    Snapshot Nome: %s\n", adj.IngredientName)
		fmt.Printf("    Snapshot Unidade: %s\n", adj.IngredientUnit)
		fmt.Printf("    Quantidade: %.4f\n", adj.Quantity)
		fmt.Printf("    Status: %s\n", adj.Status)
	}

	// VERIFICAÇÃO
	fmt.Println("\n=== VERIFICAÇÃO ===")
	if adjustmentsAfter[0].IngredientName == "Queijo Mussarela" && adjustmentsAfter[0].IngredientUnit == "kg" {
		fmt.Println("✅ SUCESSO: Snapshot histórico preservado!")
		fmt.Println("   O ajuste mostra os dados originais (Queijo Mussarela, kg)")
		fmt.Println("   mesmo após o ingrediente ter sido alterado (Queijo Mussarela Premium, g)")
	} else {
		fmt.Println("❌ FALHA: Snapshot histórico não foi preservado!")
		fmt.Printf("   Esperado: Queijo Mussarela, kg\n")
		fmt.Printf("   Recebido: %s, %s\n", adjustmentsAfter[0].IngredientName, adjustmentsAfter[0].IngredientUnit)
	}

	// Verifica o ingrediente atual
	currentIngredient, _ := productRepo.FindIngredientByID(ctx, ingredient.ID)
	fmt.Printf("\nIngrediente atual no cadastro: %s, %s\n", currentIngredient.Name, currentIngredient.Unit)
}
