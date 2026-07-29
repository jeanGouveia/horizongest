package pg

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// SQLSTATE codes for PostgreSQL errors
// Reference: https://www.postgresql.org/docs/current/errcodes-appendix.html
const (
	SQLStateUniqueViolation     = "23505" // unique_violation
	SQLStateForeignKeyViolation = "23503" // foreign_key_violation
)

// IsUniqueViolation verifica se o erro é uma violação de unique constraint PostgreSQL
// Usa SQLSTATE 23505 (unique_violation) que é padronizado pelo PostgreSQL
func IsUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	// Tenta extrair PgError do erro
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == SQLStateUniqueViolation
	}

	return false
}

// IsForeignKeyViolation verifica se o erro é uma violação de foreign key
// Usa SQLSTATE 23503 (foreign_key_violation)
func IsForeignKeyViolation(err error) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == SQLStateForeignKeyViolation
	}

	return false
}

// GetConstraintName extrai o nome da constraint de um erro PostgreSQL
// Retorna string vazia se não conseguir extrair
// Nota: Em pgx v5, o nome da constraint pode estar em diferentes campos dependendo do contexto
func GetConstraintName(err error) string {
	if err == nil {
		return ""
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// Tenta extrair do campo ConstraintName se disponível
		// Se não estiver disponível, pode estar na mensagem
		return pgErr.ConstraintName
	}

	return ""
}

// GetTableName extrai o nome da tabela de um erro PostgreSQL
// Retorna string vazia se não conseguir extrair
func GetTableName(err error) string {
	if err == nil {
		return ""
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.TableName
	}

	return ""
}

// PgErrorDetails contém detalhes estruturados de um erro PostgreSQL
type PgErrorDetails struct {
	Code         string // SQLSTATE
	Message      string
	Constraint   string
	Table        string
	Column       string
	Severity     string
	Detail       string
	Hint         string
	Schema       string
	DataTypeName string
}

// ExtractPgErrorDetails extrai todos os detalhes disponíveis de um erro PostgreSQL
// Retorna nil se não for um erro PostgreSQL
func ExtractPgErrorDetails(err error) *PgErrorDetails {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return nil
	}

	return &PgErrorDetails{
		Code:         pgErr.Code,
		Message:      pgErr.Message,
		Constraint:   pgErr.ConstraintName,
		Table:        pgErr.TableName,
		Column:       pgErr.ColumnName,
		Severity:     pgErr.Severity,
		Detail:       pgErr.Detail,
		Hint:         pgErr.Hint,
		Schema:       pgErr.SchemaName,
		DataTypeName: pgErr.DataTypeName,
	}
}
