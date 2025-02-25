package handlers

import (
	"song-library/internal/repository"

	"github.com/jinzhu/gorm"
)

// SongHandler обрабатывает запросы к песням
type SongHandler struct {
	Repo *repository.SongRepository
	DB   *gorm.DB
}

// NewSongHandler создает новый обработчик
func NewSongHandler(repo *repository.SongRepository) *SongHandler {
	return &SongHandler{Repo: repo}
}
