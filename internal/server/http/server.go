package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"song-library/internal/config"
	"song-library/pkg/logger"
	"time"
)

// Http server struct.
type Server struct {
	httpServer *http.Server
}

// Returns new configured http server.
func New(cfg *config.Config, handler http.Handler) *Server {
	logger.Info(fmt.Sprintf("http server starting on port: %d", cfg.HTTP.Port))

	return &Server{
		httpServer: &http.Server{
			Addr:           fmt.Sprintf(":%d", cfg.HTTP.Port),
			Handler:        handler,
			MaxHeaderBytes: 1 << 20, // 1048576 bytes or ~1.04Mb
			ReadTimeout:    time.Minute,
			WriteTimeout:   5 * time.Minute,
		},
	}
}

// Runs configured http server instance.
func (s *Server) MustRun() {
	err := s.httpServer.ListenAndServe()
	if err != nil {
		logger.Fatal("couldn`t start http server: " + err.Error())
	}
}

// Stops configured http server.
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
