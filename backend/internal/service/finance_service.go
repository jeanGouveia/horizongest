package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

var (
	ErrTransactionCategoryNotFound = errors.New("categoria não encontrada")
	ErrTransactionNotFound         = errors.New("transação não encontrada")
	ErrInvalidAmount               = errors.New("valor inválido")
	ErrInvalidDate                 = errors.New("data inválida")
)

type FinanceService struct {
	financeRepo ports.FinanceRepository
}

func NewFinanceService(financeRepo ports.FinanceRepository) *FinanceService {
	return &FinanceService{financeRepo: financeRepo}
}

// --- Categorias ---

// CreateTransactionCategory cria uma nova categoria financeira
func (s *FinanceService) CreateTransactionCategory(ctx context.Context, companyID uint, input CreateTransactionCategoryInput) (*domain.TransactionCategory, error) {
	category := &domain.TransactionCategory{
		CompanyID: companyID,
		Name:      input.Name,
		Type:      input.Type,
		Color:     input.Color,
		Active:    true,
	}

	if err := s.financeRepo.CreateTransactionCategory(ctx, category); err != nil {
		return nil, fmt.Errorf("FinanceService.CreateTransactionCategory: criar categoria: %w", err)
	}

	return category, nil
}

// ListTransactionCategories lista categorias financeiras
func (s *FinanceService) ListTransactionCategories(ctx context.Context, companyID uint, transactionType *domain.TransactionType, limit, offset int) ([]domain.TransactionCategory, error) {
	return s.financeRepo.ListTransactionCategories(ctx, companyID, transactionType, limit, offset)
}

// GetTransactionCategoryByID busca uma categoria por ID
func (s *FinanceService) GetTransactionCategoryByID(ctx context.Context, id uint) (*domain.TransactionCategory, error) {
	return s.financeRepo.GetTransactionCategoryByID(ctx, id)
}

// UpdateTransactionCategory atualiza uma categoria
func (s *FinanceService) UpdateTransactionCategory(ctx context.Context, category *domain.TransactionCategory) error {
	return s.financeRepo.UpdateTransactionCategory(ctx, category)
}

// DeleteTransactionCategory deleta uma categoria (soft delete)
func (s *FinanceService) DeleteTransactionCategory(ctx context.Context, id uint) error {
	return s.financeRepo.DeleteTransactionCategory(ctx, id)
}

// --- Transações ---

// CreateTransaction cria uma nova transação financeira
func (s *FinanceService) CreateTransaction(ctx context.Context, companyID, userID uint, input CreateTransactionInput) (*domain.Transaction, error) {
	// Validar valor
	if input.Amount <= 0 {
		return nil, ErrInvalidAmount
	}

	// Validar data
	if input.Date.IsZero() {
		input.Date = time.Now()
	}

	// Validar categoria
	category, err := s.financeRepo.GetTransactionCategoryByID(ctx, input.CategoryID)
	if err != nil {
		return nil, ErrTransactionCategoryNotFound
	}

	if category.CompanyID != companyID {
		return nil, errors.New("categoria não pertence a esta empresa")
	}

	// Validar tipo da categoria
	if category.Type != input.Type {
		return nil, errors.New("tipo de transação incompatível com categoria")
	}

	transaction := &domain.Transaction{
		CompanyID:   companyID,
		CategoryID:  input.CategoryID,
		Type:        input.Type,
		Amount:      domain.FromFloat64(input.Amount),
		Description: input.Description,
		Date:        input.Date,
		Reference:   input.Reference,
		CreatedBy:   userID,
	}

	if err := s.financeRepo.CreateTransaction(ctx, transaction); err != nil {
		return nil, fmt.Errorf("FinanceService.CreateTransaction: criar transação: %w", err)
	}

	return transaction, nil
}

// ListTransactions lista transações financeiras
func (s *FinanceService) ListTransactions(ctx context.Context, companyID uint, transactionType *domain.TransactionType, startDate, endDate *time.Time, limit, offset int) ([]domain.Transaction, error) {
	return s.financeRepo.ListTransactions(ctx, companyID, transactionType, startDate, endDate, limit, offset)
}

// GetTransactionByID busca uma transação por ID
func (s *FinanceService) GetTransactionByID(ctx context.Context, id uint) (*domain.Transaction, error) {
	return s.financeRepo.GetTransactionByID(ctx, id)
}

// UpdateTransaction atualiza uma transação
func (s *FinanceService) UpdateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	return s.financeRepo.UpdateTransaction(ctx, transaction)
}

// DeleteTransaction deleta uma transação (soft delete)
func (s *FinanceService) DeleteTransaction(ctx context.Context, id uint) error {
	return s.financeRepo.DeleteTransaction(ctx, id)
}

// --- Resumos ---

// GetCashFlow retorna o fluxo de caixa em um período
func (s *FinanceService) GetCashFlow(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.CashFlow, error) {
	return s.financeRepo.GetCashFlow(ctx, companyID, startDate, endDate)
}

// GetFinancialSummary retorna um resumo financeiro
func (s *FinanceService) GetFinancialSummary(ctx context.Context, companyID uint, startDate, endDate time.Time) (*domain.FinancialSummary, error) {
	return s.financeRepo.GetFinancialSummary(ctx, companyID, startDate, endDate)
}

// --- Inputs ---

type CreateTransactionCategoryInput struct {
	Name  string                 `json:"name" validate:"required"`
	Type  domain.TransactionType `json:"type" validate:"required"`
	Color string                 `json:"color"`
}

type CreateTransactionInput struct {
	CategoryID  uint                   `json:"categoryId" validate:"required"`
	Type        domain.TransactionType `json:"type" validate:"required"`
	Amount      float64                `json:"amount" validate:"required,gt=0"`
	Description string                 `json:"description"`
	Date        time.Time              `json:"date"`
	Reference   string                 `json:"reference"`
}
