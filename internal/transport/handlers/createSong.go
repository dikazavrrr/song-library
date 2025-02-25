package handlers

import (
	"net/http"

	"song-library/external"
	"song-library/pkg/domain"

	"github.com/gin-gonic/gin"
)

func (h *SongHandler) CreateSong(c *gin.Context) {
	var input struct {
		GroupName string `json:"group" binding:"required"`
		SongName  string `json:"song" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Получаем информацию о песне из внешнего API
	songInfo, err := external.FetchSongInfo(input.GroupName, input.SongName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch song info"})
		return
	}

	// Создаем объект песни для сохранения
	newSong := domain.Songs{
		SongName:    input.SongName,
		GroupName:   input.GroupName,
		ReleaseDate: songInfo.ReleaseDate,
		Lyrics:      songInfo.Text,
		Link:        songInfo.Link,
	}

	// Сохраняем в БД
	if err := h.Repo.CreateSong(&newSong); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save song"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Song added successfully", "song": newSong})
}
