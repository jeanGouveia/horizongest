package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
)

// ─── GORM models ────────────────────────────────────────────────────────────

type GormOrder struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	Status     string `gorm:"not null;default:'pending'"`
	TotalPrice float64
	Notes      string
	CreatedAt  int64 `gorm:"autoCreateTime"`
	UpdatedAt  int64 `gorm:"autoUpdateTime"`
}

func (GormOrder) TableName() string { return "orders" }

type GormOrderItem struct {
	ID        uint         `gorm:"primaryKey;autoIncrement"`
	OrderID   uint         `gorm:"not null;index"`
	ProductID uint         `gorm:"not null"`
	Quantity  float64      `gorm:"not null"`
	UnitPrice float64      `gorm:"not null"`
	Product   *GormProduct `gorm:"foreignKey:ProductID"`
}

func (GormOrderItem) TableName() string { return "order_items" }

// ─── Repository ─────────────────────────────────────────────────────────────

var _ ports.OrderRepository = (*GormOrderRepository)(nil)

type GormOrderRepository struct {
	db                  *gorm.DB
	productRepo         ports.ProductRepository
	stockAdjustmentRepo *GormStockAdjustmentRepository
}

func NewGormOrderRepository(db *gorm.DB, productRepo ports.ProductRepository, stockAdjustmentRepo *GormStockAdjustmentRepository) *GormOrderRepository {
	return &GormOrderRepository{db: db, productRepo: productRepo, stockAdjustmentRepo: stockAdjustmentRepo}
}

// CreateOrder é a operação crítica: persiste pedido + itens + baixa de estoque
// em uma única transação. Qualquer falha reverte tudo.
func (r *GormOrderRepository) CreateOrder(ctx context.Context, order *domain.Order, productIngredients map[uint][]domain.ProductIngredient) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1. Persiste o pedido
		gOrder := GormOrder{
			Status:     string(order.Status),
			TotalPrice: order.TotalPrice,
			Notes:      order.Notes,
		}
		if err := tx.Create(&gOrder).Error; err != nil {
			return fmt.Errorf("CreateOrder: criar pedido: %w", err)
		}
		order.ID = gOrder.ID
		order.CreatedAt = time.Unix(gOrder.CreatedAt, 0)

		// 2. Para cada item do pedido
		for i := range order.Items {
			item := &order.Items[i]
			item.OrderID = order.ID

			// 2a. Persiste o item
			gItem := GormOrderItem{
				OrderID:   item.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
			}
			if err := tx.Create(&gItem).Error; err != nil {
				return fmt.Errorf("CreateOrder: criar item produto_id=%d: %w", item.ProductID, err)
			}
			item.ID = gItem.ID

			// 2b. Usa os ingredientes pré-carregados (evita chamada dentro da transação)
			ingredients, ok := productIngredients[item.ProductID]
			if !ok {
				return fmt.Errorf("CreateOrder: ingredientes não pré-carregados para produto_id=%d", item.ProductID)
			}

			if len(ingredients) == 0 {
				// Produto simples sem ficha técnica: não há ingredientes para baixar
				continue
			}

			// 2c. Para cada ingrediente da ficha técnica, dá baixa proporcional
			// Consumo = quantidade_na_ficha × quantidade_vendida
			for _, pi := range ingredients {
				consumo := pi.Quantity * item.Quantity
				if err := r.productRepo.DecreaseIngredientStock(ctx, pi.IngredientID, consumo, tx, pi.Ingredient.Name, pi.Ingredient.StockQuantity); err != nil {
					return fmt.Errorf("CreateOrder: baixa estoque: %w", err)
				}
			}
		}

		return nil // commit
	})
}

func (r *GormOrderRepository) FindOrderByID(ctx context.Context, id uint) (*domain.Order, error) {
	var gOrder GormOrder
	err := r.db.WithContext(ctx).First(&gOrder, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("FindOrderByID: %w", err)
	}

	var gItems []GormOrderItem
	if err := r.db.WithContext(ctx).
		Preload("Product").
		Where("order_id = ?", id).
		Find(&gItems).Error; err != nil {
		return nil, fmt.Errorf("FindOrderByID items: %w", err)
	}

	order := orderToDomain(&gOrder)
	order.Items = make([]domain.OrderItem, len(gItems))
	for i, gi := range gItems {
		item := domain.OrderItem{
			ID: gi.ID, OrderID: gi.OrderID,
			ProductID: gi.ProductID, Quantity: gi.Quantity, UnitPrice: gi.UnitPrice,
		}
		if gi.Product != nil {
			item.Product = &domain.Product{
				ID: gi.Product.ID, Name: gi.Product.Name,
				Description: gi.Product.Description, Price: gi.Product.Price,
				IsComposto: gi.Product.IsComposto, Active: gi.Product.Active,
				CreatedAt: time.Unix(gi.Product.CreatedAt, 0),
				UpdatedAt: time.Unix(gi.Product.UpdatedAt, 0),
			}
		}
		order.Items[i] = item
	}
	return order, nil
}

func (r *GormOrderRepository) ListOrders(ctx context.Context) ([]domain.Order, error) {
	var gOrders []GormOrder
	if err := r.db.WithContext(ctx).Order("created_at desc").Find(&gOrders).Error; err != nil {
		return nil, fmt.Errorf("ListOrders: %w", err)
	}
	out := make([]domain.Order, len(gOrders))
	for i, g := range gOrders {
		out[i] = *orderToDomain(&g)
	}
	return out, nil
}

func (r *GormOrderRepository) UpdateOrderStatus(
	ctx context.Context, id uint, status domain.OrderStatus,
) error {
	if err := r.db.WithContext(ctx).Model(&GormOrder{}).
		Where("id = ?", id).Update("status", string(status)).Error; err != nil {
		return fmt.Errorf("UpdateOrderStatus: %w", err)
	}
	return nil
}

func (r *GormOrderRepository) UpdateOrderStatusWithAdjustments(
	ctx context.Context,
	id uint,
	status domain.OrderStatus,
	productIngredients map[uint][]domain.ProductIngredient,
	orderItems []domain.OrderItem,
) error {
	log.Printf("[REPO] ===== INÍCIO UpdateOrderStatusWithAdjustments =====")
	log.Printf("[REPO] order_id=%d, novo_status=%s", id, status)
	log.Printf("[REPO] Total de itens no pedido: %d", len(orderItems))
	log.Printf("[REPO] Total de produtos com ficha técnica: %d", len(productIngredients))
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		log.Printf("[REPO] Transação iniciada para order_id=%d", id)
		// 1. Atualizar status do pedido
		if err := tx.Model(&GormOrder{}).
			Where("id = ?", id).
			Update("status", string(status)).Error; err != nil {
			return fmt.Errorf("UpdateOrderStatusWithAdjustments: atualizar status: %w", err)
		}
		log.Printf("[REPO] Status do pedido atualizado para %s", status)

		// 2. Se cancelado, registrar ajustes pendentes na mesma transação
		if status == domain.OrderStatusCancelled {
			log.Printf("[REPO] Entrou em condição de cancelamento")
			log.Printf("[REPO] Total de itens no pedido: %d", len(orderItems))
			log.Printf("[REPO] Total de produtos com ficha técnica: %d", len(productIngredients))

			// Verificar se já existem ajustes pendentes para este pedido (idempotência)
			var existingCount int64
			if err := tx.Model(&GormStockAdjustmentPending{}).
				Where("order_id = ? AND status = ?", id, domain.StockAdjustmentStatusPending).
				Count(&existingCount).Error; err != nil {
				return fmt.Errorf("UpdateOrderStatusWithAdjustments: verificar ajustes existentes: %w", err)
			}
			if existingCount > 0 {
				log.Printf("[REPO] Já existem %d ajustes pendentes para o pedido %d, pulando criação", existingCount, id)
				return nil // Sucesso sem criar novos ajustes (idempotente)
			}

			ajustesCriados := 0
			itemIndex := 0
			for _, item := range orderItems {
				itemIndex++
				log.Printf("[REPO] Processando item %d: item_id=%d, product_id=%d, quantity=%.4f", itemIndex, item.ID, item.ProductID, item.Quantity)
				ingredients, ok := productIngredients[item.ProductID]
				if !ok || len(ingredients) == 0 {
					// Produto simples ou sem ficha técnica: nada a registrar
					log.Printf("[REPO] Produto %d sem ficha técnica, pulando", item.ProductID)
					continue
				}

				log.Printf("[REPO] Produto %d tem %d ingredientes", item.ProductID, len(ingredients))
				ingredientIndex := 0
				for _, pi := range ingredients {
					ingredientIndex++
					consumedQuantity := pi.Quantity * item.Quantity
					log.Printf("[REPO] ===== ANTES DE INSERT =====")
					log.Printf("[REPO] order_id=%d, item_id=%d, product_id=%d, ingredient_id=%d", id, item.ID, item.ProductID, pi.IngredientID)
					log.Printf("[REPO] quantity=%.4f, novo_status=%s", consumedQuantity, status)
					// Usar o repository oficial com suporte a transação
					adjustment := &domain.StockAdjustmentPending{
						OrderID:      id,
						IngredientID: pi.IngredientID,
						Quantity:     consumedQuantity,
						OrderStatus:  string(status),
						Status:       domain.StockAdjustmentStatusPending,
					}
					if err := r.stockAdjustmentRepo.CreateStockAdjustmentPendingWithTx(ctx, adjustment, tx); err != nil {
						log.Printf("[REPO] ===== ERRO NO INSERT =====")
						log.Printf("[REPO] order_id=%d, ingredient_id=%d", id, pi.IngredientID)
						log.Printf("[REPO] Erro: %v", err)
						return fmt.Errorf("UpdateOrderStatusWithAdjustments: criar ajuste: %w", err)
					}
					log.Printf("[REPO] INSERT bem-sucedido: order_id=%d, ingredient_id=%d", id, pi.IngredientID)
					ajustesCriados++
				}
			}
			log.Printf("[REPO] Total de ajustes criados: %d", ajustesCriados)
		}

		log.Printf("[REPO] ===== FIM UpdateOrderStatusWithAdjustments =====")
		return nil
	})
}

func orderToDomain(g *GormOrder) *domain.Order {
	return &domain.Order{
		ID: g.ID, Status: domain.OrderStatus(g.Status),
		TotalPrice: g.TotalPrice, Notes: g.Notes,
		CreatedAt: time.Unix(g.CreatedAt, 0), UpdatedAt: time.Unix(g.UpdatedAt, 0),
	}
}
