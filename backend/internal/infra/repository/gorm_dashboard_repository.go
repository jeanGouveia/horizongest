package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
)

var _ ports.DashboardRepository = (*GormDashboardRepository)(nil)

type GormDashboardRepository struct{ db *gorm.DB }

func NewGormDashboardRepository(db *gorm.DB) *GormDashboardRepository {
	return &GormDashboardRepository{db: db}
}

func (r *GormDashboardRepository) GetDashboard(ctx context.Context) (*domain.Dashboard, error) {
	dashboard := &domain.Dashboard{}

	// Métricas
	today := time.Now().Format("2006-01-02")

	// Receita de hoje
	var todayRevenue float64
	r.db.WithContext(ctx).Model(&GormOrder{}).
		Where("DATE(FROM_UNIXTIME(created_at)) = ? AND deleted_at IS NULL", today).
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&todayRevenue)
	dashboard.Metrics.TodayRevenue = todayRevenue

	// Pedidos de hoje
	var todayOrders int64
	r.db.WithContext(ctx).Model(&GormOrder{}).
		Where("DATE(FROM_UNIXTIME(created_at)) = ? AND deleted_at IS NULL", today).
		Count(&todayOrders)
	dashboard.Metrics.TodayOrders = int(todayOrders)

	// Pedidos pendentes
	var pendingOrders int64
	r.db.WithContext(ctx).Model(&GormOrder{}).
		Where("status = ? AND deleted_at IS NULL", "pending").
		Count(&pendingOrders)
	dashboard.Metrics.PendingOrders = int(pendingOrders)

	// Estoque baixo
	var lowStockCount int64
	r.db.WithContext(ctx).Model(&GormIngredient{}).
		Where("stock_quantity < min_stock AND deleted_at IS NULL").
		Count(&lowStockCount)
	dashboard.Metrics.LowStockCount = int(lowStockCount)

	// Produtos ativos
	var activeProducts int64
	r.db.WithContext(ctx).Model(&GormProduct{}).
		Where("active = ? AND deleted_at IS NULL", true).
		Count(&activeProducts)
	dashboard.Metrics.ActiveProducts = int(activeProducts)

	// Total de produtos
	var totalProducts int64
	r.db.WithContext(ctx).Model(&GormProduct{}).
		Where("deleted_at IS NULL").
		Count(&totalProducts)
	dashboard.TotalProducts = int(totalProducts)

	// Total de categorias
	var totalCategories int64
	r.db.WithContext(ctx).Model(&GormCategory{}).
		Where("deleted_at IS NULL").
		Count(&totalCategories)
	dashboard.TotalCategories = int(totalCategories)

	// Total de ingredientes
	var totalIngredients int64
	r.db.WithContext(ctx).Model(&GormIngredient{}).
		Where("deleted_at IS NULL").
		Count(&totalIngredients)
	dashboard.TotalIngredients = int(totalIngredients)

	// Pedidos recentes (últimos 10)
	type OrderResult struct {
		ID         uint    `gorm:"column:id"`
		Status     string  `gorm:"column:status"`
		TotalPrice float64 `gorm:"column:total_price"`
		CreatedAt  int64   `gorm:"column:created_at"`
	}
	var recentOrders []OrderResult
	r.db.WithContext(ctx).Model(&GormOrder{}).
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
			ID:         o.ID,
			Status:     domain.OrderStatus(o.Status),
			TotalPrice: o.TotalPrice,
			CreatedAt:  time.Unix(o.CreatedAt, 0).Format("2006-01-02 15:04"),
			ItemsCount: int(itemsCount),
		}
	}

	// Estoque baixo (ingredientes com estoque abaixo do mínimo)
	type IngredientResult struct {
		ID            uint    `gorm:"column:id"`
		Name          string  `gorm:"column:name"`
		StockQuantity float64 `gorm:"column:stock_quantity"`
		MinStock      float64 `gorm:"column:min_stock"`
		Unit          string  `gorm:"column:unit"`
	}
	var lowStockItems []IngredientResult
	r.db.WithContext(ctx).Model(&GormIngredient{}).
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

	return dashboard, nil
}
