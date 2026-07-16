import { request } from './client';
import type { Media, MediaUploadResponse } from '$lib/types/media';

export async function uploadMedia(file: File, entityType: string = 'product', entityId?: number, altText?: string): Promise<MediaUploadResponse> {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('entity_type', entityType);
  if (entityId) {
    formData.append('entity_id', entityId.toString());
  }
  if (altText) {
    formData.append('alt_text', altText);
  }

  const res = await request<MediaUploadResponse>('/media/upload', {
    method: 'POST',
    body: formData,
    headers: {} // Remove Content-Type para permitir boundary automático
  });

  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function getMedia(id: number): Promise<Media> {
  const res = await request<Media>(`/media/${id}`);
  if (res.error) throw new Error(res.error);
  return res.data!;
}

export async function deleteMedia(id: number): Promise<void> {
  const res = await request<{ message: string }>(`/media/${id}`, {
    method: 'DELETE'
  });
  if (res.error) throw new Error(res.error);
}

export async function getMediaByEntity(entityType: string, entityId: number): Promise<Media[]> {
  const res = await request<Media[]>(`/media/entity/${entityType}/${entityId}`);
  if (res.error) throw new Error(res.error);
  return res.data!;
}
