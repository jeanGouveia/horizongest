package repository

import (
	"context"

	"github.com/jeanGouveia/pratoOnline/backend/internal/domain"
	"github.com/jeanGouveia/pratoOnline/backend/internal/ports"
	"gorm.io/gorm"
)

type GormPlanRepository struct {
	db *gorm.DB
}

func NewGormPlanRepository(db *gorm.DB) ports.PlanRepository {
	return &GormPlanRepository{db: db}
}

func (r *GormPlanRepository) Create(ctx context.Context, plan *domain.Plan) error {
	return r.db.WithContext(ctx).Create(plan).Error
}

func (r *GormPlanRepository) FindByID(ctx context.Context, id uint) (*domain.Plan, error) {
	var plan domain.Plan
	err := r.db.WithContext(ctx).First(&plan, id).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *GormPlanRepository) FindBySlug(ctx context.Context, slug string) (*domain.Plan, error) {
	var plan domain.Plan
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *GormPlanRepository) Update(ctx context.Context, plan *domain.Plan) error {
	return r.db.WithContext(ctx).Save(plan).Error
}

func (r *GormPlanRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Plan{}, id).Error
}

func (r *GormPlanRepository) List(ctx context.Context) ([]*domain.Plan, error) {
	var plans []*domain.Plan
	err := r.db.WithContext(ctx).Find(&plans).Error
	return plans, err
}

func (r *GormPlanRepository) ListActive(ctx context.Context) ([]*domain.Plan, error) {
	var plans []*domain.Plan
	err := r.db.WithContext(ctx).Where("active = ?", true).Find(&plans).Error
	return plans, err
}
