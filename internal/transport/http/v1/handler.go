package v1

import (
	"song-library/internal/config"
	"song-library/internal/service"

	music "song-library/internal/transport/http/v1/music"

	"github.com/gin-gonic/gin"
)

type V1Handler struct {
	services *service.Service
}

func NewV1Handler(s *service.Service) *V1Handler {
	return &V1Handler{
		services: s,
	}
}

func (h *V1Handler) InitV1Handler(c *config.Config, api *gin.RouterGroup) {
	musicHandler := music.New(h.services)

	v1 := api.Group("/")
	{
		musicHandler.InitMusicRoute(v1)

	}
}
