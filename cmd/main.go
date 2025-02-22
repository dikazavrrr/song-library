package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"song-library/internal/config"
	"song-library/internal/repository"

	dbRepos "song-library/internal/repository/database"
	"song-library/internal/service"
	"song-library/internal/status"
	httphandler "song-library/internal/transport/http"
	postgres "song-library/pkg/database"

	"song-library/pkg/logger"
	"syscall"
)

func main() {
	logger.ZapLoggerInit()

	stat := status.NewStatus()
	ctx := context.Background()
	cfg := config.MustInit(os.Getenv("IS_PROD"))

	mainPC, mainPCErr := postgres.NewPostgresConnection(&cfg.Postgres)
	if mainPCErr == nil {
		logger.Info("postgres connection established!")
	}

	dbRepo := dbRepos.New(mainPC, cfg)
	r := repository.New(dbRepo)
	s := service.New(ctx, r, stat, cfg)
	hh := httphandler.New(s, cfg)
	logger.Info("transports, services, handlers instantiated!")

	hsrv := httpserver.New(cfg, hh.Init())
	go func() {
		hsrv.MustRun()
	}()
	defer hsrv.Stop(context.Background())

	defer func(postgresConns ...*sql.DB) {
		for i, pc := range postgresConns {
			if pc != nil {
				err := pc.Close()
				if err != nil {
					logger.Error(fmt.Sprintf("postgres: failed to close "+
						"connection '%d', err: %v", i, err.Error()))
				}
			}
		}
	}(mainPC)

	logger.Error("All services have started, ready to receive requests")

	awaitStop(ctx)
}

func awaitStop(ctx context.Context) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	osSignal := <-quit

	ctx.Done()

	logger.Info(fmt.Sprintf("program shutdown... call_type: %v", osSignal))
}
