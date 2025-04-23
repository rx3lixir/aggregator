package main

import (
	"context"
	"os"

	"github.com/rx3lixir/agg-api/config"
	"github.com/rx3lixir/agg-api/internal/db"
	"github.com/rx3lixir/agg-api/internal/lib/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	loggerConfig := logger.Config{
		Level:  logger.InfoLevel,
		Format: logger.TextFormat,
		Output: os.Stdout,
	}

	log := logger.New(loggerConfig)

	cfg, err := config.New(log)
	if err != nil {
		log.Error("Failed to initialize config", err)
	}

	conn, err := db.ConnectPostgres(ctx, cfg)
	if err != nil {
		log.Error("Failed to initialize config", err)
	}
	defer conn.Close(ctx)

	// Можно также выполнить простой SQL запрос для проверки
	var version string
	err = conn.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		log.Error("Failed to query database version", err)
		return
	}

	log.Info("Connection successfully established",
		"connection_info:", conn.Config().Config,
		"closed:", conn.IsClosed(),
		"pg_version:", version)
}
