package main

import (
	"context"
	"os"

	"github.com/rx3lixir/agg-api/config"
	"github.com/rx3lixir/agg-api/internal/api"
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

	server := api.NewAPIServer(cfg.Server.Address, log)
	server.Run()

	dbPool, err := db.CreatePostgresPool(ctx, cfg)
	if err != nil {
		log.Error("Failed to initialize config", err)
	}
	defer dbPool.Close()

	log.Info("Connection successfully established",
		"connection_info:", dbPool.Config().ConnConfig,
		"closed:", dbPool.Stat())
}
