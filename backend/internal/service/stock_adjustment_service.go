package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
)

var ErrStockAdjustmentNotFound = errors.New("ajuste de estoque não encontrado")

type StockAdjustmentService struct {
	stockAdjustmentRepo ports.StockAdjustmentRepository
	productRepo         ports.ProductRepository
}

func NewStockAdjustmentService(
	stockAdjustmentRepo ports.StockAdjustmentRepository,
	productRepo ports.ProductRepository,
) *StockAdjustmentService {
	return &StockAdjustmentService{
		stockAdjustmentRepo: stockAdjustmentRepo,
		productRepo:         productRepo,
	}
}

// ── Inputs ───────────────────────────────────────────────────────────────────

type CreateStockAdjustmentInput struct {
	OrderID      uint    `json:"order_id" validate:"required"`
	IngredientID uint    `json:"ingredient_id" validate:"required"`
	Quantity     float64 `json:"quantity" validate:"required,gt=0"`
	OrderStatus  string  `json:"order_status" validate:"required"`
}

type UpdateStockAdjustmentStatusInput struct {
	Status string `json:"status" validate:"required,oneof=pending approved rejected"`
}

// ── Operações ────────────────────────────────────────────────────────────────

// RegisterStockAdjustmentForOrder registra ajustes de estoque pendentes para um pedido cancelado
// Este método NÃO devolve estoque automaticamente, apenas registra para análise manual
func (s *StockAdjustmentService) RegisterStockAdjustmentForOrder(
	ctx context.Context,
	orderID uint,
	orderStatus domain.OrderStatus,
	productIngredients map[uint][]domain.ProductIngredient,
	orderItems []domain.OrderItem,
) error {
	// Validar se já existem ajustes pendentes para este pedido (prevenção de duplicatas)
	existing, err := s.stockAdjustmentRepo.FindPendingByOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("verificar ajustes existentes: %w", err)
	}
	if len(existing) > 0 {
		return fmt.Errorf("já existem %d ajustes pendentes para o pedido %d", len(existing), orderID)
	}

	// Para cada item do pedido
	for _, item := range orderItems {
		ingredients, ok := productIngredients[item.ProductID]
		if !ok || len(ingredients) == 0 {
			// Produto simples ou sem ficha técnica: nada a registrar
			continue
		}

		// Para cada ingrediente da ficha técnica
		for _, pi := range ingredients {
			// Calcular quantidade consumida
			consumedQuantity := pi.Quantity * item.Quantity

			// Registrar ajuste pendente
			adjustment := &domain.StockAdjustmentPending{
				OrderID:      orderID,
				IngredientID: pi.IngredientID,
				Quantity:     consumedQuantity,
				OrderStatus:  string(orderStatus),
				Status:       domain.StockAdjustmentStatusPending,
			}

			if err := s.stockAdjustmentRepo.CreateStockAdjustmentPending(ctx, adjustment); err != nil {
				return fmt.Errorf("RegisterStockAdjustmentForOrder: %w", err)
			}
		}
	}

	return nil
}

// ListPendingAdjustments lista todos os ajustes pendentes para aprovação
func (s *StockAdjustmentService) ListPendingAdjustments(
	ctx context.Context,
) ([]domain.StockAdjustmentPending, error) {
	adjustments, err := s.stockAdjustmentRepo.ListPending(ctx)
	if err != nil {
		return nil, fmt.Errorf("ListPendingAdjustments: %w", err)
	}
	return adjustments, nil
}

// ListPendingAdjustmentsWithFilters lista ajustes com filtros opcionais
func (s *StockAdjustmentService) ListPendingAdjustmentsWithFilters(
	ctx context.Context,
	statusFilter string,
	orderIDFilter *uint,
	ingredientIDFilter *uint,
) ([]domain.StockAdjustmentPending, error) {
	var adjustments []domain.StockAdjustmentPending
	var err error

	// Se não há filtros, usar método existente
	if statusFilter == "" && orderIDFilter == nil && ingredientIDFilter == nil {
		adjustments, err = s.stockAdjustmentRepo.ListPending(ctx)
	} else {
		// Aplicar filtros no nível de aplicação (simplificado para MVP)
		adjustments, err = s.stockAdjustmentRepo.ListPending(ctx)
		if err != nil {
			return nil, fmt.Errorf("ListPendingAdjustmentsWithFilters: %w", err)
		}

		// Filtrar em memória (para MVP, pode ser otimizado com query SQL no futuro)
		filtered := make([]domain.StockAdjustmentPending, 0)
		for _, adj := range adjustments {
			if statusFilter != "" && string(adj.Status) != statusFilter {
				continue
			}
			if orderIDFilter != nil && adj.OrderID != *orderIDFilter {
				continue
			}
			if ingredientIDFilter != nil && adj.IngredientID != *ingredientIDFilter {
				continue
			}
			filtered = append(filtered, adj)
		}
		adjustments = filtered
	}

	if err != nil {
		return nil, fmt.Errorf("ListPendingAdjustmentsWithFilters: %w", err)
	}
	return adjustments, nil
}

// ListPending lista todos os ajustes pendentes (sem filtros)
func (s *StockAdjustmentService) ListPending(ctx context.Context) ([]domain.StockAdjustmentPending, error) {
	return s.stockAdjustmentRepo.ListPending(ctx)
}

// GetAdjustmentsByOrder lista todos os ajustes (todos os status) para um pedido
func (s *StockAdjustmentService) GetAdjustmentsByOrder(
	ctx context.Context, orderID uint,
) ([]domain.StockAdjustmentPending, error) {
	adjustments, err := s.stockAdjustmentRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("GetAdjustmentsByOrder: %w", err)
	}
	return adjustments, nil
}

// GetPendingAdjustmentsByOrder lista apenas os ajustes pendentes para um pedido
func (s *StockAdjustmentService) GetPendingAdjustmentsByOrder(
	ctx context.Context, orderID uint,
) ([]domain.StockAdjustmentPending, error) {
	adjustments, err := s.stockAdjustmentRepo.FindPendingByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("GetPendingAdjustmentsByOrder: %w", err)
	}
	return adjustments, nil
}

// GetPendingAdjustmentsByIngredient lista ajustes pendentes para um ingrediente específico
func (s *StockAdjustmentService) GetPendingAdjustmentsByIngredient(
	ctx context.Context, ingredientID uint,
) ([]domain.StockAdjustmentPending, error) {
	adjustments, err := s.stockAdjustmentRepo.FindPendingByIngredientID(ctx, ingredientID)
	if err != nil {
		return nil, fmt.Errorf("GetPendingAdjustmentsByIngredient: %w", err)
	}
	return adjustments, nil
}

// GetAdjustmentByID busca um ajuste específico por ID
func (s *StockAdjustmentService) GetAdjustmentByID(
	ctx context.Context, id uint,
) (*domain.StockAdjustmentPending, error) {
	adjustment, err := s.stockAdjustmentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetAdjustmentByID: %w", err)
	}
	if adjustment == nil {
		return nil, ErrStockAdjustmentNotFound
	}
	return adjustment, nil
}

// ApproveAdjustment aprova um ajuste de estoque pendente e repõe o estoque do ingrediente.
func (s *StockAdjustmentService) ApproveAdjustment(
	ctx context.Context, id uint, processedBy uint, notes string,
) error {
	// Verificar se o ajuste existe e está pendente
	adjustment, err := s.stockAdjustmentRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("ApproveAdjustment: %w", err)
	}
	if adjustment == nil {
		return ErrStockAdjustmentNotFound
	}

	// Validar status
	if adjustment.Status != domain.StockAdjustmentStatusPending {
		return fmt.Errorf("ajuste já processado (status: %s)", adjustment.Status)
	}

	// Aprovar e repor estoque em transação atômica
	if err := s.stockAdjustmentRepo.ApproveAndRestoreStock(ctx, id, processedBy, notes); err != nil {
		return fmt.Errorf("ApproveAdjustment: %w", err)
	}

	return nil
}

// RejectAdjustment rejeita um ajuste de estoque pendente
// Não altera estoque, apenas registra decisão operacional
func (s *StockAdjustmentService) RejectAdjustment(
	ctx context.Context, id uint, processedBy uint, notes string,
) error {
	// Verificar se o ajuste existe e está pendente
	adjustment, err := s.stockAdjustmentRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("RejectAdjustment: %w", err)
	}
	if adjustment == nil {
		return ErrStockAdjustmentNotFound
	}

	// Validar status
	if adjustment.Status != domain.StockAdjustmentStatusPending {
		return fmt.Errorf("ajuste já processado (status: %s)", adjustment.Status)
	}

	// Rejeitar no repository
	if err := s.stockAdjustmentRepo.Reject(ctx, id, processedBy, notes); err != nil {
		return fmt.Errorf("RejectAdjustment: %w", err)
	}

	return nil
}
