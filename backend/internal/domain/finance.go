package domain

import "time"

// TransactionType representa o tipo de transação financeira
type TransactionType string

const (
	TransactionIncome  TransactionType = "income"  // Receita
	TransactionExpense TransactionType = "expense" // Despesa
)

// TransactionCategory representa uma categoria financeira
type TransactionCategory struct {
	ID        uint            `json:"id"`
	CompanyID uint            `json:"companyId" gorm:"not null;index"`
	Name      string          `json:"name" gorm:"not null"`
	Type      TransactionType `json:"type" gorm:"not null"`
	Color     string          `json:"color"`
	Active    bool            `json:"active" gorm:"default:true"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	DeletedAt *time.Time      `json:"deletedAt,omitempty" gorm:"index"`
}

// TableName especifica o nome da tabela
func (TransactionCategory) TableName() string {
	return "transaction_categories"
}

// Transaction representa uma transação financeira
type Transaction struct {
	ID          uint            `json:"id"`
	CompanyID   uint            `json:"companyId" gorm:"not null;index"`
	CategoryID  uint            `json:"categoryId" gorm:"not null;index"`
	Type        TransactionType `json:"type" gorm:"not null"`
	Amount      Money           `json:"amount" gorm:"not null"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date" gorm:"not null"`
	Reference   string          `json:"reference"` // Referência externa (ex: pedido #123)
	CreatedBy   uint            `json:"createdBy" gorm:"not null"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	DeletedAt   *time.Time      `json:"deletedAt,omitempty" gorm:"index"`

	// Relações
	Category *TransactionCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

// TableName especifica o nome da tabela
func (Transaction) TableName() string {
	return "transactions"
}

// CashFlow representa o fluxo de caixa em um período
type CashFlow struct {
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	Income         Money     `json:"income"`
	Expense        Money     `json:"expense"`
	Balance        Money     `json:"balance"`
	OpeningBalance Money     `json:"openingBalance"`
	ClosingBalance Money     `json:"closingBalance"`
}

// FinancialSummary representa um resumo financeiro
type FinancialSummary struct {
	TotalIncome      Money `json:"totalIncome"`
	TotalExpense     Money `json:"totalExpense"`
	NetBalance       Money `json:"netBalance"`
	TransactionCount int64 `json:"transactionCount"`
}
