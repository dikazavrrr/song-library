package spam

import (
	"song-library/internal/service"

	"github.com/gin-gonic/gin"
)

type MusicHandler struct {
	services *service.Service
}

func New(s *service.Service) *MusicHandler {
	return &MusicHandler{
		services: s,
	}
}

// @BasePath /api/spam/
// Inits Spam Handler router group.
func (mH *MusicHandler) InitMusicRoute(v1 *gin.RouterGroup) {
	spam := v1.Group("/info")
	{
		spam.GET("/", mH.getAllMusic)
	}
}
