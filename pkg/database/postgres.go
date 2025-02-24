package postgres

import (
	"database/sql"
	"fmt"
	"song-library/internal/config"
	"song-library/pkg/logger"
	"time"

	_ "github.com/lib/pq"

	"github.com/tinrab/retry"
	"go.uber.org/zap"
)

func NewPostgresConnection(cfg *config.Postgres) (*sql.DB, error) {
	var db *sql.DB

	var err error

	retry.ForeverSleep(time.Second*2, func(i int) error {
		fmt.Printf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.SslMode, cfg.Password)

		db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.SslMode, cfg.Password))
		if err != nil {
			logger.Error(fmt.Sprintf("postgres connection, err - %v", err),
				zap.Any("host", cfg.Host),
				zap.Any("dbname", cfg.DBName),
			)

			return err
		}

		db.SetMaxOpenConns(100)
		db.SetMaxIdleConns(5)

		logger.Info(fmt.Sprintf("ping postgres connection count - %d", i),
			zap.Any("host", cfg.Host),
			zap.Any("dbname", cfg.DBName),
		)

		err = db.Ping()
		if err != nil {
			logger.Error(fmt.Sprintf("postgres ping error - %v", err),
				zap.Any("host", cfg.Host),
				zap.Any("dbname", cfg.DBName),
			)

			return err
		}

		return nil
	})

	return db, nil
}
