export interface Ingredient {
  ID: number;
  Name: string;
  Unit: string;
  Quantity?: number; // quantidade usada no produto
}

export interface Product {
  ID: number;
  Name: string;
  Description?: string;
  Price: number;
  IsComposto: boolean;
  Active: boolean;
  PhotoURL?: string;
  CategoryID?: number;
  DisplayOrder: number;
  PreparationTimeMinutes: number;
  Featured: boolean;
  IsNew: boolean;
  PromotionPrice?: number;
  PromotionStart?: string;
  PromotionEnd?: string;
  AvailableFrom?: string;
  AvailableUntil?: string;
  SKU?: string;
  InternalNotes?: string;
  Slug?: string;
  MetaTitle?: string;
  MetaDescription?: string;
  AltImage?: string;
  Canonical?: string;
  ExternalID?: string;
  MarketplaceID?: string;
  SyncStatus?: string;
  LastSync?: string;
  CreatedAt?: string;
  UpdatedAt?: string;
  Ingredients?: Ingredient[];
}

export interface ProductCreatePayload {
  name: string;
  description?: string;
  price: number;
  is_composto: boolean;
  active: boolean;
  photo_url?: string;
  category_id?: number;
  display_order: number;
  preparation_time_minutes: number;
  featured: boolean;
  is_new: boolean;
  promotion_price?: number;
  promotion_start?: string;
  promotion_end?: string;
  available_from?: string;
  available_until?: string;
  sku?: string;
  internal_notes?: string;
  slug?: string;
  meta_title?: string;
  meta_description?: string;
  alt_image?: string;
  canonical?: string;
  external_id?: string;
  marketplace_id?: string;
  sync_status?: string;
}

export interface ProductUpdatePayload {
  name: string;
  description?: string;
  price: number;
  is_composto: boolean;
  active: boolean;
  photo_url?: string;
  category_id?: number;
  display_order: number;
  preparation_time_minutes: number;
  featured: boolean;
  is_new: boolean;
  promotion_price?: number;
  promotion_start?: string;
  promotion_end?: string;
  available_from?: string;
  available_until?: string;
  sku?: string;
  internal_notes?: string;
  slug?: string;
  meta_title?: string;
  meta_description?: string;
  alt_image?: string;
  canonical?: string;
  external_id?: string;
  marketplace_id?: string;
  sync_status?: string;
}

export interface ProductIngredientPayload {
  ingredient_id: number;
  quantity: number;
}
