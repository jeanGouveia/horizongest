package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

var ErrOrderNotFound = errors.New("pedido não encontrado")
var ErrInvalidOrderStatus = errors.New("status de pedido inválido")

// InsufficientStockError representa erro de estoque insuficiente com detalhes
type InsufficientStockError struct {
	Message     string
	Ingredients []InsufficientIngredient
}

type InsufficientIngredient struct {
	Name      string
	Available float64
	Required  float64
	Shortage  float64
	Unit      string
}

func (e *InsufficientStockError) Error() string {
	return e.Message
}

// NewInsufficientStockError cria um erro detalhado de estoque insuficiente
func NewInsufficientStockError(ingredients []InsufficientIngredient) *InsufficientStockError {
	msg := "Não foi possível concluir o pedido. Ingredientes insuficientes:\n"
	for _, ing := range ingredients {
		msg += fmt.Sprintf("• %s\n  Disponível: %.4f %s\n  Necessário: %.4f %s\n  Faltam: %.4f %s\n\n",
			ing.Name, ing.Available, ing.Unit, ing.Required, ing.Unit, ing.Shortage, ing.Unit)
	}
	return &InsufficientStockError{
		Message:     msg,
		Ingredients: ingredients,
	}
}

type OrderService struct {
	orderRepo   ports.OrderRepository
	productRepo ports.ProductRepository
}

func NewOrderService(orderRepo ports.OrderRepository, productRepo ports.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

// ── Inputs ───────────────────────────────────────────────────────────────────

type OrderItemInput struct {
	ProductID uint    `json:"product_id" validate:"required"`
	Quantity  float64 `json:"quantity"   validate:"required,gt=0"`
}

type CreateOrderInput struct {
	Items          []OrderItemInput `json:"items" validate:"required,min=1,dive"`
	Notes          string           `json:"notes"`
	IdempotencyKey *string          `json:"idempotency_key"` // Sprint 4C: Idempotency key
}

type UpdateOrderStatusInput struct {
	Status string `json:"status" validate:"required,oneof=pending confirmed preparing ready delivered cancelled"`
}

type UpdateOrderInput struct {
	Items []OrderItemInput `json:"items" validate:"required,min=1,dive"`
	Notes string           `json:"notes"`
}

// ── Operações ────────────────────────────────────────────────────────────────

func (s *OrderService) CreateOrder(ctx context.Context, in CreateOrderInput) (*domain.Order, error) {
	log.Printf("Service - Iniciando CreateOrder com %d itens", len(in.Items))

	// Sprint 4C: Verificar idempotência antes de criar pedido
	// Se a chave de idempotência for fornecida, buscar pedido existente
	if in.IdempotencyKey != nil && *in.IdempotencyKey != "" {
		// Obter companyID do contexto
		tenantCtxValue := ctx.Value("tenant")
		if tenantCtxValue == nil {
			return nil, fmt.Errorf("OrderService.CreateOrder: contexto de tenant não encontrado")
		}
		tenantCtx, ok := tenantCtxValue.(*domain.TenantContext)
		if !ok {
			return nil, fmt.Errorf("OrderService.CreateOrder: tipo de contexto de tenant inválido")
		}
		companyID := tenantCtx.CompanyID

		// Buscar pedido existente pela chave de idempotência
		existingOrder, err := s.orderRepo.FindByIdempotencyKey(ctx, companyID, *in.IdempotencyKey)
		if err != nil {
			return nil, fmt.Errorf("OrderService.CreateOrder: buscar por idempotency_key: %w", err)
		}
		if existingOrder != nil {
			log.Printf("Service - Pedido já existe (idempotência): ID=%d, Key=%s", existingOrder.ID, *in.IdempotencyKey)
			// Carregar itens do pedido existente
			orderWithItems, err := s.orderRepo.FindOrderByID(ctx, existingOrder.ID)
			if err != nil {
				return nil, fmt.Errorf("OrderService.CreateOrder: carregar itens do pedido existente: %w", err)
			}
			return orderWithItems, nil
		}
	}

	order := &domain.Order{
		Status:         domain.OrderStatusPending,
		Notes:          in.Notes,
		IdempotencyKey: in.IdempotencyKey, // Sprint 4C: Idempotency key
	}

	var total domain.Money

	// Pré-carrega produtos e fichas técnicas antes da transação para evitar context deadline
	productData := make(map[uint]*domain.Product)
	productIngredients := make(map[uint][]domain.ProductIngredient)

	for _, itemIn := range in.Items {
		// Buscar produto (fora da transação)
		p, err := s.productRepo.FindProductByID(ctx, itemIn.ProductID)
		if err != nil {
			return nil, fmt.Errorf("OrderService.CreateOrder: buscar produto: %w", err)
		}
		if p == nil || !p.Active {
			return nil, fmt.Errorf("OrderService.CreateOrder: produto id=%d não encontrado ou inativo", itemIn.ProductID)
		}
		productData[itemIn.ProductID] = p

		// Buscar ficha técnica (fora da transação)
		ingredients, err := s.productRepo.GetProductIngredients(ctx, itemIn.ProductID)
		if err != nil {
			return nil, fmt.Errorf("OrderService.CreateOrder: ficha técnica produto_id=%d: %w", itemIn.ProductID, err)
		}
		productIngredients[itemIn.ProductID] = ingredients
		log.Printf("Service - Produto %d tem %d ingredientes na ficha técnica", itemIn.ProductID, len(ingredients))
	}

	// Monta itens com snapshot completo (nome, descrição, preço, flag, campos comerciais)
	items := make([]domain.OrderItem, len(in.Items))
	for i, itemIn := range in.Items {
		p := productData[itemIn.ProductID]
		items[i] = domain.OrderItem{
			ProductID:             p.ID,
			Quantity:              itemIn.Quantity,
			UnitPrice:             p.Price,          // snapshot do preço
			ProductName:           p.Name,           // snapshot do nome
			ProductDescription:    p.Description,    // snapshot da descrição
			ProductIsComposto:     p.IsComposto,     // snapshot da flag
			ProductPhotoURL:       p.PhotoURL,       // snapshot da foto
			ProductCategoryID:     p.CategoryID,     // snapshot da categoria
			ProductPromotionPrice: p.PromotionPrice, // snapshot do preço promocional
			ProductFeatured:       p.Featured,       // snapshot do destaque
			ProductIsNew:          p.IsNew,          // snapshot do selo novo
		}
		total = total.Add(p.Price.Mul(int64(itemIn.Quantity * 100)).Div(100))
	}

	order.Items = items
	order.TotalPrice = total
	log.Printf("Service - Pedido montado: TotalPrice=%d, Items=%d", order.TotalPrice, len(order.Items))

	// Pré-validação de estoque: coletar TODOS os ingredientes insuficientes
	// antes de tentar criar o pedido, para retornar mensagem completa
	insufficientIngredients := s.validateStock(ctx, in.Items, productIngredients)
	if len(insufficientIngredients) > 0 {
		log.Printf("Service - Estoque insuficiente: %d ingredientes faltantes", len(insufficientIngredients))
		return nil, NewInsufficientStockError(insufficientIngredients)
	}

	// CreateOrder executa a baixa de estoque em transação
	// Passamos os snapshots pré-carregados para evitar chamadas dentro da transação
	// Sprint 4C: Repository retorna o pedido criado ou o pedido existente em caso de colisão
	// Sprint 2.1: Outbox event is created within the same transaction (Order + StockMovement + Outbox)
	createdOrder, err := s.orderRepo.CreateOrder(ctx, order, productIngredients)
	if err != nil {
		log.Printf("Service - Erro ao criar pedido no repository: %v", err)
		return nil, fmt.Errorf("OrderService.CreateOrder: criar pedido: %w", err)
	}

	log.Printf("Service - Pedido criado com sucesso: ID=%d", createdOrder.ID)
	return createdOrder, nil
}

// validateStock verifica se há estoque suficiente para todos os ingredientes
// e retorna uma lista de ingredientes insuficientes (não para no primeiro erro)
func (s *OrderService) validateStock(ctx context.Context, items []OrderItemInput, productIngredients map[uint][]domain.ProductIngredient) []InsufficientIngredient {
	var insufficient []InsufficientIngredient

	// Mapa para acumular consumo por ingrediente
	requiredByIngredient := make(map[uint]float64)

	// Calcular quantidade necessária de cada ingrediente
	for _, itemIn := range items {
		ingredients, ok := productIngredients[itemIn.ProductID]
		if !ok || len(ingredients) == 0 {
			// Produto simples sem ficha técnica: não há ingredientes para validar
			continue
		}

		for _, pi := range ingredients {
			required := pi.Quantity * itemIn.Quantity
			requiredByIngredient[pi.IngredientID] += required
		}
	}

	// Verificar estoque disponível para cada ingrediente necessário
	for ingredientID, required := range requiredByIngredient {
		// Sprint 4B.1 v2: Passar nil para tx (fora de transação)
		ing, err := s.productRepo.FindIngredientByID(ctx, ingredientID, nil)
		if err != nil || ing == nil {
			// Se não conseguir buscar o ingrediente, considera como insuficiente
			insufficient = append(insufficient, InsufficientIngredient{
				Name:      fmt.Sprintf("Ingrediente #%d", ingredientID),
				Available: 0,
				Required:  required,
				Shortage:  required,
				Unit:      "?",
			})
			continue
		}

		if ing.StockQuantity < required {
			insufficient = append(insufficient, InsufficientIngredient{
				Name:      ing.Name,
				Available: ing.StockQuantity,
				Required:  required,
				Shortage:  required - ing.StockQuantity,
				Unit:      ing.Unit,
			})
		}
	}

	return insufficient
}

func (s *OrderService) ListOrders(ctx context.Context) ([]domain.Order, error) {
	orders, err := s.orderRepo.ListOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("OrderService.ListOrders: listar pedidos: %w", err)
	}
	return orders, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id uint) (*domain.Order, error) {
	order, err := s.orderRepo.FindOrderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("OrderService.GetOrder: buscar pedido: %w", err)
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, id uint, in UpdateOrderStatusInput) (*domain.Order, error) {
	log.Printf("[SERVICE] UpdateOrderStatus chamado: order_id=%d, new_status=%s", id, in.Status)

	order, err := s.orderRepo.FindOrderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("OrderService.UpdateOrderStatus: buscar pedido: %w", err)
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// Validar transição de status
	newStatus := domain.OrderStatus(in.Status)
	if !isValidTransition(order.Status, newStatus) {
		return nil, fmt.Errorf("OrderService.UpdateOrderStatus: transição inválida: %s → %s", order.Status, newStatus)
	}

	log.Printf("[SERVICE] Status atual=%s, novo status=%s", order.Status, newStatus)

	// Se o pedido está sendo cancelado, usar transação atômica para:
	// 1. Atualizar status do pedido
	// 2. Registrar ajustes pendentes
	// Garante consistência: ou ambos sucedem, ou ambos falham
	if newStatus == domain.OrderStatusCancelled {
		log.Printf("[SERVICE] Entrou em condição de cancelamento")
		// Carregar fichas técnicas dos produtos do pedido
		productIngredients := make(map[uint][]domain.ProductIngredient)
		for _, item := range order.Items {
			ingredients, err := s.productRepo.GetProductIngredients(ctx, item.ProductID)
			if err != nil {
				return nil, fmt.Errorf("OrderService.UpdateOrderStatus: carregar ficha técnica produto_id=%d: %w", item.ProductID, err)
			}
			log.Printf("[SERVICE] Produto %d tem %d ingredientes", item.ProductID, len(ingredients))
			productIngredients[item.ProductID] = ingredients
		}

		log.Printf("[SERVICE] Total de produtos com ficha técnica: %d", len(productIngredients))
		log.Printf("[SERVICE] Total de itens no pedido: %d", len(order.Items))

		// Executar atualização de status e registro de ajustes em transação atômica
		if err := s.orderRepo.UpdateOrderStatusWithAdjustments(
			ctx,
			order.ID,
			newStatus,
			productIngredients,
			order.Items,
		); err != nil {
			log.Printf("[SERVICE] Erro em UpdateOrderStatusWithAdjustments: %v", err)
			return nil, fmt.Errorf("OrderService.UpdateOrderStatus: atualizar status com ajustes: %w", err)
		}
		log.Printf("[SERVICE] UpdateOrderStatusWithAdjustments concluído com sucesso")
	} else {
		// Para outros status, apenas atualizar o status
		if err := s.orderRepo.UpdateOrderStatus(ctx, id, newStatus); err != nil {
			return nil, fmt.Errorf("OrderService.UpdateOrderStatus: atualizar status: %w", err)
		}
	}

	order.Status = newStatus
	return order, nil
}

// UpdateOrder allows editing order items and notes
// Only allowed for orders in pending or confirmed status
func (s *OrderService) UpdateOrder(ctx context.Context, id uint, in UpdateOrderInput) (*domain.Order, error) {
	log.Printf("[SERVICE] UpdateOrder chamado: order_id=%d", id)

	order, err := s.orderRepo.FindOrderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("OrderService.UpdateOrder: buscar pedido: %w", err)
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}

	// Only allow editing pending or confirmed orders
	if order.Status != domain.OrderStatusPending && order.Status != domain.OrderStatusConfirmed {
		return nil, errors.New("não é possível editar pedidos que não estejam pendentes ou confirmados")
	}

	// Pre-load products and ingredients
	productData := make(map[uint]*domain.Product)
	productIngredients := make(map[uint][]domain.ProductIngredient)

	for _, itemIn := range in.Items {
		p, err := s.productRepo.FindProductByID(ctx, itemIn.ProductID)
		if err != nil {
			return nil, fmt.Errorf("OrderService.UpdateOrder: buscar produto: %w", err)
		}
		if p == nil || !p.Active {
			return nil, fmt.Errorf("OrderService.UpdateOrder: produto id=%d não encontrado ou inativo", itemIn.ProductID)
		}
		productData[itemIn.ProductID] = p

		ingredients, err := s.productRepo.GetProductIngredients(ctx, itemIn.ProductID)
		if err != nil {
			return nil, fmt.Errorf("OrderService.UpdateOrder: ficha técnica produto_id=%d: %w", itemIn.ProductID, err)
		}
		productIngredients[itemIn.ProductID] = ingredients
	}

	// Calculate new total and build items
	var total domain.Money
	items := make([]domain.OrderItem, len(in.Items))
	for i, itemIn := range in.Items {
		p := productData[itemIn.ProductID]
		items[i] = domain.OrderItem{
			ProductID:             p.ID,
			Quantity:              itemIn.Quantity,
			UnitPrice:             p.Price,
			ProductName:           p.Name,
			ProductDescription:    p.Description,
			ProductIsComposto:     p.IsComposto,
			ProductPhotoURL:       p.PhotoURL,
			ProductCategoryID:     p.CategoryID,
			ProductPromotionPrice: p.PromotionPrice,
			ProductFeatured:       p.Featured,
			ProductIsNew:          p.IsNew,
		}
		total = total.Add(p.Price.Mul(int64(itemIn.Quantity * 100)).Div(100))
	}

	// Validate stock for new items
	insufficientIngredients := s.validateStock(ctx, in.Items, productIngredients)
	if len(insufficientIngredients) > 0 {
		return nil, NewInsufficientStockError(insufficientIngredients)
	}

	// Update order with transaction to handle stock adjustments
	if err := s.orderRepo.UpdateOrder(ctx, id, items, total, in.Notes, productIngredients); err != nil {
		return nil, fmt.Errorf("OrderService.UpdateOrder: atualizar pedido: %w", err)
	}

	// Reload order to return updated state
	updatedOrder, err := s.orderRepo.FindOrderByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("OrderService.UpdateOrder: recarregar pedido: %w", err)
	}

	return updatedOrder, nil
}

// isValidTransition valida se a transição entre status é permitida
func isValidTransition(current, new domain.OrderStatus) bool {
	// Transições permitidas
	transitions := map[domain.OrderStatus][]domain.OrderStatus{
		domain.OrderStatusPending:   {domain.OrderStatusConfirmed, domain.OrderStatusCancelled},
		domain.OrderStatusConfirmed: {domain.OrderStatusPreparing, domain.OrderStatusCancelled},
		domain.OrderStatusPreparing: {domain.OrderStatusReady, domain.OrderStatusCancelled},
		domain.OrderStatusReady:     {domain.OrderStatusDelivered, domain.OrderStatusCancelled},
		domain.OrderStatusDelivered: {}, // status final, sem transições
		domain.OrderStatusCancelled: {}, // status final, sem transições
	}

	allowed, exists := transitions[current]
	if !exists {
		return false
	}

	for _, status := range allowed {
		if status == new {
			return true
		}
	}
	return false
}
