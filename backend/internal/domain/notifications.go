package domain

// Notifications representa as notificações do sistema
type Notifications struct {
	PendingOrders        int `json:"pendingOrders"`
	LowStockCount        int `json:"lowStockCount"`
	ProductsWithoutPhoto int `json:"productsWithoutPhoto"`
	ExpiredPromotions    int `json:"expiredPromotions"`
}
