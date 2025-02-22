//nolint:revive //так уже сделано
package postgres_db

import (
	"database/sql"
	"song-library/internal/config"
)

// Postgres repository.
type PostgresRepository struct {
}

// NewPostgresRepository creates a new instance of PostgresRepository.
func NewPostgresRepository(mainConn *sql.DB, cfg *config.Config) *PostgresRepository {
	return &PostgresRepository{}
}
