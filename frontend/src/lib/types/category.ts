export interface Category {
  ID: number;
  Name: string;
  Description?: string;
  DisplayOrder: number;
  Active: boolean;
  CreatedAt?: string;
  UpdatedAt?: string;
}

export interface CategoryCreatePayload {
  name: string;
  description?: string;
  display_order: number;
}

export interface CategoryUpdatePayload {
  name: string;
  description?: string;
  display_order: number;
  active: boolean;
}
