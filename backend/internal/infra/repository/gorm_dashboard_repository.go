package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

var _ ports.DashboardRepository = (*GormDashboardRepository)(nil)

type GormDashboardRepository struct{ db *gorm.DB }

func NewGormDashboardRepository(db *gorm.DB) *GormDashboardRepository {
	return &GormDashboardRepository{db: db}
}

func (r *GormDashboardRepository) GetDashboard(ctx context.Context) (*domain.Dashboard, error) {
	dashboard := &domain.Dashboard{}

	// Usar ApplyTenantFilter para isolamento de tenant
	query := ApplyTenantFilter(ctx, r.db)

	now := time.Now()
	today := now.Format("2006-01-02")
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	weekStart := now.AddDate(0, 0, -int(now.Weekday())).Format("2006-01-02")
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")

	// --- KPIs Hoje ---
	// Receita de hoje
	var todayRevenue int64
	query.WithContext(ctx).Model(&GormOrder{}).
		Where("date(datetime(created_at, 'unixepoch')) = ? AND deleted_at IS NULL", today).
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&todayRevenue)
	dashboard.Metrics.TodayRevenue = domain.Money(todayRevenue)

	// Pedidos de hoje
	var todayOrders int64
	query.WithContext(ctx).Model(&GormOrder{}).
		Where("date(datetime(created_at, 'unixepoch')) = ? AND deleted_at IS NULL", today).
		Count(&todayOrders)
	dashboard.Metrics.TodayOrders = int(todayOrders)

	// Produtos vendidos hoje
	var todayProductsSold int64
	query.WithContext(ctx).Model(&GormOrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("date(datetime(orders.created_at, 'unixepoch')) = ? AND orders.deleted_at IS NULL AND order_items.deleted_at IS NULL", today).
		Select("COALESCE(SUM(order_items.quantity), 0)").
		Scan(&todayProductsSold)
	dashboard.Metrics.TodayProductsSold = int(todayProductsSold)

	// Ticket médio hoje
	if todayOrders > 0 {
		avgTicket := domain.Money(todayRevenue).Div(int64(todayOrders))
		dashboard.Metrics.TodayAverageTicket = avgTicket
	}

	// CMV hoje (simplificado - assume custo fixo de 30%)
	todayCMV := domain.Money(todayRevenue).MulFloat(0.30)
	dashboard.Metrics.TodayCMV = todayCMV
	dashboard.Metrics.TodayGrossProfit = domain.Money(todayRevenue).Sub(todayCMV)

	// --- KPIs Ontem ---
	// Receita de ontem
	var yesterdayRevenue int64
	query.WithContext(ctx).Model(&GormOrder{}).
		Where("date(datetime(created_at, 'unixepoch')) = ? AND deleted_at IS NULL", yesterday).
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&yesterdayRevenue)
	dashboard.Metrics.YesterdayRevenue = domain.Money(yesterdayRevenue)

	// Pedidos de ontem
	var yesterdayOrders int64
	query.WithContext(ctx).Model(&GormOrder{}).
		Where("date(datetime(created_at, 'unixepoch')) = ? AND deleted_at IS NULL", yesterday).
		Count(&yesterdayOrders)
	dashboard.Metrics.YesterdayOrders = int(yesterdayOrders)

	// Produtos vendidos ontem
	var yesterdayProductsSold int64
	query.WithContext(ctx).Model(&GormOrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("date(datetime(orders.created_at, 'unixepoch')) = ? AND orders.deleted_at IS NULL AND order_items.deleted_at IS NULL", yesterday).
		Select("COALESCE(SUM(order_items.quantity), 0)").
		Scan(&yesterdayProductsSold)
	dashboard.Metrics.YesterdayProductsSold = int(yesterdayProductsSold)

	// Ticket médio ontem
	if yesterdayOrders > 0 {
		avgTicket := domain.Money(yesterdayRevenue).Div(int64(yesterdayOrders))
		dashboard.Metrics.YesterdayAverageTicket = avgTicket
	}

	// --- KPIs Semana ---
	// Receita da semana
	var weekRevenue int64
	query.WithContext(ctx).Model(&GormOrder{}).
		Where("date(datetime(created_at, 'unixepoch')) >= ? AND deleted_at IS NULL", weekStart).
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&weekRevenue)
	dashboard.Metrics.WeekRevenue = domain.Money(weekRevenue)

	// Pedidos da semana
	var weekOrders int64
	query.WithContext(ctx).Model(&GormOrder{}).
		Where("date(datetime(created_at, 'unixepoch')) >= ? AND deleted_at IS NULL", weekStart).
		Count(&weekOrders)
	dashboard.Metrics.WeekOrders = int(weekOrders)

	// Produtos vendidos na semana
	var weekProductsSold int64
	query.WithContext(ctx).Model(&GormOrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("date(datetime(orders.created_at, 'unixepoch')) >= ? AND orders.deleted_at IS NULL AND order_items.deleted_at IS NULL", weekStart).
		Select("COALESCE(SUM(order_items.quantity), 0)").
		Scan(&weekProductsSold)
	dashboard.Metrics.WeekProductsSold = int(weekProductsSold)

	// Ticket médio da semana
	if weekOrders > 0 {
		avgTicket := domain.Money(weekRevenue).Div(int64(weekOrders))
		dashboard.Metrics.WeekAverageTicket = avgTicket
	}

	// --- KPIs Mês ---
	// Receita do mês
	var monthRevenue int64
	query.WithContext(ctx).Model(&GormOrder{}).
		Where("date(datetime(created_at, 'unixepoch')) >= ? AND deleted_at IS NULL", monthStart).
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&monthRevenue)
	dashboard.Metrics.MonthRevenue = domain.Money(monthRevenue)

	// Pedidos do mês
	var monthOrders int64
	query.WithContext(ctx).Model(&GormOrder{}).
		Where("date(datetime(created_at, 'unixepoch')) >= ? AND deleted_at IS NULL", monthStart).
		Count(&monthOrders)
	dashboard.Metrics.MonthOrders = int(monthOrders)

	// Produtos vendidos no mês
	var monthProductsSold int64
	query.WithContext(ctx).Model(&GormOrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("date(datetime(orders.created_at, 'unixepoch')) >= ? AND orders.deleted_at IS NULL AND order_items.deleted_at IS NULL", monthStart).
		Select("COALESCE(SUM(order_items.quantity), 0)").
		Scan(&monthProductsSold)
	dashboard.Metrics.MonthProductsSold = int(monthProductsSold)

	// Ticket médio do mês
	if monthOrders > 0 {
		avgTicket := domain.Money(monthRevenue).Div(int64(monthOrders))
		dashboard.Metrics.MonthAverageTicket = avgTicket
	}

	// --- KPIs Gerais ---
	// Pedidos pendentes
	var pendingOrders int64
	query.WithContext(ctx).Model(&GormOrder{}).
		Where("status = ? AND deleted_at IS NULL", "pending").
		Count(&pendingOrders)
	dashboard.Metrics.PendingOrders = int(pendingOrders)

	// Pedidos cancelados
	var cancelledOrders int64
	query.WithContext(ctx).Model(&GormOrder{}).
		Where("status = ? AND deleted_at IS NULL", "cancelled").
		Count(&cancelledOrders)
	dashboard.Metrics.CancelledOrders = int(cancelledOrders)

	// Estoque baixo (inclui zerado)
	var lowStockCount int64
	query.WithContext(ctx).Model(&GormIngredient{}).
		Where("stock_quantity < min_stock AND deleted_at IS NULL").
		Count(&lowStockCount)
	dashboard.Metrics.LowStockCount = int(lowStockCount)

	// Estoque zerado
	var zeroStockCount int64
	query.WithContext(ctx).Model(&GormIngredient{}).
		Where("stock_quantity = 0 AND deleted_at IS NULL").
		Count(&zeroStockCount)
	dashboard.Metrics.ZeroStockCount = int(zeroStockCount)

	// Produtos ativos
	var activeProducts int64
	query.WithContext(ctx).Model(&GormProduct{}).
		Where("active = ? AND deleted_at IS NULL", true).
		Count(&activeProducts)
	dashboard.Metrics.ActiveProducts = int(activeProducts)

	// Total de produtos
	var totalProducts int64
	query.WithContext(ctx).Model(&GormProduct{}).
		Where("deleted_at IS NULL").
		Count(&totalProducts)
	dashboard.TotalProducts = int(totalProducts)

	// Total de categorias
	var totalCategories int64
	query.WithContext(ctx).Model(&GormCategory{}).
		Where("deleted_at IS NULL").
		Count(&totalCategories)
	dashboard.TotalCategories = int(totalCategories)

	// Total de ingredientes
	var totalIngredients int64
	query.WithContext(ctx).Model(&GormIngredient{}).
		Where("deleted_at IS NULL").
		Count(&totalIngredients)
	dashboard.TotalIngredients = int(totalIngredients)

	// --- Pedidos recentes (últimos 10) ---
	type OrderResult struct {
		ID          uint      `gorm:"column:id"`
		OrderNumber int       `gorm:"column:order_number"`
		Status      string    `gorm:"column:status"`
		TotalPrice  int64     `gorm:"column:total_price"`
		CreatedAt   time.Time `gorm:"column:created_at"`
	}
	var recentOrders []OrderResult
	query.WithContext(ctx).Model(&GormOrder{}).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(10).
		Find(&recentOrders)

	dashboard.RecentOrders = make([]domain.RecentOrder, len(recentOrders))
	for i, o := range recentOrders {
		// Contar itens
		var itemsCount int64
		r.db.WithContext(ctx).Model(&GormOrderItem{}).
			Where("order_id = ? AND deleted_at IS NULL", o.ID).
			Count(&itemsCount)

		dashboard.RecentOrders[i] = domain.RecentOrder{
			ID:          o.ID,
			OrderNumber: o.OrderNumber,
			Status:      domain.OrderStatus(o.Status),
			TotalPrice:  domain.Money(o.TotalPrice),
			CreatedAt:   o.CreatedAt.Format("2006-01-02 15:04"),
			ItemsCount:  int(itemsCount),
		}
	}

	// --- Estoque baixo (ingredientes com estoque abaixo do mínimo) ---
	type IngredientResult struct {
		ID            uint    `gorm:"column:id"`
		Name          string  `gorm:"column:name"`
		StockQuantity float64 `gorm:"column:stock_quantity"`
		MinStock      float64 `gorm:"column:min_stock"`
		Unit          string  `gorm:"column:unit"`
	}
	var lowStockItems []IngredientResult
	query.WithContext(ctx).Model(&GormIngredient{}).
		Where("stock_quantity < min_stock AND deleted_at IS NULL").
		Order("stock_quantity ASC").
		Limit(10).
		Find(&lowStockItems)

	dashboard.LowStock = make([]domain.LowStockItem, len(lowStockItems))
	for i, item := range lowStockItems {
		dashboard.LowStock[i] = domain.LowStockItem{
			ID:            item.ID,
			Name:          item.Name,
			StockQuantity: item.StockQuantity,
			MinStock:      item.MinStock,
			Unit:          item.Unit,
		}
	}

	// --- Estoque zerado ---
	var zeroStockItems []IngredientResult
	query.WithContext(ctx).Model(&GormIngredient{}).
		Where("stock_quantity = 0 AND deleted_at IS NULL").
		Order("name ASC").
		Limit(10).
		Find(&zeroStockItems)

	dashboard.ZeroStock = make([]domain.LowStockItem, len(zeroStockItems))
	for i, item := range zeroStockItems {
		dashboard.ZeroStock[i] = domain.LowStockItem{
			ID:            item.ID,
			Name:          item.Name,
			StockQuantity: item.StockQuantity,
			MinStock:      item.MinStock,
			Unit:          item.Unit,
		}
	}

	// --- Gráficos ---
	// Vendas por dia (últimos 7 dias)
	dashboard.Charts.SalesByDay = r.getSalesByDay(ctx, now)

	// Vendas por hora (hoje)
	dashboard.Charts.SalesByHour = r.getSalesByHour(ctx, today)

	// Produtos mais vendidos (últimos 30 dias)
	dashboard.Charts.TopProducts = r.getTopProducts(ctx, now)

	// Categorias mais vendidas (últimos 30 dias)
	dashboard.Charts.TopCategories = r.getTopCategories(ctx, now)

	return dashboard, nil
}

// getSalesByDay retorna vendas por dia (últimos 7 dias)
func (r *GormDashboardRepository) getSalesByDay(ctx context.Context, now time.Time) []domain.ChartPoint {
	type DayResult struct {
		Date  string `gorm:"column:date"`
		Total int64  `gorm:"column:total"`
	}

	query := ApplyTenantFilter(ctx, r.db)

	var results []DayResult
	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		var total int64
		query.WithContext(ctx).Model(&GormOrder{}).
			Where("date(datetime(created_at, 'unixepoch')) = ? AND deleted_at IS NULL", date).
			Select("COALESCE(SUM(total_price), 0)").
			Scan(&total)
		results = append(results, DayResult{Date: date, Total: total})
	}

	chartPoints := make([]domain.ChartPoint, len(results))
	for i, r := range results {
		// Format date as "DD/MM"
		parsedTime, _ := time.Parse("2006-01-02", r.Date)
		label := parsedTime.Format("02/01")
		chartPoints[i] = domain.ChartPoint{Label: label, Value: domain.Money(r.Total)}
	}

	return chartPoints
}

// getSalesByHour retorna vendas por hora (hoje)
func (r *GormDashboardRepository) getSalesByHour(ctx context.Context, today string) []domain.ChartPoint {
	type HourResult struct {
		Hour  int   `gorm:"column:hour"`
		Total int64 `gorm:"column:total"`
	}

	query := ApplyTenantFilter(ctx, r.db)

	var results []HourResult
	query.WithContext(ctx).Model(&GormOrder{}).
		Select("strftime('%H', datetime(created_at, 'unixepoch')) as hour, COALESCE(SUM(total_price), 0) as total").
		Where("date(datetime(created_at, 'unixepoch')) = ? AND deleted_at IS NULL", today).
		Group("hour").
		Order("hour ASC").
		Scan(&results)

	chartPoints := make([]domain.ChartPoint, len(results))
	for i, r := range results {
		label := fmt.Sprintf("%02d:00", r.Hour)
		chartPoints[i] = domain.ChartPoint{Label: label, Value: domain.Money(r.Total)}
	}

	return chartPoints
}

// getTopProducts retorna produtos mais vendidos (últimos 30 dias)
func (r *GormDashboardRepository) getTopProducts(ctx context.Context, now time.Time) []domain.TopItem {
	type ProductResult struct {
		ID       uint   `gorm:"column:id"`
		Name     string `gorm:"column:name"`
		Quantity int64  `gorm:"column:quantity"`
		Total    int64  `gorm:"column:total"`
	}

	query := ApplyTenantFilter(ctx, r.db)

	monthAgo := now.AddDate(0, 0, -30).Format("2006-01-02")
	var results []ProductResult
	query.WithContext(ctx).Model(&GormOrderItem{}).
		Select("products.id, products.name, SUM(order_items.quantity) as quantity, SUM(order_items.quantity * order_items.unit_price) as total").
		Joins("JOIN products ON products.id = order_items.product_id").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("date(datetime(orders.created_at, 'unixepoch')) >= ? AND orders.deleted_at IS NULL AND order_items.deleted_at IS NULL AND products.deleted_at IS NULL", monthAgo).
		Group("products.id, products.name").
		Order("quantity DESC").
		Limit(10).
		Scan(&results)

	topItems := make([]domain.TopItem, len(results))
	for i, r := range results {
		topItems[i] = domain.TopItem{
			ID:    r.ID,
			Name:  r.Name,
			Value: domain.Money(r.Total),
			Count: int(r.Quantity),
		}
	}

	return topItems
}

// getTopCategories retorna categorias mais vendidas (últimos 30 dias)
func (r *GormDashboardRepository) getTopCategories(ctx context.Context, now time.Time) []domain.TopItem {
	type CategoryResult struct {
		ID       uint   `gorm:"column:id"`
		Name     string `gorm:"column:name"`
		Quantity int64  `gorm:"column:quantity"`
		Total    int64  `gorm:"column:total"`
	}

	query := ApplyTenantFilter(ctx, r.db)

	monthAgo := now.AddDate(0, 0, -30).Format("2006-01-02")
	var results []CategoryResult
	query.WithContext(ctx).Model(&GormOrderItem{}).
		Select("categories.id, categories.name, SUM(order_items.quantity) as quantity, SUM(order_items.quantity * order_items.unit_price) as total").
		Joins("JOIN products ON products.id = order_items.product_id").
		Joins("JOIN categories ON categories.id = products.category_id").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("date(datetime(orders.created_at, 'unixepoch')) >= ? AND orders.deleted_at IS NULL AND order_items.deleted_at IS NULL AND products.deleted_at IS NULL AND categories.deleted_at IS NULL", monthAgo).
		Group("categories.id, categories.name").
		Order("quantity DESC").
		Limit(10).
		Scan(&results)

	topItems := make([]domain.TopItem, len(results))
	for i, r := range results {
		topItems[i] = domain.TopItem{
			ID:    r.ID,
			Name:  r.Name,
			Value: domain.Money(r.Total),
			Count: int(r.Quantity),
		}
	}

	return topItems
}
