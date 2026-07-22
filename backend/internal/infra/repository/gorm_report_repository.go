package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

var _ ports.ReportRepository = (*GormReportRepository)(nil)

type GormReportRepository struct {
	db *gorm.DB
}

func NewGormReportRepository(db *gorm.DB) *GormReportRepository {
	return &GormReportRepository{db: db}
}

// GetSalesReport retorna relatório de vendas
func (r *GormReportRepository) GetSalesReport(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.SalesReport, error) {
	report := &domain.SalesReport{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Total de pedidos
	var totalOrders int64
	r.db.WithContext(ctx).Model(&GormOrder{}).
		Where("company_id = ? AND DATE(FROM_UNIXTIME(created_at)) >= ? AND DATE(FROM_UNIXTIME(created_at)) <= ? AND deleted_at IS NULL", companyID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Count(&totalOrders)
	report.TotalOrders = int(totalOrders)

	// Receita total
	var totalRevenue float64
	r.db.WithContext(ctx).Model(&GormOrder{}).
		Where("company_id = ? AND DATE(FROM_UNIXTIME(created_at)) >= ? AND DATE(FROM_UNIXTIME(created_at)) <= ? AND deleted_at IS NULL", companyID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&totalRevenue)
	report.TotalRevenue = totalRevenue

	// Ticket médio
	if totalOrders > 0 {
		report.AverageTicket = totalRevenue / float64(totalOrders)
	}

	// Produtos vendidos
	var productsSold int64
	r.db.WithContext(ctx).Model(&GormOrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.company_id = ? AND DATE(FROM_UNIXTIME(orders.created_at)) >= ? AND DATE(FROM_UNIXTIME(orders.created_at)) <= ? AND orders.deleted_at IS NULL AND order_items.deleted_at IS NULL", companyID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Select("COALESCE(SUM(order_items.quantity), 0)").
		Scan(&productsSold)
	report.ProductsSold = int(productsSold)

	// Pedidos cancelados
	var cancelledOrders int64
	r.db.WithContext(ctx).Model(&GormOrder{}).
		Where("company_id = ? AND status = ? AND DATE(FROM_UNIXTIME(created_at)) >= ? AND DATE(FROM_UNIXTIME(created_at)) <= ? AND deleted_at IS NULL", companyID, "cancelled", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Count(&cancelledOrders)
	report.CancelledOrders = int(cancelledOrders)

	// Top produtos (simplificado - retorna até 10)
	type ProductResult struct {
		ID       uint    `gorm:"column:id"`
		Name     string  `gorm:"column:name"`
		Quantity int64   `gorm:"column:quantity"`
		Total    float64 `gorm:"column:total"`
	}
	var topProducts []ProductResult
	r.db.WithContext(ctx).Model(&GormOrderItem{}).
		Select("products.id, products.name, SUM(order_items.quantity) as quantity, SUM(order_items.quantity * order_items.unit_price) as total").
		Joins("JOIN products ON products.id = order_items.product_id").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.company_id = ? AND DATE(FROM_UNIXTIME(orders.created_at)) >= ? AND DATE(FROM_UNIXTIME(orders.created_at)) <= ? AND orders.deleted_at IS NULL AND order_items.deleted_at IS NULL AND products.deleted_at IS NULL", companyID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Group("products.id, products.name").
		Order("quantity DESC").
		Limit(10).
		Scan(&topProducts)

	report.TopProducts = make([]domain.TopItem, len(topProducts))
	for i, p := range topProducts {
		report.TopProducts[i] = domain.TopItem{
			ID:    p.ID,
			Name:  p.Name,
			Value: p.Total,
			Count: int(p.Quantity),
		}
	}

	// Vendas por dia (simplificado - últimos 7 dias do período)
	report.SalesByDay = r.getSalesByDayForReport(ctx, companyID, startDate, endDate)

	return report, nil
}

// GetProductsReport retorna relatório de produtos
func (r *GormReportRepository) GetProductsReport(ctx context.Context, companyID uint) (*domain.ProductsReport, error) {
	report := &domain.ProductsReport{}

	// Total de produtos
	var totalProducts int64
	r.db.WithContext(ctx).Model(&GormProduct{}).
		Where("company_id = ? AND deleted_at IS NULL", companyID).
		Count(&totalProducts)
	report.TotalProducts = int(totalProducts)

	// Produtos ativos
	var activeProducts int64
	r.db.WithContext(ctx).Model(&GormProduct{}).
		Where("company_id = ? AND active = ? AND deleted_at IS NULL", companyID, true).
		Count(&activeProducts)
	report.ActiveProducts = int(activeProducts)

	// Produtos arquivados
	report.ArchivedProducts = report.TotalProducts - report.ActiveProducts

	// Produtos sem ficha técnica (simplificado - assume que produtos sem ingredientes não têm ficha)
	// Esta query é complexa, vamos simplificar
	report.ProductsWithoutRecipe = 0 // Placeholder

	return report, nil
}

// GetCMVReport retorna relatório de CMV
func (r *GormReportRepository) GetCMVReport(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.CMVReport, error) {
	report := &domain.CMVReport{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Receita total (reutilizando lógica de vendas)
	var totalRevenue float64
	r.db.WithContext(ctx).Model(&GormOrder{}).
		Where("company_id = ? AND DATE(FROM_UNIXTIME(created_at)) >= ? AND DATE(FROM_UNIXTIME(created_at)) <= ? AND deleted_at IS NULL", companyID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&totalRevenue)
	report.TotalRevenue = totalRevenue

	// CMV simplificado (30% fixo - placeholder)
	report.TotalCMV = totalRevenue * 0.30
	if totalRevenue > 0 {
		report.CMVPercentage = (report.TotalCMV / totalRevenue) * 100
	}
	report.GrossProfit = totalRevenue - report.TotalCMV
	if totalRevenue > 0 {
		report.ProfitMargin = (report.GrossProfit / totalRevenue) * 100
	}

	// Por produto (simplificado - placeholder)
	report.ByProduct = []domain.ProductCMVItem{}

	return report, nil
}

// GetProfitReport retorna relatório de lucro
func (r *GormReportRepository) GetProfitReport(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.ProfitReport, error) {
	report := &domain.ProfitReport{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Receita total
	var totalRevenue float64
	r.db.WithContext(ctx).Model(&GormOrder{}).
		Where("company_id = ? AND DATE(FROM_UNIXTIME(created_at)) >= ? AND DATE(FROM_UNIXTIME(created_at)) <= ? AND deleted_at IS NULL", companyID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&totalRevenue)
	report.TotalRevenue = totalRevenue

	// Despesa total (do módulo financeiro)
	var totalExpense float64
	r.db.WithContext(ctx).Model(&GormTransaction{}).
		Where("company_id = ? AND type = 'expense' AND date >= ? AND date <= ? AND deleted_at IS NULL", companyID, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense)
	report.TotalExpense = totalExpense

	// Lucro líquido
	report.NetProfit = totalRevenue - totalExpense

	// Margem de lucro
	if totalRevenue > 0 {
		report.ProfitMargin = (report.NetProfit / totalRevenue) * 100
	}

	// Por categoria (simplificado - placeholder)
	report.ByCategory = []domain.CategoryProfitItem{}

	return report, nil
}

// GetStockReport retorna relatório de estoque
func (r *GormReportRepository) GetStockReport(ctx context.Context, companyID uint) (*domain.StockReport, error) {
	report := &domain.StockReport{}

	// Total de ingredientes
	var totalIngredients int64
	r.db.WithContext(ctx).Model(&GormIngredient{}).
		Where("company_id = ? AND deleted_at IS NULL", companyID).
		Count(&totalIngredients)
	report.TotalIngredients = int(totalIngredients)

	// Estoque baixo
	var lowStockCount int64
	r.db.WithContext(ctx).Model(&GormIngredient{}).
		Where("company_id = ? AND stock_quantity < min_stock AND stock_quantity > 0 AND deleted_at IS NULL", companyID).
		Count(&lowStockCount)
	report.LowStockCount = int(lowStockCount)

	// Estoque zerado
	var zeroStockCount int64
	r.db.WithContext(ctx).Model(&GormIngredient{}).
		Where("company_id = ? AND stock_quantity = 0 AND deleted_at IS NULL", companyID).
		Count(&zeroStockCount)
	report.ZeroStockCount = int(zeroStockCount)

	// Valor total do estoque (simplificado - assume custo unitário de 0)
	report.TotalStockValue = 0 // Placeholder

	// Itens com estoque baixo (reutilizar lógica do dashboard)
	type IngredientResult struct {
		ID            uint    `gorm:"column:id"`
		Name          string  `gorm:"column:name"`
		StockQuantity float64 `gorm:"column:stock_quantity"`
		MinStock      float64 `gorm:"column:min_stock"`
		Unit          string  `gorm:"column:unit"`
	}
	var lowStockItems []IngredientResult
	r.db.WithContext(ctx).Model(&GormIngredient{}).
		Where("company_id = ? AND stock_quantity < min_stock AND stock_quantity > 0 AND deleted_at IS NULL", companyID).
		Order("stock_quantity ASC").
		Limit(10).
		Find(&lowStockItems)

	report.LowStockItems = make([]domain.LowStockItem, len(lowStockItems))
	for i, item := range lowStockItems {
		report.LowStockItems[i] = domain.LowStockItem{
			ID:            item.ID,
			Name:          item.Name,
			StockQuantity: item.StockQuantity,
			MinStock:      item.MinStock,
			Unit:          item.Unit,
		}
	}

	// Itens com estoque zerado
	var zeroStockItems []IngredientResult
	r.db.WithContext(ctx).Model(&GormIngredient{}).
		Where("company_id = ? AND stock_quantity = 0 AND deleted_at IS NULL", companyID).
		Order("name ASC").
		Limit(10).
		Find(&zeroStockItems)

	report.ZeroStockItems = make([]domain.LowStockItem, len(zeroStockItems))
	for i, item := range zeroStockItems {
		report.ZeroStockItems[i] = domain.LowStockItem{
			ID:            item.ID,
			Name:          item.Name,
			StockQuantity: item.StockQuantity,
			MinStock:      item.MinStock,
			Unit:          item.Unit,
		}
	}

	// Itens com alto valor (placeholder)
	report.HighValueItems = []domain.StockValueItem{}

	return report, nil
}

// GetPurchasesReport retorna relatório de compras
func (r *GormReportRepository) GetPurchasesReport(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.PurchasesReport, error) {
	report := &domain.PurchasesReport{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Total de pedidos (placeholder - tabela de compras não integrada ainda)
	report.TotalOrders = 0
	report.TotalAmount = 0
	report.PendingOrders = 0
	report.ReceivedOrders = 0
	report.TopSuppliers = []domain.TopItem{}
	report.PurchasesByDay = []domain.ChartPoint{}

	return report, nil
}

// GetFinancialReport retorna relatório financeiro
func (r *GormReportRepository) GetFinancialReport(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.FinancialReport, error) {
	report := &domain.FinancialReport{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Receita total
	var totalIncome float64
	r.db.WithContext(ctx).Model(&GormTransaction{}).
		Where("company_id = ? AND type = 'income' AND date >= ? AND date <= ? AND deleted_at IS NULL", companyID, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome)
	report.TotalIncome = totalIncome

	// Despesa total
	var totalExpense float64
	r.db.WithContext(ctx).Model(&GormTransaction{}).
		Where("company_id = ? AND type = 'expense' AND date >= ? AND date <= ? AND deleted_at IS NULL", companyID, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense)
	report.TotalExpense = totalExpense

	// Saldo líquido
	report.NetBalance = totalIncome - totalExpense

	// Por categoria (placeholder)
	report.ByCategory = []domain.CategoryFinancialItem{}

	// Fluxo de caixa (placeholder)
	report.CashFlow = nil

	return report, nil
}

// Helper method para vendas por dia
func (r *GormReportRepository) getSalesByDayForReport(ctx context.Context, companyID uint, startDate, endDate time.Time) []domain.ChartPoint {
	type DayResult struct {
		Date  string  `gorm:"column:date"`
		Total float64 `gorm:"column:total"`
	}

	var results []DayResult
	currentDate := startDate
	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		dateStr := currentDate.Format("2006-01-02")
		var total float64
		r.db.WithContext(ctx).Model(&GormOrder{}).
			Where("company_id = ? AND DATE(FROM_UNIXTIME(created_at)) = ? AND deleted_at IS NULL", companyID, dateStr).
			Select("COALESCE(SUM(total_price), 0)").
			Scan(&total)
		results = append(results, DayResult{Date: dateStr, Total: total})
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	chartPoints := make([]domain.ChartPoint, len(results))
	for i, r := range results {
		parsedTime, _ := time.Parse("2006-01-02", r.Date)
		label := parsedTime.Format("02/01")
		chartPoints[i] = domain.ChartPoint{Label: label, Value: r.Total}
	}

	return chartPoints
}
