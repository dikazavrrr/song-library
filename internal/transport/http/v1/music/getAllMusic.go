package spam

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (sh *MusicHandler) getAllMusic(c *gin.Context) {

	c.JSON(http.StatusOK, "ok")
}
