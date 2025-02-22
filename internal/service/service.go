package service

import (
	"context"
	"song-library/internal/config"
	"song-library/internal/repository"
	"song-library/internal/status"
)

// Service layer.
type Service struct {
}

func New(ctx context.Context, r *repository.Repository, s status.APIStatusReader, c *config.Config) *Service {
	return &Service{}
}
