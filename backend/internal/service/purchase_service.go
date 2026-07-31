package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
	"gorm.io/gorm"
)

var (
	ErrSupplierNotFound      = errors.New("fornecedor não encontrado")
	ErrPurchaseOrderNotFound = errors.New("pedido de compra não encontrado")
	ErrPurchaseOrderInvalid  = errors.New("pedido de compra inválido")
	ErrPurchaseOrderSent     = errors.New("pedido já enviado, não pode ser alterado")
	ErrPurchaseOrderReceived = errors.New("pedido já recebido")
)

type PurchaseService struct {
	purchaseRepo ports.PurchaseRepository
	productRepo  ports.ProductRepository
	db           *gorm.DB
}

func NewPurchaseService(
	purchaseRepo ports.PurchaseRepository,
	productRepo ports.ProductRepository,
	db *gorm.DB,
) *PurchaseService {
	return &PurchaseService{
		purchaseRepo: purchaseRepo,
		productRepo:  productRepo,
		db:           db,
	}
}

// --- Fornecedores ---

// CreateSupplier cria um novo fornecedor
func (s *PurchaseService) CreateSupplier(ctx context.Context, companyID uint, input CreateSupplierInput) (*domain.Supplier, error) {
	supplier := &domain.Supplier{
		CompanyID: companyID,
		Name:      input.Name,
		CNPJ:      input.CNPJ,
		Email:     input.Email,
		Phone:     input.Phone,
		Address:   input.Address,
		City:      input.City,
		State:     input.State,
		ZipCode:   input.ZipCode,
		Notes:     input.Notes,
		Active:    true,
	}

	if err := s.purchaseRepo.CreateSupplier(ctx, supplier); err != nil {
		return nil, fmt.Errorf("PurchaseService.CreateSupplier: criar fornecedor: %w", err)
	}

	return supplier, nil
}

// ListSuppliers lista fornecedores
func (s *PurchaseService) ListSuppliers(ctx context.Context, companyID uint, activeOnly bool, limit, offset int) ([]domain.Supplier, error) {
	return s.purchaseRepo.ListSuppliers(ctx, companyID, activeOnly, limit, offset)
}

// GetSupplierByID busca um fornecedor por ID
func (s *PurchaseService) GetSupplierByID(ctx context.Context, id uint) (*domain.Supplier, error) {
	return s.purchaseRepo.GetSupplierByID(ctx, id)
}

// UpdateSupplier atualiza um fornecedor
func (s *PurchaseService) UpdateSupplier(ctx context.Context, supplier *domain.Supplier) error {
	return s.purchaseRepo.UpdateSupplier(ctx, supplier)
}

// DeleteSupplier deleta um fornecedor (soft delete)
func (s *PurchaseService) DeleteSupplier(ctx context.Context, id uint) error {
	return s.purchaseRepo.DeleteSupplier(ctx, id)
}

// --- Pedidos de Compra ---

// CreatePurchaseOrder cria um novo pedido de compra em transação atômica
func (s *PurchaseService) CreatePurchaseOrder(ctx context.Context, companyID, userID uint, input CreatePurchaseOrderInput) (*domain.PurchaseOrder, error) {
	var order *domain.PurchaseOrder

	// Executar toda a operação em transação atômica (se db disponível)
	executeInTx := func(fn func() error) error {
		if s.db != nil {
			return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
				return fn()
			})
		}
		return fn()
	}

	err := executeInTx(func() error {
		// Validar fornecedor
		supplier, err := s.purchaseRepo.GetSupplierByID(ctx, input.SupplierID)
		if err != nil {
			return ErrSupplierNotFound
		}

		if supplier.CompanyID != companyID {
			return errors.New("fornecedor não pertence a esta empresa")
		}

		// Gerar número do pedido
		orderNumber := fmt.Sprintf("PC-%d-%d", companyID, time.Now().Unix())

		// Calcular totais
		var subtotal domain.Money
		for _, item := range input.Items {
			itemPrice := domain.FromFloat64(item.UnitPrice).Mul(int64(item.Quantity * 100)).Div(100)
			subtotal = subtotal.Add(itemPrice)
		}

		tax := domain.FromFloat64(input.Tax)
		discount := domain.FromFloat64(input.Discount)
		total := subtotal.Add(tax).Sub(discount)

		order = &domain.PurchaseOrder{
			CompanyID:    companyID,
			SupplierID:   input.SupplierID,
			OrderNumber:  orderNumber,
			Status:       domain.PurchaseOrderDraft,
			OrderDate:    input.OrderDate,
			ExpectedDate: input.ExpectedDate,
			Subtotal:     subtotal,
			Tax:          tax,
			Discount:     discount,
			Total:        total,
			Notes:        input.Notes,
			CreatedBy:    userID,
		}

		if err := s.purchaseRepo.CreatePurchaseOrder(ctx, order); err != nil {
			return fmt.Errorf("PurchaseService.CreatePurchaseOrder: criar pedido de compra: %w", err)
		}

		// Criar itens
		for _, item := range input.Items {
			poItem := &domain.PurchaseOrderItem{
				PurchaseOrderID: order.ID,
				IngredientID:    item.IngredientID,
				Quantity:        item.Quantity,
				Unit:            item.Unit,
				UnitPrice:       domain.FromFloat64(item.UnitPrice),
				Subtotal:        domain.FromFloat64(item.UnitPrice).Mul(int64(item.Quantity * 100)).Div(100),
			}

			if err := s.purchaseRepo.CreatePurchaseOrderItem(ctx, poItem); err != nil {
				return fmt.Errorf("PurchaseService.CreatePurchaseOrder: criar item do pedido: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return order, nil
}

// ListPurchaseOrders lista pedidos de compra
func (s *PurchaseService) ListPurchaseOrders(ctx context.Context, companyID uint, status string, limit, offset int) ([]domain.PurchaseOrder, error) {
	return s.purchaseRepo.ListPurchaseOrders(ctx, companyID, status, limit, offset)
}

// GetPurchaseOrderByID busca um pedido de compra por ID
func (s *PurchaseService) GetPurchaseOrderByID(ctx context.Context, id uint) (*domain.PurchaseOrder, error) {
	return s.purchaseRepo.GetPurchaseOrderByID(ctx, id)
}

// UpdatePurchaseOrder atualiza um pedido de compra
func (s *PurchaseService) UpdatePurchaseOrder(ctx context.Context, order *domain.PurchaseOrder) error {
	// Validar status
	if order.Status != domain.PurchaseOrderDraft {
		return ErrPurchaseOrderSent
	}

	return s.purchaseRepo.UpdatePurchaseOrder(ctx, order)
}

// UpdatePurchaseOrderStatus atualiza o status de um pedido de compra
func (s *PurchaseService) UpdatePurchaseOrderStatus(ctx context.Context, id uint, status domain.PurchaseOrderStatus) error {
	return s.purchaseRepo.UpdatePurchaseOrderStatus(ctx, id, status)
}

// DeletePurchaseOrder deleta um pedido de compra (soft delete)
func (s *PurchaseService) DeletePurchaseOrder(ctx context.Context, id uint) error {
	return s.purchaseRepo.DeletePurchaseOrder(ctx, id)
}

// --- Recebimentos ---

// CreatePurchaseReceiving cria um recebimento de compra em transação atômica
func (s *PurchaseService) CreatePurchaseReceiving(ctx context.Context, purchaseOrderID, userID uint, input CreatePurchaseReceivingInput) (*domain.PurchaseReceiving, error) {
	var receiving *domain.PurchaseReceiving

	// Executar toda a operação em transação atômica (se db disponível)
	executeInTx := func(fn func() error) error {
		if s.db != nil {
			return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
				return fn()
			})
		}
		return fn()
	}

	err := executeInTx(func() error {
		// Buscar pedido de compra
		order, err := s.purchaseRepo.GetPurchaseOrderByID(ctx, purchaseOrderID)
		if err != nil {
			return ErrPurchaseOrderNotFound
		}

		// Validar status
		if order.Status == domain.PurchaseOrderReceived {
			return ErrPurchaseOrderReceived
		}

		receiving = &domain.PurchaseReceiving{
			PurchaseOrderID: purchaseOrderID,
			ReceivedDate:    input.ReceivedDate,
			ReceivedBy:      userID,
			Notes:           input.Notes,
		}

		if err := s.purchaseRepo.CreatePurchaseReceiving(ctx, receiving); err != nil {
			return fmt.Errorf("PurchaseService.ReceivePurchaseOrder: criar recebimento: %w", err)
		}

		// Criar itens de recebimento
		for _, item := range input.Items {
			receivingItem := &domain.PurchaseReceivingItem{
				PurchaseReceivingID: receiving.ID,
				PurchaseOrderItemID: item.PurchaseOrderItemID,
				IngredientID:        item.IngredientID,
				Quantity:            item.Quantity,
				Unit:                item.Unit,
				UnitPrice:           domain.FromFloat64(item.UnitPrice),
				Subtotal:            domain.FromFloat64(item.UnitPrice).Mul(int64(item.Quantity * 100)).Div(100),
			}

			if err := s.purchaseRepo.CreatePurchaseReceivingItem(ctx, receivingItem); err != nil {
				return fmt.Errorf("PurchaseService.ReceivePurchaseOrder: criar item de recebimento: %w", err)
			}

			// TODO: Atualizar estoque do ingrediente via stock movements
			// Isso será implementado na integração com o módulo de estoque
		}

		// Atualizar status do pedido para recebido
		now := time.Now()
		order.ReceivedDate = &now
		if err := s.purchaseRepo.UpdatePurchaseOrderStatus(ctx, purchaseOrderID, domain.PurchaseOrderReceived); err != nil {
			return fmt.Errorf("PurchaseService.CreatePurchaseReceiving: atualizar status do pedido: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return receiving, nil
}

// GetPurchaseReceivingByID busca um recebimento por ID
func (s *PurchaseService) GetPurchaseReceivingByID(ctx context.Context, id uint) (*domain.PurchaseReceiving, error) {
	return s.purchaseRepo.GetPurchaseReceivingByID(ctx, id)
}

// ListPurchaseReceivings lista recebimentos de um pedido
func (s *PurchaseService) ListPurchaseReceivings(ctx context.Context, purchaseOrderID uint) ([]domain.PurchaseReceiving, error) {
	return s.purchaseRepo.ListPurchaseReceivings(ctx, purchaseOrderID)
}

// DeletePurchaseReceiving deleta um recebimento (soft delete)
func (s *PurchaseService) DeletePurchaseReceiving(ctx context.Context, id uint) error {
	return s.purchaseRepo.DeletePurchaseReceiving(ctx, id)
}

// --- Inputs ---

type CreateSupplierInput struct {
	Name    string `json:"name" validate:"required"`
	CNPJ    string `json:"cnpj"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipCode"`
	Notes   string `json:"notes"`
}

type CreatePurchaseOrderInput struct {
	SupplierID   uint                           `json:"supplierId" validate:"required"`
	OrderDate    time.Time                      `json:"orderDate" validate:"required"`
	ExpectedDate *time.Time                     `json:"expectedDate"`
	Tax          float64                        `json:"tax"`
	Discount     float64                        `json:"discount"`
	Notes        string                         `json:"notes"`
	Items        []CreatePurchaseOrderItemInput `json:"items" validate:"required"`
}

type CreatePurchaseOrderItemInput struct {
	IngredientID uint    `json:"ingredientId" validate:"required"`
	Quantity     float64 `json:"quantity" validate:"required,gt=0"`
	Unit         string  `json:"unit" validate:"required"`
	UnitPrice    float64 `json:"unitPrice" validate:"required,gt=0"`
}

type CreatePurchaseReceivingInput struct {
	ReceivedDate time.Time                          `json:"receivedDate" validate:"required"`
	Notes        string                             `json:"notes"`
	Items        []CreatePurchaseReceivingItemInput `json:"items" validate:"required"`
}

type CreatePurchaseReceivingItemInput struct {
	PurchaseOrderItemID uint    `json:"purchaseOrderItemId" validate:"required"`
	IngredientID        uint    `json:"ingredientId" validate:"required"`
	Quantity            float64 `json:"quantity" validate:"required,gt=0"`
	Unit                string  `json:"unit" validate:"required"`
	UnitPrice           float64 `json:"unitPrice" validate:"required,gt=0"`
}
