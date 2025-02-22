package repository

import (
	"song-library/internal/repository/database"
)

// Repository Layer.
type Repository struct {
	DB *database.Repository
}

func New(dbRepo *database.Repository) *Repository {
	return &Repository{
		DB: dbRepo,
	}
}
