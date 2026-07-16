import { request } from './client';
import type { Product, ProductCreatePayload, ProductUpdatePayload, ProductIngredientPayload } from '$lib/types/product';
import type { Ingredient, IngredientCreatePayload } from '$lib/types/ingredient';

// ── Produtos ────────────────────────────────────────────────────────────────

export async function getProducts(): Promise<Product[]> {
  const res = await request<Product[]>('/products');
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function getActiveProducts(): Promise<Product[]> {
  const res = await request<Product[]>('/products/active');
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function getProduct(id: number): Promise<Product> {
  const res = await request<Product>(`/products/${id}`);
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function createProduct(payload: ProductCreatePayload): Promise<Product> {
  const res = await request<Product>('/products', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function updateProduct(id: number, payload: ProductUpdatePayload): Promise<Product> {
  const res = await request<Product>(`/products/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload)
  });
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function deleteProduct(id: number): Promise<void> {
  const res = await request<{ message: string }>(`/products/${id}`, {
    method: 'DELETE'
  });
  if (res.error) throw new Error(res.error);
}

export async function getProductIngredients(id: number): Promise<Ingredient[]> {
  const res = await request<Ingredient[]>(`/products/${id}/ingredients`);
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function updateProductIngredients(
  id: number,
  ingredients: ProductIngredientPayload[]
): Promise<void> {
  const res = await request<{ message: string }>(`/products/${id}/ingredients`, {
    method: 'PUT',
    body: JSON.stringify({ items: ingredients })
  });
  if (res.error) throw new Error(res.error);
}

// ── Ingredientes ─────────────────────────────────────────────────────────────

export async function getIngredients(): Promise<Ingredient[]> {
  const res = await request<Ingredient[]>('/ingredients');
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function createIngredient(payload: IngredientCreatePayload): Promise<Ingredient> {
  const res = await request<Ingredient>('/ingredients', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function getIngredient(id: number): Promise<Ingredient> {
  const res = await request<Ingredient>(`/ingredients/${id}`);
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function updateIngredient(id: number, payload: IngredientCreatePayload): Promise<Ingredient> {
  const res = await request<Ingredient>(`/ingredients/${id}`, {
    method: 'PUT',
    body: JSON.stringify(payload)
  });
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function deleteIngredient(id: number): Promise<void> {
  const res = await request<{ message: string }>(`/ingredients/${id}`, {
    method: 'DELETE'
  });
  if (res.error) throw new Error(res.error);
}

export async function updateIngredientStock(id: number, quantity: number): Promise<Ingredient> {
  const res = await request<Ingredient>(`/ingredients/${id}/stock`, {
    method: 'PATCH',
    body: JSON.stringify({ quantity })
  });
  if (res.error) throw new Error(res.error);
  return res.data!;
}
