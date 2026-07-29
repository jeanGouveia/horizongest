package repository

import (
	"errors"
	"fmt"

	"github.com/jeanGouveia/horizongest/backend/internal/infra/pg"
	"gorm.io/gorm"
)

// Repository-specific error types

// ErrDuplicateKey representa erro de chave duplicada no repository
type ErrDuplicateKey struct {
	Constraint string
	Table      string
	Err        error
}

func (e *ErrDuplicateKey) Error() string {
	if e.Constraint != "" {
		return fmt.Sprintf("duplicate key violates constraint %q", e.Constraint)
	}
	return "duplicate key violation"
}

func (e *ErrDuplicateKey) Unwrap() error {
	return e.Err
}

// ErrForeignKeyViolation representa erro de violação de foreign key
type ErrForeignKeyViolation struct {
	Constraint string
	Table      string
	Err        error
}

func (e *ErrForeignKeyViolation) Error() string {
	if e.Constraint != "" {
		return fmt.Sprintf("foreign key violation on constraint %q", e.Constraint)
	}
	return "foreign key violation"
}

func (e *ErrForeignKeyViolation) Unwrap() error {
	return e.Err
}

// IsDuplicateKeyError verifica se o erro é uma violação de unique constraint
// Esta função é a interface principal para detecção de duplicate key no repository
// Ela tenta múltiplas estratégias em ordem de robustez:
// 1. Verifica custom ErrDuplicateKey (nosso tipo)
// 2. Verifica pg.IsUniqueViolation (SQLSTATE 23505 via pgconn.PgError)
// 3. Verifica gorm.ErrDuplicatedKey (se TranslateError estiver habilitado)
func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	// 1. Verifica nosso custom error type
	var dupErr *ErrDuplicateKey
	if errors.As(err, &dupErr) {
		return true
	}

	// 2. Verifica SQLSTATE 23505 via pgconn.PgError (mais robusto)
	if pg.IsUniqueViolation(err) {
		return true
	}

	// 3. Verifica gorm.ErrDuplicatedKey (fallback para TranslateError)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	return false
}

// WrapDuplicateKeyError converte um erro de banco em ErrDuplicateKey com detalhes
func WrapDuplicateKeyError(err error) error {
	if err == nil {
		return err
	}

	if !pg.IsUniqueViolation(err) {
		return err
	}

	details := pg.ExtractPgErrorDetails(err)
	if details == nil {
		return &ErrDuplicateKey{Err: err}
	}

	return &ErrDuplicateKey{
		Constraint: details.Constraint,
		Table:      details.Table,
		Err:        err,
	}
}

// IsForeignKeyError verifica se o erro é uma violação de foreign key
func IsForeignKeyError(err error) bool {
	if err == nil {
		return false
	}

	var fkErr *ErrForeignKeyViolation
	if errors.As(err, &fkErr) {
		return true
	}

	if pg.IsForeignKeyViolation(err) {
		return true
	}

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return true
	}

	return false
}

// WrapForeignKeyError converte um erro de banco em ErrForeignKeyViolation com detalhes
func WrapForeignKeyError(err error) error {
	if err == nil {
		return err
	}

	if !pg.IsForeignKeyViolation(err) {
		return err
	}

	details := pg.ExtractPgErrorDetails(err)
	if details == nil {
		return &ErrForeignKeyViolation{Err: err}
	}

	return &ErrForeignKeyViolation{
		Constraint: details.Constraint,
		Table:      details.Table,
		Err:        err,
	}
}
