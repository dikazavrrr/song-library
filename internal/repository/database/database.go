package database

import (
	"database/sql"

	"song-library/internal/config"
	"song-library/internal/repository/database/postgres_db"
)

// Repository Layer.
type Repository struct {
	PostgresDB *postgres_db.PostgresRepository
}

func New(mainConn *sql.DB, cfg *config.Config) *Repository {
	return &Repository{
		PostgresDB: postgres_db.NewPostgresRepository(mainConn, cfg),
	}
}
