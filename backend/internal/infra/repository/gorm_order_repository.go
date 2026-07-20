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
	CompanyID  uint   `gorm:"index;not null"` // Sprint 3: NOT NULL
	DeletedAt  *int64 `gorm:"index"`
	CreatedAt  int64  `gorm:"autoCreateTime"`
	UpdatedAt  int64  `gorm:"autoUpdateTime"`
}

func (GormOrder) TableName() string { return "orders" }

type GormOrderItem struct {
	ID                    uint         `gorm:"primaryKey;autoIncrement"`
	OrderID               uint         `gorm:"not null;index"`
	ProductID             uint         `gorm:"not null"`
	Quantity              float64      `gorm:"not null"`
	UnitPrice             float64      `gorm:"not null"`
	ProductName           string       `gorm:"not null"`               // snapshot do nome
	ProductDescription    string       `gorm:"type:text"`              // snapshot da descrição
	ProductIsComposto     bool         `gorm:"not null;default:false"` // snapshot da flag
	ProductPhotoURL       string       // snapshot da foto
	ProductCategoryID     *uint        // snapshot da categoria
	ProductPromotionPrice *float64     // snapshot do preço promocional
	ProductFeatured       bool         `gorm:"not null;default:false"` // snapshot do destaque
	ProductIsNew          bool         `gorm:"not null;default:false"` // snapshot do selo novo
	DeletedAt             *int64       `gorm:"index"`
	Product               *GormProduct `gorm:"foreignKey:ProductID"`
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
	// Auto-fill CompanyID from tenant context
	companyID, err := GetCompanyIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("CreateOrder: %w", err)
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1. Persiste o pedido
		gOrder := GormOrder{
			Status:     string(order.Status),
			TotalPrice: order.TotalPrice,
			Notes:      order.Notes,
			CompanyID:  companyID, // Auto-filled from context
		}
		if err := tx.Create(&gOrder).Error; err != nil {
			return fmt.Errorf("CreateOrder: criar pedido: %w", err)
		}
		order.ID = gOrder.ID
		order.CompanyID = gOrder.CompanyID
		order.CreatedAt = time.Unix(gOrder.CreatedAt, 0)

		// 2. Para cada item do pedido
		for i := range order.Items {
			item := &order.Items[i]
			item.OrderID = order.ID

			// 2a. Persiste o item com snapshot pré-carregado (princípio #4: Histórico é imutável)
			// O snapshot já foi montado no service para evitar chamadas dentro da transação
			gItem := GormOrderItem{
				OrderID:               item.OrderID,
				ProductID:             item.ProductID,
				Quantity:              item.Quantity,
				UnitPrice:             item.UnitPrice,
				ProductName:           item.ProductName,           // snapshot do nome
				ProductDescription:    item.ProductDescription,    // snapshot da descrição
				ProductIsComposto:     item.ProductIsComposto,     // snapshot da flag
				ProductPhotoURL:       item.ProductPhotoURL,       // snapshot da foto
				ProductCategoryID:     item.ProductCategoryID,     // snapshot da categoria
				ProductPromotionPrice: item.ProductPromotionPrice, // snapshot do preço promocional
				ProductFeatured:       item.ProductFeatured,       // snapshot do destaque
				ProductIsNew:          item.ProductIsNew,          // snapshot do selo novo
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
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	err := query.Where("deleted_at IS NULL").First(&gOrder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("FindOrderByID: %w", err)
	}

	var gItems []GormOrderItem
	if err := r.db.WithContext(ctx).
		Preload("Product").
		Where("order_id = ? AND deleted_at IS NULL", id).
		Find(&gItems).Error; err != nil {
		return nil, fmt.Errorf("FindOrderByID items: %w", err)
	}

	order := orderToDomain(&gOrder)
	order.Items = make([]domain.OrderItem, len(gItems))
	for i, gi := range gItems {
		item := domain.OrderItem{
			ID:                    gi.ID,
			OrderID:               gi.OrderID,
			ProductID:             gi.ProductID,
			Quantity:              gi.Quantity,
			UnitPrice:             gi.UnitPrice,
			ProductName:           gi.ProductName,           // snapshot
			ProductDescription:    gi.ProductDescription,    // snapshot
			ProductIsComposto:     gi.ProductIsComposto,     // snapshot
			ProductPhotoURL:       gi.ProductPhotoURL,       // snapshot
			ProductCategoryID:     gi.ProductCategoryID,     // snapshot
			ProductPromotionPrice: gi.ProductPromotionPrice, // snapshot
			ProductFeatured:       gi.ProductFeatured,       // snapshot
			ProductIsNew:          gi.ProductIsNew,          // snapshot
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
	query := ApplyTenantFilter(ctx, r.db)
	if err := query.WithContext(ctx).Where("deleted_at IS NULL").Order("created_at desc").Find(&gOrders).Error; err != nil {
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
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	if err := query.WithContext(ctx).Model(&GormOrder{}).
		Where("deleted_at IS NULL").Update("status", string(status)).Error; err != nil {
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
		query := ApplyTenantFilterWithID(ctx, tx, id)
		if err := query.Model(&GormOrder{}).
			Where("deleted_at IS NULL").Update("status", string(status)).Error; err != nil {
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

func (r *GormOrderRepository) ValidateStock(ctx context.Context, items []domain.OrderItem, productIngredients map[uint][]domain.ProductIngredient) (*domain.StockValidationResponse, error) {
	response := &domain.StockValidationResponse{
		Valid:             true,
		InsufficientStock: []domain.InsufficientIngredient{},
	}

	// Mapa para acumular quantidade necessária por ingrediente
	requiredByIngredient := make(map[uint]float64)

	for _, item := range items {
		ingredients, exists := productIngredients[item.ProductID]
		if !exists || len(ingredients) == 0 {
			// Produto simples ou sem ficha técnica - não afeta estoque
			continue
		}

		for _, pi := range ingredients {
			requiredQty := pi.Quantity * item.Quantity
			requiredByIngredient[pi.IngredientID] += requiredQty
		}
	}

	// Verificar estoque atual de cada ingrediente necessário
	for ingredientID, requiredQty := range requiredByIngredient {
		ingredient, err := r.productRepo.FindIngredientByID(ctx, ingredientID)
		if err != nil {
			return nil, fmt.Errorf("ValidateStock: buscar ingrediente %d: %w", ingredientID, err)
		}
		if ingredient == nil {
			return nil, fmt.Errorf("ValidateStock: ingrediente %d não encontrado", ingredientID)
		}

		if ingredient.StockQuantity < requiredQty {
			response.Valid = false
			response.InsufficientStock = append(response.InsufficientStock, domain.InsufficientIngredient{
				IngredientID:   ingredientID,
				IngredientName: ingredient.Name,
				Required:       requiredQty,
				Available:      ingredient.StockQuantity,
				Shortage:       requiredQty - ingredient.StockQuantity,
				Unit:           ingredient.Unit,
			})
		}
	}

	return response, nil
}

// UpdateOrder updates order items and notes with stock adjustment
func (r *GormOrderRepository) UpdateOrder(
	ctx context.Context,
	id uint,
	items []domain.OrderItem,
	total float64,
	notes string,
	productIngredients map[uint][]domain.ProductIngredient,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get current order to calculate stock adjustments
		var gOrder GormOrder
		query := ApplyTenantFilterWithID(ctx, tx, id)
		if err := query.Where("deleted_at IS NULL").First(&gOrder).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("pedido não encontrado")
			}
			return fmt.Errorf("UpdateOrder: buscar pedido: %w", err)
		}

		// Get current items
		var gItems []GormOrderItem
		if err := tx.Where("order_id = ? AND deleted_at IS NULL", id).Find(&gItems).Error; err != nil {
			return fmt.Errorf("UpdateOrder: buscar itens atuais: %w", err)
		}

		// Build map of current items for easy lookup
		currentItems := make(map[uint]domain.OrderItem)
		for _, gi := range gItems {
			currentItems[gi.ID] = domain.OrderItem{
				ID:        gi.ID,
				ProductID: gi.ProductID,
				Quantity:  gi.Quantity,
				UnitPrice: gi.UnitPrice,
			}
		}

		// Step 1: Add back stock for removed items or reduced quantities
		for _, gi := range gItems {
			// Find if this item still exists in new items
			var newItem *domain.OrderItem
			for _, item := range items {
				if item.ProductID == gi.ProductID {
					newItem = &item
					break
				}
			}

			if newItem == nil {
				// Item was completely removed - add back full stock
				ingredients, ok := productIngredients[gi.ProductID]
				if ok && len(ingredients) > 0 {
					for _, pi := range ingredients {
						consumo := pi.Quantity * gi.Quantity
						if err := r.productRepo.IncreaseIngredientStock(ctx, pi.IngredientID, consumo, tx); err != nil {
							return fmt.Errorf("UpdateOrder: restaurar estoque item removido: %w", err)
						}
					}
				}
			} else if newItem.Quantity < gi.Quantity {
				// Quantity reduced - add back difference
				difference := gi.Quantity - newItem.Quantity
				ingredients, ok := productIngredients[gi.ProductID]
				if ok && len(ingredients) > 0 {
					for _, pi := range ingredients {
						consumo := pi.Quantity * difference
						if err := r.productRepo.IncreaseIngredientStock(ctx, pi.IngredientID, consumo, tx); err != nil {
							return fmt.Errorf("UpdateOrder: restaurar estoque redução: %w", err)
						}
					}
				}
			}
		}

		// Step 2: Deduct stock for new items or increased quantities
		for _, item := range items {
			var oldItem *domain.OrderItem
			for _, gi := range gItems {
				if gi.ProductID == item.ProductID {
					oldItem = &domain.OrderItem{
						ID:       gi.ID,
						Quantity: gi.Quantity,
					}
					break
				}
			}

			if oldItem == nil {
				// New item added - deduct full stock
				ingredients, ok := productIngredients[item.ProductID]
				if ok && len(ingredients) > 0 {
					for _, pi := range ingredients {
						consumo := pi.Quantity * item.Quantity
						if err := r.productRepo.DecreaseIngredientStock(ctx, pi.IngredientID, consumo, tx, pi.Ingredient.Name, pi.Ingredient.StockQuantity); err != nil {
							return fmt.Errorf("UpdateOrder: baixar estoque novo item: %w", err)
						}
					}
				}
			} else if item.Quantity > oldItem.Quantity {
				// Quantity increased - deduct difference
				difference := item.Quantity - oldItem.Quantity
				ingredients, ok := productIngredients[item.ProductID]
				if ok && len(ingredients) > 0 {
					for _, pi := range ingredients {
						consumo := pi.Quantity * difference
						if err := r.productRepo.DecreaseIngredientStock(ctx, pi.IngredientID, consumo, tx, pi.Ingredient.Name, pi.Ingredient.StockQuantity); err != nil {
							return fmt.Errorf("UpdateOrder: baixar estoque aumento: %w", err)
						}
					}
				}
			}
		}

		// Step 3: Soft delete old items
		if err := tx.Where("order_id = ?", id).Delete(&GormOrderItem{}).Error; err != nil {
			return fmt.Errorf("UpdateOrder: deletar itens antigos: %w", err)
		}

		// Step 4: Insert new items
		for i := range items {
			item := &items[i]
			item.OrderID = id
			gItem := GormOrderItem{
				OrderID:               item.OrderID,
				ProductID:             item.ProductID,
				Quantity:              item.Quantity,
				UnitPrice:             item.UnitPrice,
				ProductName:           item.ProductName,
				ProductDescription:    item.ProductDescription,
				ProductIsComposto:     item.ProductIsComposto,
				ProductPhotoURL:       item.ProductPhotoURL,
				ProductCategoryID:     item.ProductCategoryID,
				ProductPromotionPrice: item.ProductPromotionPrice,
				ProductFeatured:       item.ProductFeatured,
				ProductIsNew:          item.ProductIsNew,
			}
			if err := tx.Create(&gItem).Error; err != nil {
				return fmt.Errorf("UpdateOrder: criar novo item: %w", err)
			}
			item.ID = gItem.ID
		}

		// Step 5: Update order total and notes
		if err := query.Model(&GormOrder{}).
			Where("deleted_at IS NULL").
			Updates(map[string]interface{}{
				"total_price": total,
				"notes":       notes,
			}).Error; err != nil {
			return fmt.Errorf("UpdateOrder: atualizar pedido: %w", err)
		}

		return nil
	})
}

func orderToDomain(g *GormOrder) *domain.Order {
	var deletedAt *time.Time
	if g.DeletedAt != nil {
		dt := time.Unix(*g.DeletedAt, 0)
		deletedAt = &dt
	}
	return &domain.Order{
		ID:         g.ID,
		Status:     domain.OrderStatus(g.Status),
		TotalPrice: g.TotalPrice,
		Notes:      g.Notes,
		CompanyID:  g.CompanyID,
		DeletedAt:  deletedAt,
		CreatedAt:  time.Unix(g.CreatedAt, 0),
		UpdatedAt:  time.Unix(g.UpdatedAt, 0),
	}
}
