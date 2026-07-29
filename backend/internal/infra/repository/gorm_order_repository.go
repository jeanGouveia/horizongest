package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/infra/pg"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
	"gorm.io/gorm"
)

// ─── GORM models ────────────────────────────────────────────────────────────

type GormOrder struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	OrderNumber    int    `gorm:"not null;default:0;index:idx_orders_company_order_number,priority:1"`
	Status         string `gorm:"not null;default:'pending'"`
	TotalPrice     int64
	Notes          string
	CompanyID      uint       `gorm:"index;not null;index:idx_orders_company_order_number,priority:2"` // Sprint 3: NOT NULL
	IdempotencyKey *string    `gorm:"type:varchar(255);index"`                                         // Sprint 4C: Idempotency key
	DeletedAt      *time.Time `gorm:"index"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`
}

func (GormOrder) TableName() string { return "orders" }

type GormOrderItem struct {
	ID                    uint         `gorm:"primaryKey;autoIncrement"`
	OrderID               uint         `gorm:"not null;index"`
	ProductID             uint         `gorm:"not null"`
	Quantity              float64      `gorm:"not null"`
	UnitPrice             int64        `gorm:"not null"`
	ProductName           string       `gorm:"not null"`               // snapshot do nome
	ProductDescription    string       `gorm:"type:text"`              // snapshot da descrição
	ProductIsComposto     bool         `gorm:"not null;default:false"` // snapshot da flag
	ProductPhotoURL       string       // snapshot da foto
	ProductCategoryID     *uint        // snapshot da categoria
	ProductPromotionPrice *int64       // snapshot do preço promocional
	ProductFeatured       bool         `gorm:"not null;default:false"` // snapshot do destaque
	ProductIsNew          bool         `gorm:"not null;default:false"` // snapshot do selo novo
	DeletedAt             *time.Time   `gorm:"index"`
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
// Sprint 4C: Retorna o pedido criado ou o pedido existente em caso de colisão de idempotency_key
func (r *GormOrderRepository) CreateOrder(ctx context.Context, order *domain.Order, productIngredients map[uint][]domain.ProductIngredient) (*domain.Order, error) {
	// Auto-fill CompanyID from tenant context
	companyID, err := GetCompanyIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("CreateOrder: %w", err)
	}

	var createdOrder *domain.Order

	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1. Gerar o próximo número de pedido para esta empresa (isolamento de tenant)
		// SELECT MAX(order_number) + 1 FROM orders WHERE company_id = ?
		var nextOrderNumber int
		if err := tx.Model(&GormOrder{}).
			Where("company_id = ?", companyID).
			Select("COALESCE(MAX(order_number), 0) + 1").
			Scan(&nextOrderNumber).Error; err != nil {
			return fmt.Errorf("CreateOrder: gerar order_number: %w", err)
		}

		// 2. Persiste o pedido
		gOrder := GormOrder{
			OrderNumber:    nextOrderNumber,
			Status:         string(order.Status),
			TotalPrice:     int64(order.TotalPrice),
			Notes:          order.Notes,
			CompanyID:      companyID,            // Auto-filled from context
			IdempotencyKey: order.IdempotencyKey, // Sprint 4C: Idempotency key
		}
		if err := tx.Create(&gOrder).Error; err != nil {
			// BUG 2 FIX: Verificar se é erro de unique constraint usando SQLSTATE 23505
			if pg.IsUniqueViolation(err) {
				// Colisão de idempotency_key - retornar erro customizado para buscar fora da transação
				return WrapDuplicateKeyError(err)
			}
			return fmt.Errorf("CreateOrder: criar pedido: %w", err)
		}
		order.ID = gOrder.ID
		order.OrderNumber = gOrder.OrderNumber
		order.CompanyID = gOrder.CompanyID
		order.CreatedAt = gOrder.CreatedAt
		createdOrder = order // Pedido criado com sucesso

		// 2. Sprint 4B.3: Coletar todos os ingredientes únicos do pedido para ordenação de locks
		// Isso previne deadlock garantindo ordem global determinística
		type ingredientConsumption struct {
			ingredientID uint
			totalQty     float64
			name         string
			currentStock float64
		}

		ingredientMap := make(map[uint]*ingredientConsumption)

		for i := range order.Items {
			item := &order.Items[i]
			ingredients, ok := productIngredients[item.ProductID]
			if !ok {
				return fmt.Errorf("CreateOrder: ingredientes não pré-carregados para produto_id=%d", item.ProductID)
			}

			for _, pi := range ingredients {
				consumo := pi.Quantity * item.Quantity
				if existing, found := ingredientMap[pi.IngredientID]; found {
					existing.totalQty += consumo
				} else {
					ingredientMap[pi.IngredientID] = &ingredientConsumption{
						ingredientID: pi.IngredientID,
						totalQty:     consumo,
						name:         pi.Ingredient.Name,
						currentStock: pi.Ingredient.StockQuantity,
					}
				}
			}
		}

		// Converter para slice e ordenar por IngredientID (ordem global determinística)
		ingredientList := make([]*ingredientConsumption, 0, len(ingredientMap))
		for _, ic := range ingredientMap {
			ingredientList = append(ingredientList, ic)
		}
		sort.Slice(ingredientList, func(i, j int) bool {
			return ingredientList[i].ingredientID < ingredientList[j].ingredientID
		})

		// 3. Sprint 4B.3: Adquirir locks em ordem determinística
		// Isso garante que todas as transações adquiram locks na mesma ordem
		for _, ic := range ingredientList {
			consumo := ic.totalQty
			if err := r.productRepo.DecreaseIngredientStock(ctx, ic.ingredientID, consumo, tx, ic.name, ic.currentStock); err != nil {
				return fmt.Errorf("CreateOrder: baixa estoque ingrediente_id=%d: %w", ic.ingredientID, err)
			}
		}

		// 4. Persistir os itens do pedido (após locks já adquiridos)
		for i := range order.Items {
			item := &order.Items[i]
			item.OrderID = order.ID

			// 4a. Persiste o item com snapshot pré-carregado (princípio #4: Histórico é imutável)
			// O snapshot já foi montado no service para evitar chamadas dentro da transação
			gItem := GormOrderItem{
				OrderID:               item.OrderID,
				ProductID:             item.ProductID,
				Quantity:              item.Quantity,
				UnitPrice:             int64(item.UnitPrice),
				ProductName:           item.ProductName,                                      // snapshot do nome
				ProductDescription:    item.ProductDescription,                               // snapshot da descrição
				ProductIsComposto:     item.ProductIsComposto,                                // snapshot da flag
				ProductPhotoURL:       item.ProductPhotoURL,                                  // snapshot da foto
				ProductCategoryID:     item.ProductCategoryID,                                // snapshot da categoria
				ProductPromotionPrice: convertMoneyPtrToInt64Ptr(item.ProductPromotionPrice), // snapshot do preço promocional
				ProductFeatured:       item.ProductFeatured,                                  // snapshot do destaque
				ProductIsNew:          item.ProductIsNew,                                     // snapshot do selo novo
			}
			if err := tx.Create(&gItem).Error; err != nil {
				return fmt.Errorf("CreateOrder: criar item produto_id=%d: %w", item.ProductID, err)
			}
			item.ID = gItem.ID
		}

		return nil // commit
	})

	// Sprint 4C: Se houve erro de unique constraint, buscar pedido existente
	// BUG 2 FIX: Usar IsDuplicateKeyError para detectar erro de colisão
	if IsDuplicateKeyError(err) {
		// Buscar pedido existente pela idempotency_key
		if order.IdempotencyKey != nil {
			existingOrder, findErr := r.FindByIdempotencyKey(ctx, companyID, *order.IdempotencyKey)
			if findErr != nil {
				return nil, fmt.Errorf("CreateOrder: buscar pedido após colisão: %w", findErr)
			}
			if existingOrder != nil {
				// Carregar itens do pedido existente
				orderWithItems, findErr := r.FindOrderByID(ctx, existingOrder.ID)
				if findErr != nil {
					return nil, fmt.Errorf("CreateOrder: carregar itens após colisão: %w", findErr)
				}
				return orderWithItems, nil
			}
		}
		// Se não encontrou, retornar erro original
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return createdOrder, nil
}

func (r *GormOrderRepository) FindByIdempotencyKey(ctx context.Context, companyID uint, idempotencyKey string) (*domain.Order, error) {
	var gOrder GormOrder
	query := r.db.WithContext(ctx).
		Where("company_id = ? AND idempotency_key = ? AND deleted_at IS NULL", companyID, idempotencyKey).
		First(&gOrder)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if query.Error != nil {
		return nil, fmt.Errorf("FindByIdempotencyKey: %w", query.Error)
	}

	return orderToDomain(&gOrder), nil
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
			UnitPrice:             domain.Money(gi.UnitPrice),
			ProductName:           gi.ProductName,                                      // snapshot
			ProductDescription:    gi.ProductDescription,                               // snapshot
			ProductIsComposto:     gi.ProductIsComposto,                                // snapshot
			ProductPhotoURL:       gi.ProductPhotoURL,                                  // snapshot
			ProductCategoryID:     gi.ProductCategoryID,                                // snapshot
			ProductPromotionPrice: convertInt64PtrToMoneyPtr(gi.ProductPromotionPrice), // snapshot
			ProductFeatured:       gi.ProductFeatured,                                  // snapshot
			ProductIsNew:          gi.ProductIsNew,                                     // snapshot
		}
		if gi.Product != nil {
			item.Product = &domain.Product{
				ID: gi.Product.ID, Name: gi.Product.Name,
				Description: gi.Product.Description, Price: domain.Money(gi.Product.Price),
				IsComposto: gi.Product.IsComposto, Active: gi.Product.Active,
				CreatedAt: gi.Product.CreatedAt,
				UpdatedAt: gi.Product.UpdatedAt,
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

			// Sprint 4B.1: Idempotência garantida por unique constraint no banco
			// Migration 00005: uk_stock_adjustments_order_ingredient_pending
			// Se houver violação de constraint, trataremos como idempotência (sucesso)
			// Isso evita race condition entre COUNT e INSERT

			ajustesCriados := 0
			ajustesPulados := 0 // Contador para idempotência
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
						// Sprint 4B.1: Tratar violação de unique constraint como idempotência
						// Se o ajuste já existe, não é erro - é idempotência
						if IsDuplicateKeyError(err) {
							log.Printf("[REPO] Ajuste já existe (idempotência): order_id=%d, ingredient_id=%d", id, pi.IngredientID)
							ajustesPulados++
							continue // Não é erro, apenas idempotente
						}
						log.Printf("[REPO] ===== ERRO NO INSERT =====")
						log.Printf("[REPO] order_id=%d, ingredient_id=%d", id, pi.IngredientID)
						log.Printf("[REPO] Erro: %v", err)
						return fmt.Errorf("UpdateOrderStatusWithAdjustments: criar ajuste: %w", err)
					}
					log.Printf("[REPO] INSERT bem-sucedido: order_id=%d, ingredient_id=%d", id, pi.IngredientID)
					ajustesCriados++
				}
			}
			log.Printf("[REPO] Total de ajustes criados: %d, pulados (idempotência): %d", ajustesCriados, ajustesPulados)
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
		// Sprint 4B.1 v2: Passar nil para tx (fora de transação)
		ingredient, err := r.productRepo.FindIngredientByID(ctx, ingredientID, nil)
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
	total domain.Money,
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
		// Sprint 4A: NOTA - Query sem filtro de tenant mas dentro de transação
		// NÃO é IDOR pois o order_id já foi validado por ApplyTenantFilterWithID
		// na linha anterior (query.ApplyTenantFilterWithID(ctx, tx, id))
		// Isso garante que o pedido pertence ao tenant antes de buscar itens
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
				UnitPrice: domain.Money(gi.UnitPrice),
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
		// Sprint 4A: NOTA - Delete sem filtro de tenant mas dentro de transação
		// NÃO é IDOR pois o order_id já foi validado por ApplyTenantFilterWithID
		// no início da transação. Só deleta itens do pedido validado.
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
				UnitPrice:             int64(item.UnitPrice),
				ProductName:           item.ProductName,
				ProductDescription:    item.ProductDescription,
				ProductIsComposto:     item.ProductIsComposto,
				ProductPhotoURL:       item.ProductPhotoURL,
				ProductCategoryID:     item.ProductCategoryID,
				ProductPromotionPrice: convertMoneyPtrToInt64Ptr(item.ProductPromotionPrice),
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
		dt := *g.DeletedAt
		deletedAt = &dt
	}
	return &domain.Order{
		ID:          g.ID,
		OrderNumber: g.OrderNumber,
		Status:      domain.OrderStatus(g.Status),
		TotalPrice:  domain.Money(g.TotalPrice),
		Notes:       g.Notes,
		CompanyID:   g.CompanyID,
		DeletedAt:   deletedAt,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
}

// isDuplicateKeyError verifica se o erro é uma violação de unique constraint
