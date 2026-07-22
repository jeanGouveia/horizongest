package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/infra/repository"
	"gorm.io/gorm"
)

// Seed cria dados iniciais no banco de dados
func Seed(db *gorm.DB) error {
	ctx := context.Background()

	// Criar repositórios
	platformUserRepo := repository.NewGormPlatformUserRepository(db)

	// Verificar se já existe usuário admin
	existing, err := platformUserRepo.FindByEmail(ctx, "admin@platform.com")
	if err != nil {
		return fmt.Errorf("erro ao verificar usuário admin: %w", err)
	}
	if existing != nil {
		log.Println("Usuário admin já existe, pulando seed")
		return nil
	}

	// Criar usuário admin padrão
	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("erro ao gerar hash da senha: %w", err)
	}

	adminUser := &domain.PlatformUser{
		Name:         "Administrador",
		Email:        "admin@platform.com",
		PasswordHash: string(hash),
		Role:         domain.PlatformRoleAdmin,
		Active:       true,
	}

	if err := platformUserRepo.Create(ctx, adminUser); err != nil {
		return fmt.Errorf("erro ao criar usuário admin: %w", err)
	}

	log.Println("Usuário admin criado com sucesso: admin@platform.com / admin123")
	return nil
}

// RunSeed executa o seed se a variável de ambiente RUN_SEED=true
func RunSeed(db *gorm.DB) {
	if os.Getenv("RUN_SEED") == "true" {
		log.Println("Executando seed...")
		if err := Seed(db); err != nil {
			log.Fatalf("Erro ao executar seed: %v", err)
		}
	}
}
