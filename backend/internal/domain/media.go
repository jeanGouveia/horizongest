package domain

import "time"

type MediaType string

const (
	MediaTypeImage    MediaType = "image"
	MediaTypeVideo    MediaType = "video"
	MediaTypeDocument MediaType = "document"
)

type Media struct {
	ID            uint
	FileName      string
	OriginalName  string
	FilePath      string
	ThumbnailPath string
	FileSize      int64
	MimeType      string
	Width         *int
	Height        *int
	AltText       string
	EntityType    string // "product", "category", etc.
	EntityID      *uint
	CompanyID     uint // ID da empresa/tenant (obrigatório - Sprint 3)
	DeletedAt     *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
