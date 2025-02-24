package httphandler

import (
	"fmt"
	"song-library/docs"
	"song-library/internal/config"
	"song-library/internal/service"
	v1 "song-library/internal/transport/http/v1"
	"song-library/pkg/limiter"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/unrolled/secure"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
)

// Handler struct.
type Handler struct {
	service *service.Service
	cfg     *config.Config
}

// New создает новый обработчик HTTP.
func New(s *service.Service, cfg *config.Config) *Handler {
	return &Handler{
		service: s,
		cfg:     cfg,
	}
}

// Init инициализирует маршруты.
func (h *Handler) Init() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	secureMiddleware := secure.New(secure.Options{
		FrameDeny:          true,
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
	})

	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			if err := secureMiddleware.Process(c.Writer, c.Request); err != nil {
				c.Abort()
				return
			}
			if status := c.Writer.Status(); status > 300 && status < 399 {
				c.Abort()
			}
		}
	}()

	li := limiter.New()
	r.Use(
		gin.Recovery(),
		secureFunc,
		limiter.GinLimit(li),
	)

	docs.SwaggerInfo.Host = h.cfg.Environment.RouterHost
	docs.SwaggerInfo.BasePath = h.cfg.Environment.RootRouter

	r.GET(fmt.Sprintf("%s/swagger/*any", h.cfg.Environment.RootRouter), ginSwagger.WrapHandler(swaggerFiles.Handler))

	h.initApi(h.cfg, r)

	return r
}

func (h *Handler) initApi(c *config.Config, router *gin.Engine) {
	v1Handler := v1.NewV1Handler(h.service)
	api := router.Group(h.cfg.Environment.RootRouter)
	{
		v1Handler.InitV1Handler(c, api)
	}
}
