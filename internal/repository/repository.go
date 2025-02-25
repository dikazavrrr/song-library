package repository

import (
	"song-library/pkg/domain"

	"github.com/jinzhu/gorm"
)

type SongRepository struct {
	DB *gorm.DB
}

func NewSongRepository(db *gorm.DB) *SongRepository {
	return &SongRepository{DB: db}
}

func (r *SongRepository) GetAllSongs() ([]domain.Songs, error) {
	var songs []domain.Songs
	err := r.DB.Find(&songs).Error
	return songs, err
}

func (r *SongRepository) DeleteSong(id uint) error {
	if err := r.DB.Delete(&domain.Songs{}, id).Error; err != nil {
		return err
	}
	return nil
}

// GetSongByID получает песню по ID
func (r *SongRepository) GetSongByID(id uint64) (*domain.Songs, error) {
	var song domain.Songs
	if err := r.DB.First(&song, id).Error; err != nil {
		return nil, err
	}
	return &song, nil
}

// CreateSong добавляет новую песню в базу данных
func (r *SongRepository) CreateSong(song *domain.Songs) error {
	return r.DB.Create(song).Error
}

// UpdateSong обновляет данные существующей песни
func (r *SongRepository) UpdateSong(song domain.Songs) (*domain.Songs, error) {
	if err := r.DB.Save(&song).Error; err != nil {
		return nil, err
	}
	return &song, nil
}

func (r *SongRepository) CreateSong1(song domain.Songs) (*domain.Songs, error) {
	if err := r.DB.Create(&song).Error; err != nil {
		return nil, err
	}
	return &song, nil
}
