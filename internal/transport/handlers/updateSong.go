package handlers

import (
	"net/http"
	"song-library/pkg/domain"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UpdateSong обновляет информацию о песне.
//
// @Summary Обновить песню
// @Description Обновляет данные песни в базе данных по её ID.
// @Tags Songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param song body domain.Songs true "Данные песни"
// @Success 200 {object} domain.Songs
// @Failure 400 {object} map[string]string "Ошибка: некорректный ввод"
// @Failure 500 {object} map[string]string "Ошибка при обновлении"
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	var updatedSong domain.Songs
	if err := c.ShouldBindJSON(&updatedSong); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedSong.ID = uint64(id)
	song, err := h.Repo.UpdateSong(updatedSong)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}

	c.JSON(http.StatusOK, song)
}
