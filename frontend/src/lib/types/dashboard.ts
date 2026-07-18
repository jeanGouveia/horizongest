export interface DashboardMetrics {
	todayRevenue: number;
	todayOrders: number;
	pendingOrders: number;
	lowStockCount: number;
	activeProducts: number;
}

export interface RecentOrder {
	id: number;
	status: string;
	totalPrice: number;
	createdAt: string;
	itemsCount: number;
}

export interface LowStockItem {
	id: number;
	name: string;
	stockQuantity: number;
	minStock: number;
	unit: string;
}

export interface Dashboard {
	metrics: DashboardMetrics;
	recentOrders: RecentOrder[];
	lowStock: LowStockItem[];
	totalProducts: number;
	totalCategories: number;
	totalIngredients: number;
}
