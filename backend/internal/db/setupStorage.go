package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rx3lixir/agg-api/config"
)

func CreatePostgresPool(ctx context.Context, cfg *config.AppConfig) (*pgxpool.Pool, error) {
	c, cancel := context.WithTimeout(ctx, cfg.DB.ConnectTimeout)
	defer cancel()

	pool, err := pgxpool.New(c, cfg.DB.DSN())
	if err != nil {
		return nil, err
	}

	// Проверяем соединение
	if err = pool.Ping(c); err != nil {
		return nil, err
	}

	return pool, nil
}
