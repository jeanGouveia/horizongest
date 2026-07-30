package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

var _ ports.FinanceRepository = (*GormFinanceRepository)(nil)

type GormFinanceRepository struct {
	db *gorm.DB
}

func NewGormFinanceRepository(db *gorm.DB) *GormFinanceRepository {
	return &GormFinanceRepository{db: db}
}

// GORM models
type GormTransactionCategory struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	CompanyID uint   `gorm:"not null;index"`
	Name      string `gorm:"not null"`
	Type      string `gorm:"not null"`
	Color     string
	Active    bool       `gorm:"default:true"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"`
}

func (GormTransactionCategory) TableName() string { return "transaction_categories" }

type GormTransaction struct {
	ID          uint    `gorm:"primaryKey;autoIncrement"`
	CompanyID   uint    `gorm:"not null;index"`
	CategoryID  uint    `gorm:"not null;index"`
	Type        string  `gorm:"not null"`
	Amount      float64 `gorm:"not null"`
	Description string
	Date        time.Time `gorm:"not null"`
	Reference   string
	CreatedBy   uint       `gorm:"not null"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `gorm:"index"`
}

func (GormTransaction) TableName() string { return "transactions" }

// --- Categorias ---

func (r *GormFinanceRepository) CreateTransactionCategory(ctx context.Context, category *domain.TransactionCategory) error {
	companyID, err := GetCompanyIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("CreateTransactionCategory: %w", err)
	}
	category.CompanyID = companyID
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *GormFinanceRepository) ListTransactionCategories(ctx context.Context, companyID uint, transactionType *domain.TransactionType, limit, offset int) ([]domain.TransactionCategory, error) {
	var categories []domain.TransactionCategory
	query := r.db.WithContext(ctx).Where("company_id = ? AND deleted_at IS NULL", companyID)

	if transactionType != nil {
		query = query.Where("type = ?", *transactionType)
	}

	query = query.Order("name ASC").Limit(limit).Offset(offset)

	err := query.Find(&categories).Error
	return categories, fmt.Errorf("ListTransactionCategories: %w", err)
}

func (r *GormFinanceRepository) GetTransactionCategoryByID(ctx context.Context, id uint) (*domain.TransactionCategory, error) {
	var category domain.TransactionCategory
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	err := query.Where("deleted_at IS NULL").First(&category).Error
	return &category, fmt.Errorf("GetTransactionCategoryByID: %w", err)
}

func (r *GormFinanceRepository) UpdateTransactionCategory(ctx context.Context, category *domain.TransactionCategory) error {
	query := ApplyTenantFilterWithID(ctx, r.db, category.ID)
	return query.WithContext(ctx).Save(category).Error
}

func (r *GormFinanceRepository) DeleteTransactionCategory(ctx context.Context, id uint) error {
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	return query.Delete(&domain.TransactionCategory{}).Error
}

// --- Transações ---

func (r *GormFinanceRepository) CreateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	companyID, err := GetCompanyIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("CreateTransaction: %w", err)
	}
	transaction.CompanyID = companyID
	return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *GormFinanceRepository) ListTransactions(ctx context.Context, companyID uint, transactionType *domain.TransactionType, startDate, endDate *time.Time, limit, offset int) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	query := r.db.WithContext(ctx).Where("company_id = ? AND deleted_at IS NULL", companyID)

	if transactionType != nil {
		query = query.Where("type = ?", *transactionType)
	}

	if startDate != nil {
		query = query.Where("date >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("date <= ?", *endDate)
	}

	query = query.Order("date DESC").Limit(limit).Offset(offset)

	err := query.Preload("Category").Find(&transactions).Error
	return transactions, fmt.Errorf("ListTransactions: %w", err)
}

func (r *GormFinanceRepository) GetTransactionByID(ctx context.Context, id uint) (*domain.Transaction, error) {
	var transaction domain.Transaction
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	err := query.Where("deleted_at IS NULL").
		Preload("Category").
		First(&transaction).Error
	return &transaction, fmt.Errorf("GetTransactionByID: %w", err)
}

func (r *GormFinanceRepository) UpdateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	query := ApplyTenantFilterWithID(ctx, r.db, transaction.ID)
	return query.WithContext(ctx).Save(transaction).Error
}

func (r *GormFinanceRepository) DeleteTransaction(ctx context.Context, id uint) error {
	query := ApplyTenantFilterWithID(ctx, r.db, id)
	return query.Delete(&domain.Transaction{}).Error
}

// --- Resumos ---

func (r *GormFinanceRepository) GetCashFlow(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.CashFlow, error) {
	cashFlow := &domain.CashFlow{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Calcular saldo de abertura (transações antes da data inicial)
	var openingBalance int64
	r.db.WithContext(ctx).Model(&domain.Transaction{}).
		Where("company_id = ? AND date < ? AND deleted_at IS NULL", companyID, startDate).
		Select("COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE -amount END)), 0)").
		Scan(&openingBalance)
	cashFlow.OpeningBalance = domain.Money(openingBalance)

	// Calcular receitas no período
	var income int64
	r.db.WithContext(ctx).Model(&domain.Transaction{}).
		Where("company_id = ? AND type = 'income' AND date >= ? AND date <= ? AND deleted_at IS NULL", companyID, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&income)
	cashFlow.Income = domain.Money(income)

	// Calcular despesas no período
	var expense int64
	r.db.WithContext(ctx).Model(&domain.Transaction{}).
		Where("company_id = ? AND type = 'expense' AND date >= ? AND date <= ? AND deleted_at IS NULL", companyID, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&expense)
	cashFlow.Expense = domain.Money(expense)

	// Calcular saldo no período
	cashFlow.Balance = cashFlow.Income.Sub(cashFlow.Expense)

	// Calcular saldo de fechamento
	cashFlow.ClosingBalance = cashFlow.OpeningBalance.Add(cashFlow.Balance)

	return cashFlow, nil
}

func (r *GormFinanceRepository) GetFinancialSummary(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.FinancialSummary, error) {
	summary := &domain.FinancialSummary{}

	// Calcular receitas
	r.db.WithContext(ctx).Model(&domain.Transaction{}).
		Where("company_id = ? AND type = 'income' AND date >= ? AND date <= ? AND deleted_at IS NULL", companyID, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&summary.TotalIncome)

	// Calcular despesas
	r.db.WithContext(ctx).Model(&domain.Transaction{}).
		Where("company_id = ? AND type = 'expense' AND date >= ? AND date <= ? AND deleted_at IS NULL", companyID, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&summary.TotalExpense)

	// Calcular saldo líquido
	summary.NetBalance = summary.TotalIncome - summary.TotalExpense

	// Contar transações
	r.db.WithContext(ctx).Model(&domain.Transaction{}).
		Where("company_id = ? AND date >= ? AND date <= ? AND deleted_at IS NULL", companyID, startDate, endDate).
		Count(&summary.TransactionCount)

	return summary, nil
}
