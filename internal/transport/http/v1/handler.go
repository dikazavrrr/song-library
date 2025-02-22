package httpserver

import (
	"context"
	"net/http"
	config "song-library/internal/config"
	"time"

	"song-library/pkg/logger"
)

// Server структура для HTTP-сервера
type Server struct {
	httpServer *http.Server
}

// New создает новый HTTP-сервер
func New(cfg config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.HTTPServer.Address,
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// MustRun запускает сервер и падает с фатальной ошибкой в случае неудачи
func (s *Server) MustRun() {
	logger.Info("Starting HTTP server on " + s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("HTTP server error: " + err.Error())
	}
}

// Stop корректно останавливает сервер
func (s *Server) Stop(ctx context.Context) {
	logger.Info("Shutting down HTTP server...")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Error("HTTP server shutdown error: " + err.Error())
	}
}
