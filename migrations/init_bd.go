package db

import (
	"fmt"
	"log"
	"os"

	"song-library/pkg/domain"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// InitDB инициализирует соединение с БД
func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Автомиграция схемы БД
	if err := db.AutoMigrate(&domain.Songs{}).Error; err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	return db, nil
}
