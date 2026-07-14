export type StockAdjustmentStatus = 'pending' | 'approved' | 'rejected';

export interface StockAdjustment {
  id: number;
  order_id: number;
  ingredient_id: number;
  quantity: number;
  order_status: string;
  status: StockAdjustmentStatus;
  created_at: string;
  processed_at?: string;
  processed_by?: number;
  processing_notes?: string;
}

export type StockAdjustmentActionPayload = {
  notes?: string;
};

export const STOCK_ADJUSTMENT_STATUS_LABEL: Record<StockAdjustmentStatus, string> = {
  pending: 'Pendente',
  approved: 'Aprovado',
  rejected: 'Rejeitado'
};

export const STOCK_ADJUSTMENT_STATUS_COLOR: Record<StockAdjustmentStatus, string> = {
  pending: 'badge-warning',
  approved: 'badge-success',
  rejected: 'badge-error'
};
