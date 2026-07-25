export interface DashboardMetrics {
  todayRevenue: number;
  todayOrders: number;
  todayProductsSold: number;
  todayAverageTicket: number;
  todayCMV: number;
  todayGrossProfit: number;
  yesterdayRevenue: number;
  yesterdayOrders: number;
  yesterdayProductsSold: number;
  yesterdayAverageTicket: number;
  weekRevenue: number;
  weekOrders: number;
  weekProductsSold: number;
  weekAverageTicket: number;
  monthRevenue: number;
  monthOrders: number;
  monthProductsSold: number;
  monthAverageTicket: number;
  pendingOrders: number;
  cancelledOrders: number;
  lowStockCount: number;
  zeroStockCount: number;
  activeProducts: number;
}

export interface RecentOrder {
  id: number;
  orderNumber: number; // Número comercial (sequencial por empresa)
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
  zeroStock: LowStockItem[];
  totalProducts: number;
  totalCategories: number;
  totalIngredients: number;
  charts: {
    salesByDay: { label: string; value: number }[];
    salesByHour: { label: string; value: number }[];
    topProducts: { id: number; name: string; value: number; count: number }[];
    topCategories: { id: number; name: string; value: number; count: number }[];
  };
}
