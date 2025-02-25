package main

import (
	"log"
	"song-library/internal/config"
	router "song-library/internal/transport"
	db "song-library/migrations"

	"song-library/pkg/domain"

	_ "song-library/docs"

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Song API
// @version 1.0
// @description API для управления песнями.
// @host localhost:8080
// @BasePath /
func main() {
	log.Println("[INFO] Starting Song API service...")

	// Загрузка переменных окружения из .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	config.LoadConfig()
	log.Println("External API URL:", config.AppConfig.API.ExternalURL)

	log.Println("[INFO] .env file loaded successfully")

	// Инициализация базы данных
	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbConn.Close()
	log.Println("[INFO] Connected to database")

	sqlDB := dbConn.DB()
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("[ERROR] Database is not reachable: %v", err)
	}

	// Автоматическая миграция модели Song
	dbConn.AutoMigrate(&domain.Songs{})

	log.Println("[INFO] Database migrated successfully")
	// Настройка маршрутов
	r := router.SetupRoutes(dbConn)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Запуск сервера на указанном порту (например, 8080)
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
