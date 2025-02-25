package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteSong удаляет песню по ID.
//
// @Summary Удалить песню
// @Description Удаляет песню из базы данных по её ID.
// @Tags Songs
// @Param id path int true "ID песни"
// @Success 200 {object} map[string]string "Сообщение об успешном удалении"
// @Failure 400 {object} map[string]string "Ошибка: некорректный ID"
// @Failure 500 {object} map[string]string "Ошибка при удалении"
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	if err := h.Repo.DeleteSong(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}
