package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
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
}

func NewStockMovementService(
	stockMovementRepo ports.StockMovementRepository,
	productRepo ports.ProductRepository,
) *StockMovementService {
	return &StockMovementService{
		stockMovementRepo: stockMovementRepo,
		productRepo:       productRepo,
	}
}

// --- Movimentações ---

// CreateStockMovement cria uma movimentação de estoque
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

	// Buscar ingrediente atual através do product repository
	ingredient, err := s.productRepo.FindIngredientByID(ctx, input.IngredientID)
	if err != nil {
		return nil, fmt.Errorf("ingrediente não encontrado: %w", err)
	}

	// Validar company
	if ingredient.CompanyID != companyID {
		return nil, errors.New("ingrediente não pertence a esta empresa")
	}

	// Calcular novo estoque
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
			return nil, errors.New("estoque não pode ser negativo")
		}
	}

	// Criar movimentação
	movement := &domain.StockMovement{
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

	if err := s.stockMovementRepo.Create(ctx, movement); err != nil {
		return nil, fmt.Errorf("erro ao criar movimentação: %w", err)
	}

	// Atualizar estoque do ingrediente
	ingredient.StockQuantity = newStock
	if err := s.productRepo.UpdateIngredient(ctx, ingredient); err != nil {
		return nil, fmt.Errorf("erro ao atualizar estoque: %w", err)
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

	if err := s.stockMovementRepo.CreateInventory(ctx, inventory); err != nil {
		return nil, fmt.Errorf("erro ao criar inventário: %w", err)
	}

	return inventory, nil
}

// GetInventoryByID busca um inventário por ID
func (s *StockMovementService) GetInventoryByID(ctx context.Context, id uint) (*domain.StockInventory, error) {
	return s.stockMovementRepo.GetInventoryByID(ctx, id)
}

// ListInventories lista inventários
func (s *StockMovementService) ListInventories(ctx context.Context, companyID uint, status string, limit, offset int) ([]domain.StockInventory, error) {
	return s.stockMovementRepo.ListInventories(ctx, companyID, status, limit, offset)
}

// AddInventoryItem adiciona um item ao inventário
func (s *StockMovementService) AddInventoryItem(ctx context.Context, inventoryID, ingredientID uint, expectedStock, actualStock float64, reason string) (*domain.StockInventoryItem, error) {
	// Buscar inventário
	inventory, err := s.stockMovementRepo.GetInventoryByID(ctx, inventoryID)
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

	if err := s.stockMovementRepo.CreateInventoryItem(ctx, item); err != nil {
		return nil, fmt.Errorf("erro ao criar item de inventário: %w", err)
	}

	return item, nil
}

// CompleteInventory completa um inventário e ajusta o estoque
func (s *StockMovementService) CompleteInventory(ctx context.Context, inventoryID, userID uint) error {
	// Buscar inventário
	inventory, err := s.stockMovementRepo.GetInventoryByID(ctx, inventoryID)
	if err != nil {
		return ErrStockInventoryNotFound
	}

	// Validar status
	if inventory.Status != "draft" {
		return ErrStockInventoryCompleted
	}

	// Buscar itens
	items, err := s.stockMovementRepo.ListInventoryItems(ctx, inventoryID)
	if err != nil {
		return fmt.Errorf("erro ao buscar itens do inventário: %w", err)
	}

	// Ajustar estoque para cada item
	for _, item := range items {
		if item.Difference != 0 {
			// Criar movimentação de ajuste
			input := CreateStockMovementInput{
				IngredientID:  item.IngredientID,
				Type:          domain.StockMovementInventory,
				Quantity:      item.Difference,
				Reason:        fmt.Sprintf("Ajuste por inventário #%d: %s", inventoryID, item.Reason),
				ReferenceType: "inventory",
				ReferenceID:   inventoryID,
			}

			_, err := s.CreateStockMovement(ctx, inventory.CompanyID, userID, input)
			if err != nil {
				return fmt.Errorf("erro ao ajustar estoque do ingrediente %d: %w", item.IngredientID, err)
			}
		}
	}

	// Atualizar status do inventário
	if err := s.stockMovementRepo.UpdateInventoryStatus(ctx, inventoryID, "completed"); err != nil {
		return fmt.Errorf("erro ao atualizar status do inventário: %w", err)
	}

	return nil
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
