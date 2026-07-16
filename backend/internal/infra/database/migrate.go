package database

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/jeanGouveia/pratoOnline/backend/internal/infra/repository"
)

// RunMigrations executa o AutoMigrate do GORM para todas as tabelas do sistema.
// Utiliza exclusivamente os GormModels definidos nos repositories (única fonte de verdade).
// Em produção, substitua por Goose com migrations SQL versionadas.
func RunMigrations(db *gorm.DB) error {
	models := []interface{}{
		&repository.GormUserModel{},
		&repository.GormProduct{},
		&repository.GormIngredient{},
		&repository.GormProductIngredient{},
		&repository.GormOrder{},
		&repository.GormOrderItem{},
		&repository.GormStockAdjustmentPending{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("RunMigrations: %w", err)
	}

	return nil
}
