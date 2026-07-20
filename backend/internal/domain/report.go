package domain

import "time"

// ReportType representa o tipo de relatório
type ReportType string

const (
	ReportSales      ReportType = "sales"
	ReportProducts   ReportType = "products"
	ReportCMV        ReportType = "cmv"
	ReportProfit     ReportType = "profit"
	ReportStock      ReportType = "stock"
	ReportPurchases  ReportType = "purchases"
	ReportFinancial  ReportType = "financial"
)

// ReportFormat representa o formato de exportação
type ReportFormat string

const (
	ReportFormatCSV  ReportFormat = "csv"
	ReportFormatJSON ReportFormat = "json"
	ReportFormatPDF  ReportFormat = "pdf"
)

// SalesReport representa um relatório de vendas
type SalesReport struct {
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	TotalOrders    int       `json:"totalOrders"`
	TotalRevenue   float64   `json:"totalRevenue"`
	AverageTicket  float64   `json:"averageTicket"`
	ProductsSold   int       `json:"productsSold"`
	CancelledOrders int      `json:"cancelledOrders"`
	TopProducts    []TopItem `json:"topProducts"`
	SalesByDay     []ChartPoint `json:"salesByDay"`
}

// ProductsReport representa um relatório de produtos
type ProductsReport struct {
	TotalProducts     int          `json:"totalProducts"`
	ActiveProducts    int          `json:"activeProducts"`
	ArchivedProducts  int          `json:"archivedProducts"`
	ProductsWithoutRecipe int      `json:"productsWithoutRecipe"`
	TopSellingProducts []TopItem   `json:"topSellingProducts"`
	LowMarginProducts  []ProductMarginItem `json:"lowMarginProducts"`
}

// ProductMarginItem representa um produto com margem baixa
type ProductMarginItem struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Cost   float64 `json:"cost"`
	Price  float64 `json:"price"`
	Margin float64 `json:"margin"`
}

// CMVReport representa um relatório de CMV
type CMVReport struct {
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
	TotalRevenue float64   `json:"totalRevenue"`
	TotalCMV     float64   `json:"totalCMV"`
	CMVPercentage float64  `json:"cmvPercentage"`
	GrossProfit  float64   `json:"grossProfit"`
	ProfitMargin float64   `json:"profitMargin"`
	ByProduct    []ProductCMVItem `json:"byProduct"`
}

// ProductCMVItem representa o CMV de um produto
type ProductCMVItem struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Revenue     float64 `json:"revenue"`
	CMV         float64 `json:"cmv"`
	CMVPercent  float64 `json:"cmvPercent"`
	GrossProfit float64 `json:"grossProfit"`
}

// ProfitReport representa um relatório de lucro
type ProfitReport struct {
	StartDate     time.Time `json:"startDate"`
	EndDate       time.Time `json:"endDate"`
	TotalRevenue  float64   `json:"totalRevenue"`
	TotalExpense  float64   `json:"totalExpense"`
	NetProfit    float64   `json:"netProfit"`
	ProfitMargin float64   `json:"profitMargin"`
	ByCategory   []CategoryProfitItem `json:"byCategory"`
}

// CategoryProfitItem representa o lucro por categoria
type CategoryProfitItem struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Revenue     float64 `json:"revenue"`
	CMV         float64 `json:"cmv"`
	GrossProfit float64 `json:"grossProfit"`
}

// StockReport representa um relatório de estoque
type StockReport struct {
	TotalIngredients    int               `json:"totalIngredients"`
	LowStockCount      int               `json:"lowStockCount"`
	ZeroStockCount     int               `json:"zeroStockCount"`
	TotalStockValue    float64           `json:"totalStockValue"`
	LowStockItems      []LowStockItem    `json:"lowStockItems"`
	ZeroStockItems     []LowStockItem    `json:"zeroStockItems"`
	HighValueItems     []StockValueItem  `json:"highValueItems"`
}

// StockValueItem representa um item com alto valor de estoque
type StockValueItem struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	StockQuantity float64 `json:"stockQuantity"`
	Unit         string  `json:"unit"`
	UnitCost     float64 `json:"unitCost"`
	TotalValue   float64 `json:"totalValue"`
}

// PurchasesReport representa um relatório de compras
type PurchasesReport struct {
	StartDate       time.Time `json:"startDate"`
	EndDate         time.Time `json:"endDate"`
	TotalOrders    int       `json:"totalOrders"`
	TotalAmount    float64   `json:"totalAmount"`
	PendingOrders  int       `json:"pendingOrders"`
	ReceivedOrders int       `json:"receivedOrders"`
	TopSuppliers   []TopItem `json:"topSuppliers"`
	PurchasesByDay []ChartPoint `json:"purchasesByDay"`
}

// FinancialReport representa um relatório financeiro
type FinancialReport struct {
	StartDate      time.Time `json:"startDate"`
	EndDate        time.Time `json:"endDate"`
	TotalIncome   float64 `json:"totalIncome"`
	TotalExpense  float64 `json:"totalExpense"`
	NetBalance    float64 `json:"netBalance"`
	ByCategory    []CategoryFinancialItem `json:"byCategory"`
	CashFlow      *CashFlow `json:"cashFlow"`
}

// CategoryFinancialItem representa valores por categoria financeira
type CategoryFinancialItem struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
}
