import { request } from './client';
import type { Order, OrderCreatePayload, OrderStatus } from '$lib/types/order';

export async function getOrders(): Promise<Order[]> {
  const res = await request<Order[]>('/orders');
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function getOrder(id: number): Promise<Order> {
  const res = await request<Order>(`/orders/${id}`);
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function createOrder(payload: OrderCreatePayload): Promise<Order> {
  const res = await request<Order>('/orders', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function updateOrderStatus(id: number, status: OrderStatus): Promise<Order> {
  const res = await request<Order>(`/orders/${id}/status`, {
    method: 'PATCH',
    body: JSON.stringify({ status })
  });
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function updateOrder(id: number, payload: OrderCreatePayload): Promise<Order> {
  const res = await request<Order>(`/orders/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload)
  });
  if (res.error) throw new Error(res.error);
  return res.data!;
}
