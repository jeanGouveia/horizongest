export interface Media {
  ID: number;
  FileName: string;
  OriginalName: string;
  FilePath: string;
  ThumbnailPath: string;
  FileSize: number;
  MimeType: string;
  Width?: number;
  Height?: number;
  AltText: string;
  EntityType: string;
  EntityID?: number;
  CreatedAt?: string;
  UpdatedAt?: string;
}

export interface MediaUploadResponse {
  ID: number;
  FileName: string;
  OriginalName: string;
  FilePath: string;
  ThumbnailPath: string;
  FileSize: number;
  MimeType: string;
  Width?: number;
  Height?: number;
  AltText: string;
  EntityType: string;
  EntityID?: number;
  CreatedAt: string;
  UpdatedAt: string;
}
