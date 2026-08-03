package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/ports"
)

type GormCompanyModel struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	Name           string `gorm:"not null"`
	Slug           string `gorm:"uniqueIndex;not null"`
	Description    string `gorm:"type:text"`
	Active         bool   `gorm:"not null;default:true"`
	LogoURL        string `gorm:"type:varchar(500)"`
	PrimaryColor   string `gorm:"type:varchar(7);default:'#3b82f6'"`
	SecondaryColor string `gorm:"type:varchar(7);default:'#1e40af'"`
	// Business Engine fields (Sprint 3)
	BusinessType string     `gorm:"type:text;default:'generic'"`
	Locale       string     `gorm:"type:text;default:'pt-BR'"`
	Currency     string     `gorm:"type:text;default:'BRL'"`
	Timezone     string     `gorm:"type:text;default:'America/Sao_Paulo'"`
	DeletedAt    *time.Time `gorm:"index"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime"`
}

func (GormCompanyModel) TableName() string { return "companies" }

var _ ports.CompanyRepository = (*GormCompanyRepository)(nil)

type GormCompanyRepository struct {
	db *gorm.DB
}

func NewGormCompanyRepository(db *gorm.DB) *GormCompanyRepository {
	return &GormCompanyRepository{db: db}
}

func (r *GormCompanyRepository) getDB(ctx context.Context, tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *GormCompanyRepository) Create(ctx context.Context, company *domain.Company) error {
	return r.CreateWithTx(ctx, company, nil)
}

func (r *GormCompanyRepository) CreateWithTx(ctx context.Context, company *domain.Company, tx *gorm.DB) error {
	model := GormCompanyModel{
		Name:           company.Name,
		Slug:           company.Slug,
		Description:    company.Description,
		Active:         company.Active,
		LogoURL:        company.LogoURL,
		PrimaryColor:   company.PrimaryColor,
		SecondaryColor: company.SecondaryColor,
		BusinessType:   string(company.BusinessType),
		Locale:         company.Locale,
		Currency:       company.Currency,
		Timezone:       company.Timezone,
	}
	if err := r.getDB(ctx, tx).Create(&model).Error; err != nil {
		return fmt.Errorf("CompanyRepository.Create: %w", err)
	}
	company.ID = model.ID
	company.CreatedAt = model.CreatedAt
	company.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *GormCompanyRepository) FindByID(ctx context.Context, id uint) (*domain.Company, error) {
	// FORENSIC: Log CompanyID received
	log.Printf("[FORENSIC] CompanyRepository.FindByID - CompanyID recebido: %d", id)

	var model GormCompanyModel

	// FORENSIC: Log SQL that will be executed
	log.Printf("[FORENSIC] CompanyRepository.FindByID - SQL a ser executado: SELECT * FROM companies WHERE id = %d AND deleted_at IS NULL", id)

	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&model, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[FORENSIC] CompanyRepository.FindByID - Registro não encontrado para CompanyID: %d", id)
		return nil, nil
	}
	if err != nil {
		log.Printf("[FORENSIC] CompanyRepository.FindByID - Erro ao executar SQL: %v", err)
		return nil, fmt.Errorf("CompanyRepository.FindByID: %w", err)
	}

	// FORENSIC: Log company returned
	log.Printf("[FORENSIC] CompanyRepository.FindByID - Empresa retornada: ID=%d, Name=%s", model.ID, model.Name)

	return toDomainCompany(&model), nil
}

func (r *GormCompanyRepository) FindBySlug(ctx context.Context, slug string) (*domain.Company, error) {
	var model GormCompanyModel
	err := r.db.WithContext(ctx).Where("slug = ? AND deleted_at IS NULL", slug).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("CompanyRepository.FindBySlug: %w", err)
	}
	return toDomainCompany(&model), nil
}

func (r *GormCompanyRepository) List(ctx context.Context) ([]domain.Company, error) {
	var models []GormCompanyModel
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("CompanyRepository.List: %w", err)
	}
	companies := make([]domain.Company, len(models))
	for i, model := range models {
		companies[i] = *toDomainCompany(&model)
	}
	return companies, nil
}

func (r *GormCompanyRepository) Update(ctx context.Context, company *domain.Company) error {
	model := GormCompanyModel{
		ID:             company.ID,
		Name:           company.Name,
		Slug:           company.Slug,
		Description:    company.Description,
		Active:         company.Active,
		LogoURL:        company.LogoURL,
		PrimaryColor:   company.PrimaryColor,
		SecondaryColor: company.SecondaryColor,
		BusinessType:   string(company.BusinessType),
		Locale:         company.Locale,
		Currency:       company.Currency,
		Timezone:       company.Timezone,
	}
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return fmt.Errorf("CompanyRepository.Update: %w", err)
	}
	company.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *GormCompanyRepository) Delete(ctx context.Context, id uint) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&GormCompanyModel{}).Where("id = ?", id).Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("CompanyRepository.Delete: %w", err)
	}
	return nil
}

func toDomainCompany(m *GormCompanyModel) *domain.Company {
	return &domain.Company{
		ID:             m.ID,
		Name:           m.Name,
		Slug:           m.Slug,
		Description:    m.Description,
		Active:         m.Active,
		LogoURL:        m.LogoURL,
		PrimaryColor:   m.PrimaryColor,
		SecondaryColor: m.SecondaryColor,
		BusinessType:   domain.BusinessType(m.BusinessType),
		Locale:         m.Locale,
		Currency:       m.Currency,
		Timezone:       m.Timezone,
		DeletedAt:      m.DeletedAt,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}
