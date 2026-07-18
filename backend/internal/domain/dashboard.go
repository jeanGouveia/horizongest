package domain

// Dashboard representa os dados do dashboard executivo
type Dashboard struct {
	Metrics          DashboardMetrics `json:"metrics"`
	RecentOrders     []RecentOrder    `json:"recentOrders"`
	LowStock         []LowStockItem   `json:"lowStock"`
	TotalProducts    int              `json:"totalProducts"`
	TotalCategories  int              `json:"totalCategories"`
	TotalIngredients int              `json:"totalIngredients"`
}

// DashboardMetrics representa os KPIs executivos
type DashboardMetrics struct {
	TodayRevenue   float64 `json:"todayRevenue"`
	TodayOrders    int     `json:"todayOrders"`
	PendingOrders  int     `json:"pendingOrders"`
	LowStockCount  int     `json:"lowStockCount"`
	ActiveProducts int     `json:"activeProducts"`
}

// RecentOrder representa um pedido recente para o dashboard
type RecentOrder struct {
	ID         uint        `json:"id"`
	Status     OrderStatus `json:"status"`
	TotalPrice float64     `json:"totalPrice"`
	CreatedAt  string      `json:"createdAt"`
	ItemsCount int         `json:"itemsCount"`
}

// LowStockItem representa um ingrediente com estoque baixo
type LowStockItem struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	StockQuantity float64 `json:"stockQuantity"`
	MinStock      float64 `json:"minStock"`
	Unit          string  `json:"unit"`
}
