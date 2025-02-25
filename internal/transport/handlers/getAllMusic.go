package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllSongs получает список всех песен.
//
// @Summary Получить все песни
// @Description Возвращает список всех песен в базе данных.
// @Tags Songs
// @Produce json
// @Success 200 {array} domain.Songs
// @Router /songs [get]
func (h *SongHandler) GetAllSongs(c *gin.Context) {
	songs, err := h.Repo.GetAllSongs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve songs"})
		return
	}
	c.JSON(http.StatusOK, songs)
}
