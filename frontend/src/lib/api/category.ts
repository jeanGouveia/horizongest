import { request } from './client';
import type { Category, CategoryCreatePayload, CategoryUpdatePayload } from '$lib/types/category';

export async function getCategories(): Promise<Category[]> {
  const res = await request<Category[]>('/categories');
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function getCategory(id: number): Promise<Category> {
  const res = await request<Category>(`/categories/${id}`);
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function createCategory(payload: CategoryCreatePayload): Promise<Category> {
  const res = await request<Category>('/categories', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function updateCategory(id: number, payload: CategoryUpdatePayload): Promise<Category> {
  const res = await request<Category>(`/categories/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload)
  });
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function deleteCategory(id: number): Promise<void> {
  const res = await request<{ message: string }>(`/categories/${id}`, {
    method: 'DELETE'
  });
  if (res.error) throw new Error(res.error);
}
