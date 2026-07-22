package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

var (
	ErrCategoryNotFound = errors.New("categoria não encontrada")
)

type CategoryService struct {
	repo ports.CategoryRepository
}

func NewCategoryService(repo ports.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

type CreateCategoryInput struct {
	Name         string `json:"name" validate:"required,min=2,max=120"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"display_order" validate:"gte=0"`
}

type UpdateCategoryInput struct {
	Name         string `json:"name" validate:"required,min=2,max=120"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"display_order" validate:"gte=0"`
	Active       *bool  `json:"active"`
}

func (s *CategoryService) CreateCategory(ctx context.Context, in CreateCategoryInput) (*domain.Category, error) {
	c := &domain.Category{
		Name:         in.Name,
		Description:  in.Description,
		DisplayOrder: in.DisplayOrder,
		Active:       true,
	}
	if err := s.repo.CreateCategory(ctx, c); err != nil {
		return nil, fmt.Errorf("CategoryService.CreateCategory: %w", err)
	}
	return c, nil
}

func (s *CategoryService) ListCategories(ctx context.Context) ([]domain.Category, error) {
	categories, err := s.repo.ListCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("CategoryService.ListCategories: %w", err)
	}
	return categories, nil
}

func (s *CategoryService) GetCategory(ctx context.Context, id uint) (*domain.Category, error) {
	c, err := s.repo.FindCategoryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("CategoryService.GetCategory: %w", err)
	}
	if c == nil {
		return nil, ErrCategoryNotFound
	}
	return c, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id uint, in UpdateCategoryInput) (*domain.Category, error) {
	c, err := s.repo.FindCategoryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("CategoryService.UpdateCategory: %w", err)
	}
	if c == nil {
		return nil, ErrCategoryNotFound
	}

	c.Name = in.Name
	c.Description = in.Description
	c.DisplayOrder = in.DisplayOrder
	if in.Active != nil {
		c.Active = *in.Active
	}

	if err := s.repo.UpdateCategory(ctx, c); err != nil {
		return nil, fmt.Errorf("CategoryService.UpdateCategory: %w", err)
	}
	return c, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id uint) error {
	c, err := s.repo.FindCategoryByID(ctx, id)
	if err != nil {
		return fmt.Errorf("CategoryService.DeleteCategory: %w", err)
	}
	if c == nil {
		return ErrCategoryNotFound
	}
	return s.repo.DeleteCategory(ctx, id)
}
