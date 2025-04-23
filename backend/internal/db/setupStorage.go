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

	// Создаем таблцу, если не существует
	if err = initSchema(c, pool); err != nil {
		return nil, err
	}

	return pool, nil
}

func initSchema(ctx context.Context, pool *pgxpool.Pool) error {
	tableCreate := `
		CREATE TABLE IF NOT EXISTS events (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			date TEXT,
			type TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`

	_, err := pool.Exec(ctx, tableCreate)

	return err
}
