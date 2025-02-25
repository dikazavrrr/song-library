package router

import (
	"song-library/internal/repository"
	songHandler "song-library/internal/transport/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	songRepo := repository.NewSongRepository(db)
	songHandler := songHandler.NewSongHandler(songRepo)

	songs := r.Group("/songs")
	{
		// Получить все песни
		songs.GET("/", songHandler.GetAllSongs)

		// Получить песню по ID
		songs.GET("/:id", songHandler.GetSongByID)

		// Обновить данные песни
		songs.PUT("/:id", songHandler.UpdateSong)

		// Удалить песню
		songs.DELETE("/:id", songHandler.DeleteSong)

		// Получить текст песни с пагинацией по куплетам
		songs.GET("/info", songHandler.GetSongLyricsWithPagination)
	}

	return r
}
