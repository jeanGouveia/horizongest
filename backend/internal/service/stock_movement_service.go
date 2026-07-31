package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
	"gorm.io/gorm"
)

var (
	ErrStockMovementInvalidQuantity = errors.New("quantidade inválida")
	ErrStockMovementInvalidType     = errors.New("tipo de movimentação inválido")
	ErrStockInventoryNotFound       = errors.New("inventário não encontrado")
	ErrStockInventoryCompleted      = errors.New("inventário já concluído")
)

type StockMovementService struct {
	stockMovementRepo ports.StockMovementRepository
	productRepo       ports.ProductRepository
	db                *gorm.DB
}

func NewStockMovementService(
	stockMovementRepo ports.StockMovementRepository,
	productRepo ports.ProductRepository,
	db *gorm.DB,
) *StockMovementService {
	return &StockMovementService{
		stockMovementRepo: stockMovementRepo,
		productRepo:       productRepo,
		db:                db,
	}
}

// --- Movimentações ---

// CreateStockMovement cria uma movimentação de estoque em transação atômica
// Sprint 4B.1 v2: Corrigido com transação real e SELECT FOR UPDATE
func (s *StockMovementService) CreateStockMovement(ctx context.Context, companyID, userID uint, input CreateStockMovementInput) (*domain.StockMovement, error) {
	// Validar quantidade
	if input.Quantity == 0 {
		return nil, ErrStockMovementInvalidQuantity
	}

	// Validar tipo
	if input.Type != domain.StockMovementEntry &&
		input.Type != domain.StockMovementExit &&
		input.Type != domain.StockMovementAdjust &&
		input.Type != domain.StockMovementInventory {
		return nil, ErrStockMovementInvalidType
	}

	var movement *domain.StockMovement

	// Executar toda a operação em transação atômica
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Buscar ingrediente com SELECT FOR UPDATE (lock pessimista real)
		// Sprint 4B.1 v2: Usa FindIngredientByIDForUpdate com Clauses(clause.Locking{Strength: "UPDATE"})
		ingredient, err := s.productRepo.FindIngredientByIDForUpdate(ctx, input.IngredientID, tx)
		if err != nil {
			return fmt.Errorf("StockMovementService.CreateMovement: buscar ingrediente: %w", err)
		}
		if ingredient == nil {
			return errors.New("ingrediente não encontrado")
		}

		// Validar company
		if ingredient.CompanyID != companyID {
			return errors.New("ingrediente não pertence a esta empresa")
		}

		// 2. Calcular novo estoque
		previousStock := ingredient.StockQuantity
		var newStock float64
		var quantity float64

		switch input.Type {
		case domain.StockMovementEntry, domain.StockMovementInventory:
			quantity = input.Quantity
			newStock = previousStock + quantity
		case domain.StockMovementExit, domain.StockMovementAdjust:
			quantity = -input.Quantity
			newStock = previousStock - quantity
			if newStock < 0 {
				return errors.New("estoque não pode ser negativo")
			}
		}

		// 3. Criar movimentação DENTRO da transação (passando tx)
		movement = &domain.StockMovement{
			CompanyID:     companyID,
			IngredientID:  input.IngredientID,
			Type:          input.Type,
			Quantity:      quantity,
			PreviousStock: previousStock,
			NewStock:      newStock,
			Reason:        input.Reason,
			ReferenceType: input.ReferenceType,
			ReferenceID:   input.ReferenceID,
			PerformedBy:   userID,
			PerformedAt:   time.Now(),
		}

		if err := s.stockMovementRepo.Create(ctx, movement, tx); err != nil {
			return fmt.Errorf("StockMovementService.CreateMovement: criar movimentação: %w", err)
		}

		// 4. Atualizar estoque do ingrediente DENTRO da transação (passando tx)
		ingredient.StockQuantity = newStock
		if err := s.productRepo.UpdateIngredient(ctx, ingredient, tx); err != nil {
			return fmt.Errorf("StockMovementService.CreateMovement: atualizar estoque: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return movement, nil
}

// ListStockMovements lista movimentações de estoque
func (s *StockMovementService) ListStockMovements(ctx context.Context, companyID uint, ingredientID *uint, limit, offset int) ([]domain.StockMovement, error) {
	return s.stockMovementRepo.List(ctx, companyID, ingredientID, limit, offset)
}

// GetStockMovementByID busca uma movimentação por ID
func (s *StockMovementService) GetStockMovementByID(ctx context.Context, id uint) (*domain.StockMovement, error) {
	return s.stockMovementRepo.GetByID(ctx, id)
}

// DeleteStockMovement deleta uma movimentação (soft delete)
func (s *StockMovementService) DeleteStockMovement(ctx context.Context, id uint) error {
	return s.stockMovementRepo.Delete(ctx, id)
}

// --- Inventários ---

// CreateInventory cria um inventário de estoque
func (s *StockMovementService) CreateInventory(ctx context.Context, companyID, userID uint, input CreateInventoryInput) (*domain.StockInventory, error) {
	inventory := &domain.StockInventory{
		CompanyID:     companyID,
		InventoryDate: input.InventoryDate,
		Status:        "draft",
		Notes:         input.Notes,
		PerformedBy:   userID,
	}

	// Sprint 4B.1 v2: Passar nil para tx (fora de transação)
	if err := s.stockMovementRepo.CreateInventory(ctx, inventory, nil); err != nil {
		return nil, fmt.Errorf("StockMovementService.CreateInventory: criar inventário: %w", err)
	}

	return inventory, nil
}

// GetInventoryByID busca um inventário por ID
// Sprint 4B.2: Passar nil para tx (fora de transação)
func (s *StockMovementService) GetInventoryByID(ctx context.Context, id uint) (*domain.StockInventory, error) {
	return s.stockMovementRepo.GetInventoryByID(ctx, id, nil)
}

// ListInventories lista inventários
func (s *StockMovementService) ListInventories(ctx context.Context, companyID uint, status string, limit, offset int) ([]domain.StockInventory, error) {
	return s.stockMovementRepo.ListInventories(ctx, companyID, status, limit, offset)
}

// AddInventoryItem adiciona um item ao inventário
// Sprint 4B.2: Passar nil para tx (fora de transação)
func (s *StockMovementService) AddInventoryItem(ctx context.Context, inventoryID, ingredientID uint, expectedStock, actualStock float64, reason string) (*domain.StockInventoryItem, error) {
	// Buscar inventário
	inventory, err := s.stockMovementRepo.GetInventoryByID(ctx, inventoryID, nil)
	if err != nil {
		return nil, ErrStockInventoryNotFound
	}

	// Validar status
	if inventory.Status != "draft" {
		return nil, ErrStockInventoryCompleted
	}

	// Calcular diferença
	difference := actualStock - expectedStock

	item := &domain.StockInventoryItem{
		InventoryID:   inventoryID,
		IngredientID:  ingredientID,
		ExpectedStock: expectedStock,
		ActualStock:   actualStock,
		Difference:    difference,
		Reason:        reason,
	}

	// Sprint 4B.1 v2: Passar nil para tx (fora de transação)
	if err := s.stockMovementRepo.CreateInventoryItem(ctx, item, nil); err != nil {
		return nil, fmt.Errorf("StockMovementService.CreateInventoryItem: criar item de inventário: %w", err)
	}

	return item, nil
}

// CompleteInventory completa um inventário e ajusta o estoque em transação atômica
// Sprint 4B.1 v2: Corrigido com transação real, SELECT FOR UPDATE e ordenação de locks
// Sprint 4B.2: Corrigido para propagar tx em GetInventoryByID e ListInventoryItems
// Sprint 4B.5: Adicionado SELECT FOR UPDATE no inventário para prevenir double completion
func (s *StockMovementService) CompleteInventory(ctx context.Context, inventoryID, userID uint) error {
	// Executar toda a operação em transação atômica
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Sprint 4B.5: Buscar inventário com SELECT FOR UPDATE DENTRO da transação
		// Isso previne double completion e modificações concorrentes
		inventory, err := s.stockMovementRepo.FindInventoryByIDForUpdate(ctx, inventoryID, tx)
		if err != nil {
			return ErrStockInventoryNotFound
		}

		// Validar status
		if inventory.Status != "draft" {
			return ErrStockInventoryCompleted
		}

		// 2. Buscar itens DENTRO da transação (passando tx)
		// Sprint 4B.5: Não é necessário SELECT FOR UPDATE aqui pois o inventário já está travado
		items, err := s.stockMovementRepo.ListInventoryItems(ctx, inventoryID, tx)
		if err != nil {
			return fmt.Errorf("StockMovementService.CompleteInventory: buscar itens do inventário: %w", err)
		}

		// Sprint 4B.1 v2: Ordenar itens por IngredientID para evitar deadlock
		// Isso garante ordem determinística de locks
		sort.Slice(items, func(i, j int) bool {
			return items[i].IngredientID < items[j].IngredientID
		})

		// 3. Ajustar estoque para cada item na mesma transação
		for _, item := range items {
			if item.Difference != 0 {
				// Sprint 4B.1 v2: Buscar ingrediente com SELECT FOR UPDATE dentro da transação
				ingredient, err := s.productRepo.FindIngredientByIDForUpdate(ctx, item.IngredientID, tx)
				if err != nil {
					return fmt.Errorf("StockMovementService.CompleteInventory: buscar ingrediente: %w", err)
				}
				if ingredient == nil {
					return fmt.Errorf("StockMovementService.CompleteInventory: ingrediente id=%d não encontrado", item.IngredientID)
				}

				// Validar company
				if ingredient.CompanyID != inventory.CompanyID {
					return errors.New("ingrediente não pertence a esta empresa")
				}

				// Calcular novo estoque
				previousStock := ingredient.StockQuantity
				var newStock float64
				var quantity float64

				if item.Difference > 0 {
					quantity = item.Difference
					newStock = previousStock + quantity
				} else {
					quantity = item.Difference // já é negativo
					newStock = previousStock + quantity
					if newStock < 0 {
						return fmt.Errorf("StockMovementService.CompleteInventory: estoque não pode ser negativo para ingrediente %d", item.IngredientID)
					}
				}

				// Criar movimentação DENTRO da transação (passando tx)
				movement := &domain.StockMovement{
					CompanyID:     inventory.CompanyID,
					IngredientID:  item.IngredientID,
					Type:          domain.StockMovementInventory,
					Quantity:      quantity,
					PreviousStock: previousStock,
					NewStock:      newStock,
					Reason:        fmt.Sprintf("Ajuste por inventário #%d: %s", inventoryID, item.Reason),
					ReferenceType: "inventory",
					ReferenceID:   inventoryID,
					PerformedBy:   userID,
					PerformedAt:   time.Now(),
				}

				if err := s.stockMovementRepo.Create(ctx, movement, tx); err != nil {
					return fmt.Errorf("StockMovementService.CompleteInventory: criar movimentação: %w", err)
				}

				// Atualizar estoque do ingrediente DENTRO da transação (passando tx)
				ingredient.StockQuantity = newStock
				if err := s.productRepo.UpdateIngredient(ctx, ingredient, tx); err != nil {
					return fmt.Errorf("StockMovementService.CompleteInventory: atualizar estoque: %w", err)
				}
			}
		}

		// Sprint 4B.5: Validação de status após o loop (defesa em profundidade)
		// Isso garante que o inventário não foi modificado durante o processamento
		currentInventory, err := s.stockMovementRepo.FindInventoryByIDForUpdate(ctx, inventoryID, tx)
		if err != nil {
			return fmt.Errorf("StockMovementService.CompleteInventory: verificar status do inventário: %w", err)
		}
		if currentInventory.Status != "draft" {
			return ErrStockInventoryCompleted
		}

		// 4. Atualizar status do inventário DENTRO da transação (passando tx)
		if err := s.stockMovementRepo.UpdateInventoryStatus(ctx, inventoryID, "completed", tx); err != nil {
			return fmt.Errorf("StockMovementService.CompleteInventory: atualizar status do inventário: %w", err)
		}

		return nil
	})
}

// DeleteInventory deleta um inventário
func (s *StockMovementService) DeleteInventory(ctx context.Context, id uint) error {
	return s.stockMovementRepo.DeleteInventory(ctx, id)
}

// --- Inputs ---

type CreateStockMovementInput struct {
	IngredientID  uint                     `json:"ingredientId" validate:"required"`
	Type          domain.StockMovementType `json:"type" validate:"required"`
	Quantity      float64                  `json:"quantity" validate:"required,gt=0"`
	Reason        string                   `json:"reason"`
	ReferenceType string                   `json:"referenceType"`
	ReferenceID   uint                     `json:"referenceId"`
}

type CreateInventoryInput struct {
	InventoryDate time.Time `json:"inventoryDate" validate:"required"`
	Notes         string    `json:"notes"`
}
