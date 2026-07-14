import { request } from './client';
import type { StockAdjustment, StockAdjustmentActionPayload } from '$lib/types/stock-adjustment';

export async function getPendingAdjustments(filters?: {
  status?: string;
  order_id?: number;
  ingredient_id?: number;
}): Promise<StockAdjustment[]> {
  const params = new URLSearchParams();
  if (filters?.status) params.set('status', filters.status);
  if (filters?.order_id) params.set('order_id', filters.order_id.toString());
  if (filters?.ingredient_id) params.set('ingredient_id', filters.ingredient_id.toString());

  const res = await request<StockAdjustment[]>(`/stock-adjustments/pending?${params.toString()}`);
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function approveAdjustment(id: number, payload: StockAdjustmentActionPayload): Promise<StockAdjustment> {
  const res = await request<StockAdjustment>(`/stock-adjustments/${id}/approve`, {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function rejectAdjustment(id: number, payload: StockAdjustmentActionPayload): Promise<StockAdjustment> {
  const res = await request<StockAdjustment>(`/stock-adjustments/${id}/reject`, {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  if (res.error) throw new Error(res.error);
  return res.data!;
}
