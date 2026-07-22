package service

import (
	"context"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

type ReportService struct {
	reportRepo ports.ReportRepository
}

func NewReportService(reportRepo ports.ReportRepository) *ReportService {
	return &ReportService{reportRepo: reportRepo}
}

// GetSalesReport retorna relatório de vendas
func (s *ReportService) GetSalesReport(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.SalesReport, error) {
	return s.reportRepo.GetSalesReport(ctx, companyID, startDate, endDate)
}

// GetProductsReport retorna relatório de produtos
func (s *ReportService) GetProductsReport(ctx context.Context, companyID uint) (*domain.ProductsReport, error) {
	return s.reportRepo.GetProductsReport(ctx, companyID)
}

// GetCMVReport retorna relatório de CMV
func (s *ReportService) GetCMVReport(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.CMVReport, error) {
	return s.reportRepo.GetCMVReport(ctx, companyID, startDate, endDate)
}

// GetProfitReport retorna relatório de lucro
func (s *ReportService) GetProfitReport(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.ProfitReport, error) {
	return s.reportRepo.GetProfitReport(ctx, companyID, startDate, endDate)
}

// GetStockReport retorna relatório de estoque
func (s *ReportService) GetStockReport(ctx context.Context, companyID uint) (*domain.StockReport, error) {
	return s.reportRepo.GetStockReport(ctx, companyID)
}

// GetPurchasesReport retorna relatório de compras
func (s *ReportService) GetPurchasesReport(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.PurchasesReport, error) {
	return s.reportRepo.GetPurchasesReport(ctx, companyID, startDate, endDate)
}

// GetFinancialReport retorna relatório financeiro
func (s *ReportService) GetFinancialReport(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.FinancialReport, error) {
	return s.reportRepo.GetFinancialReport(ctx, companyID, startDate, endDate)
}
