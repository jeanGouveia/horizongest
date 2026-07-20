package ports

import (
	"context"
	"time"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
)

// ReportRepository define a interface para repositório de relatórios
type ReportRepository interface {
	// Relatório de Vendas
	GetSalesReport(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.SalesReport, error)
	
	// Relatório de Produtos
	GetProductsReport(ctx context.Context, companyID uint) (*domain.ProductsReport, error)
	
	// Relatório de CMV
	GetCMVReport(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.CMVReport, error)
	
	// Relatório de Lucro
	GetProfitReport(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.ProfitReport, error)
	
	// Relatório de Estoque
	GetStockReport(ctx context.Context, companyID uint) (*domain.StockReport, error)
	
	// Relatório de Compras
	GetPurchasesReport(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.PurchasesReport, error)
	
	// Relatório Financeiro
	GetFinancialReport(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.FinancialReport, error)
}
