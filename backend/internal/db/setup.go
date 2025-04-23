package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rx3lixir/agg-api/config"
)

// ConnectionConfig это структура для хранения
// данных подключения к базе постгрес
type ConnectionConfig struct {
	URI            string
	Database       string
	CollectionName string
	Username       string
	Password       string
	Timeout        time.Duration
}

// NewDefaultConfig создает конфигурацию по умолчанию
func NewDefaultConfig(uri, database, collection string) *ConnectionConfig {
	return &ConnectionConfig{
		URI:            uri,
		Database:       database,
		CollectionName: collection,
		Timeout:        10 * time.Second,
	}
}

func ConnectPostgres(ctx context.Context, cfg *config.AppConfig) (*pgx.Conn, error) {
	// Контекст для подключения с таймаутом
	ctx, cancel := context.WithTimeout(ctx, cfg.DB.ConnectTimeout)
	defer cancel()

	// Подключаемся с контекстом
	conn, err := pgx.Connect(ctx, cfg.DB.DSN())
	if err != nil {
		return nil, err
	}

	// Пингуем базу данных
	if err = conn.Ping(ctx); err != nil {
		return nil, err
	}

	return conn, nil
}
