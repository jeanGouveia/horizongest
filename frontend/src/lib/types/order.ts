export type OrderStatus = 'pending' | 'confirmed' | 'preparing' | 'ready' | 'delivered' | 'cancelled';

export interface OrderItem {
  ID: number;
  OrderID: number;
  ProductID: number;
  Quantity: number;
  UnitPrice: number;
  Product?: any; // pode ser null ou objeto completo
}

export interface Order {
  ID: number;
  OrderNumber: number; // Número comercial do pedido (sequencial por empresa)
  Status: OrderStatus;
  TotalPrice: number;
  Notes?: string;
  Items: OrderItem[];
  TableNumber?: number;
  CreatedAt?: string;
  UpdatedAt?: string;
}

export type OrderCreatePayload = {
  notes?: string;
  table_number?: number;
  items: { product_id: number; quantity: number }[];
};

export const ORDER_STATUS_LABEL: Record<OrderStatus, string> = {
  pending:    'Pendente',
  confirmed:  'Confirmado',
  preparing:  'Preparando',
  ready:      'Pronto',
  delivered:  'Entregue',
  cancelled:  'Cancelado',
};

export const ORDER_STATUS_COLOR: Record<OrderStatus, string> = {
  pending:   'badge-warning',
  confirmed: 'badge-info',
  preparing: 'badge-info',
  ready:     'badge-success',
  delivered: 'badge-neutral',
  cancelled: 'badge-error',
};
