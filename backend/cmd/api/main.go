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

	// Загрузка и инициализация конфигурации приложения
	cfg, err := config.New(log)
	if err != nil {
		log.Error("Failed to initialize config", err)
	}

	// Создаение пула подключений Postgres
	pool, err := db.CreatePostgresPool(ctx, cfg)
	if err != nil {
		log.Error("Failed to initialize config", err)
	}
	defer pool.Close()

	// Создаение хранилища с инициализированным пулом подключений
	store := db.NewPosgresStore(pool)

	log.Info("Хранилище инициализированно", "db", store)

	// Инициализация и запуск сервера с заданными параметрами
	server := api.NewAPIServer(cfg.Server.Address, log, store)
	server.Run()
}
