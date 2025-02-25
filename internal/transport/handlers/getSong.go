package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetSongByID получает информацию о песне по ID.
//
// @Summary Получить песню по ID
// @Description Возвращает данные о конкретной песне по её ID.
// @Tags Songs
// @Param id path int true "ID песни"
// @Success 200 {object} domain.Songs
// @Failure 400 {object} map[string]string "Ошибка: некорректный ID"
// @Failure 404 {object} map[string]string "Ошибка: песня не найдена"
// @Router /songs/{id} [get]
func (h *SongHandler) GetSongByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	song, err := h.Repo.GetSongByID(uint64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	c.JSON(http.StatusOK, song)
}
