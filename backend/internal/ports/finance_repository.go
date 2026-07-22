package ports

import (
	"context"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
)

// FinanceRepository define a interface para repositório financeiro
type FinanceRepository interface {
	// Categorias
	CreateTransactionCategory(ctx context.Context, category *domain.TransactionCategory) error
	ListTransactionCategories(ctx context.Context, companyID uint, transactionType *domain.TransactionType, limit, offset int) ([]domain.TransactionCategory, error)
	GetTransactionCategoryByID(ctx context.Context, id uint) (*domain.TransactionCategory, error)
	UpdateTransactionCategory(ctx context.Context, category *domain.TransactionCategory) error
	DeleteTransactionCategory(ctx context.Context, id uint) error

	// Transações
	CreateTransaction(ctx context.Context, transaction *domain.Transaction) error
	ListTransactions(ctx context.Context, companyID uint, transactionType *domain.TransactionType, startDate, endDate *time.Time, limit, offset int) ([]domain.Transaction, error)
	GetTransactionByID(ctx context.Context, id uint) (*domain.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *domain.Transaction) error
	DeleteTransaction(ctx context.Context, id uint) error

	// Resumos
	GetCashFlow(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.CashFlow, error)
	GetFinancialSummary(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.FinancialSummary, error)
}
