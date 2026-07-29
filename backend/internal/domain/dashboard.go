package domain

// Dashboard representa os dados do dashboard executivo
type Dashboard struct {
	Metrics          DashboardMetrics `json:"metrics"`
	RecentOrders     []RecentOrder    `json:"recentOrders"`
	LowStock         []LowStockItem   `json:"lowStock"`
	ZeroStock        []LowStockItem   `json:"zeroStock"`
	TotalProducts    int              `json:"totalProducts"`
	TotalCategories  int              `json:"totalCategories"`
	TotalIngredients int              `json:"totalIngredients"`
	Charts           DashboardCharts  `json:"charts"`
}

// DashboardMetrics representa os KPIs executivos
type DashboardMetrics struct {
	// Hoje
	TodayRevenue       Money `json:"todayRevenue"`
	TodayOrders        int   `json:"todayOrders"`
	TodayProductsSold  int   `json:"todayProductsSold"`
	TodayAverageTicket Money `json:"todayAverageTicket"`
	TodayCMV           Money `json:"todayCMV"`
	TodayGrossProfit   Money `json:"todayGrossProfit"`

	// Ontem
	YesterdayRevenue       Money `json:"yesterdayRevenue"`
	YesterdayOrders        int   `json:"yesterdayOrders"`
	YesterdayProductsSold  int   `json:"yesterdayProductsSold"`
	YesterdayAverageTicket Money `json:"yesterdayAverageTicket"`

	// Semana
	WeekRevenue       Money `json:"weekRevenue"`
	WeekOrders        int   `json:"weekOrders"`
	WeekProductsSold  int   `json:"weekProductsSold"`
	WeekAverageTicket Money `json:"weekAverageTicket"`

	// Mês
	MonthRevenue       Money `json:"monthRevenue"`
	MonthOrders        int   `json:"monthOrders"`
	MonthProductsSold  int   `json:"monthProductsSold"`
	MonthAverageTicket Money `json:"monthAverageTicket"`

	// Geral
	PendingOrders   int `json:"pendingOrders"`
	CancelledOrders int `json:"cancelledOrders"`
	LowStockCount   int `json:"lowStockCount"`
	ZeroStockCount  int `json:"zeroStockCount"`
	ActiveProducts  int `json:"activeProducts"`
}

// DashboardCharts representa os dados dos gráficos
type DashboardCharts struct {
	SalesByDay    []ChartPoint `json:"salesByDay"`
	SalesByHour   []ChartPoint `json:"salesByHour"`
	TopProducts   []TopItem    `json:"topProducts"`
	TopCategories []TopItem    `json:"topCategories"`
}

// ChartPoint representa um ponto em um gráfico
type ChartPoint struct {
	Label string `json:"label"`
	Value Money  `json:"value"`
}

// TopItem representa um item em ranking
type TopItem struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Value Money  `json:"value"`
	Count int    `json:"count"`
}

// RecentOrder representa um pedido recente para o dashboard
type RecentOrder struct {
	ID          uint        `json:"id"`
	OrderNumber int         `json:"orderNumber"` // Número comercial (sequencial por empresa)
	Status      OrderStatus `json:"status"`
	TotalPrice  Money       `json:"totalPrice"`
	CreatedAt   string      `json:"createdAt"`
	ItemsCount  int         `json:"itemsCount"`
}

// LowStockItem representa um ingrediente com estoque baixo
type LowStockItem struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	StockQuantity float64 `json:"stockQuantity"`
	MinStock      float64 `json:"minStock"`
	Unit          string  `json:"unit"`
}
