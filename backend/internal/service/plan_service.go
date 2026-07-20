package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
)

var (
	ErrPlanNotFound = errors.New("plan not found")
	ErrPlanAlreadyExists = errors.New("plan with this slug already exists")
)

type PlanService struct {
	planRepo ports.PlanRepository
}

func NewPlanService(planRepo ports.PlanRepository) *PlanService {
	return &PlanService{
		planRepo: planRepo,
	}
}

type CreatePlanInput struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Slug        string  `json:"slug" validate:"required,min=2,max=50"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"omitempty,min=0"`
	Currency    string  `json:"currency" validate:"omitempty,len=3"`
	Interval    string  `json:"interval" validate:"omitempty,oneof=monthly yearly"`
	MaxUsers    int     `json:"max_users" validate:"omitempty,min=1"`
	MaxProducts int     `json:"max_products" validate:"omitempty,min=1"`
	Features    string  `json:"features"`
}

type UpdatePlanInput struct {
	Name        *string  `json:"name" validate:"omitempty,min=2,max=100"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price" validate:"omitempty,min=0"`
	Currency    *string  `json:"currency" validate:"omitempty,len=3"`
	Interval    *string  `json:"interval" validate:"omitempty,oneof=monthly yearly"`
	MaxUsers    *int     `json:"max_users" validate:"omitempty,min=1"`
	MaxProducts *int     `json:"max_products" validate:"omitempty,min=1"`
	Features    *string  `json:"features"`
	Active      *bool    `json:"active"`
}

func (s *PlanService) CreatePlan(ctx context.Context, input CreatePlanInput) (*domain.Plan, error) {
	// Check if plan with slug already exists
	existing, err := s.planRepo.FindBySlug(ctx, input.Slug)
	if err == nil && existing != nil {
		return nil, ErrPlanAlreadyExists
	}

	plan := &domain.Plan{
		Name:        input.Name,
		Slug:        input.Slug,
		Description: input.Description,
		Price:       input.Price,
		Currency:    input.Currency,
		Interval:    input.Interval,
		MaxUsers:    input.MaxUsers,
		MaxProducts: input.MaxProducts,
		Features:    input.Features,
		Active:      true,
	}

	if err := s.planRepo.Create(ctx, plan); err != nil {
		return nil, fmt.Errorf("CreatePlan: failed to create plan: %w", err)
	}

	return plan, nil
}

func (s *PlanService) GetPlan(ctx context.Context, id uint) (*domain.Plan, error) {
	plan, err := s.planRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetPlan: failed to get plan: %w", err)
	}
	if plan == nil {
		return nil, ErrPlanNotFound
	}
	return plan, nil
}

func (s *PlanService) ListPlans(ctx context.Context) ([]*domain.Plan, error) {
	plans, err := s.planRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("ListPlans: failed to list plans: %w", err)
	}
	return plans, nil
}

func (s *PlanService) ListActivePlans(ctx context.Context) ([]*domain.Plan, error) {
	plans, err := s.planRepo.ListActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("ListActivePlans: failed to list active plans: %w", err)
	}
	return plans, nil
}

func (s *PlanService) UpdatePlan(ctx context.Context, id uint, input UpdatePlanInput) (*domain.Plan, error) {
	plan, err := s.planRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("UpdatePlan: failed to get plan: %w", err)
	}
	if plan == nil {
		return nil, ErrPlanNotFound
	}

	// Update fields if provided
	if input.Name != nil {
		plan.Name = *input.Name
	}
	if input.Description != nil {
		plan.Description = *input.Description
	}
	if input.Price != nil {
		plan.Price = *input.Price
	}
	if input.Currency != nil {
		plan.Currency = *input.Currency
	}
	if input.Interval != nil {
		plan.Interval = *input.Interval
	}
	if input.MaxUsers != nil {
		plan.MaxUsers = *input.MaxUsers
	}
	if input.MaxProducts != nil {
		plan.MaxProducts = *input.MaxProducts
	}
	if input.Features != nil {
		plan.Features = *input.Features
	}
	if input.Active != nil {
		plan.Active = *input.Active
	}

	if err := s.planRepo.Update(ctx, plan); err != nil {
		return nil, fmt.Errorf("UpdatePlan: failed to update plan: %w", err)
	}

	return plan, nil
}

func (s *PlanService) DeletePlan(ctx context.Context, id uint) error {
	plan, err := s.planRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("DeletePlan: failed to get plan: %w", err)
	}
	if plan == nil {
		return ErrPlanNotFound
	}

	if err := s.planRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("DeletePlan: failed to delete plan: %w", err)
	}

	return nil
}
