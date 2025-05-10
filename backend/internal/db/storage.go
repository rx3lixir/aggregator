package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rx3lixir/agg-api/config"
	"github.com/rx3lixir/agg-api/internal/models"
)

// Интерфейс для абстракции методов базы данных от pgxpool
type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// SessionStorage определяет интерфейс для работы с сессиями
type SessionStorage interface {
	CreateSession(ctx context.Context, session *models.Session) (*models.Session, error)
	GetSession(ctx context.Context, id string) (*models.Session, error)
	RevokeSession(ctx context.Context, id string) error
	DeleteSession(ctx context.Context, id string) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
	Close() error
}

// Методы для User
type UserStore interface {
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	GetUsers(ctx context.Context) ([]*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(parentCtx context.Context, email string) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error
}

// Методы для Event
type EventStore interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	UpdateEvent(ctx context.Context, event *models.Event) error
	GetEvents(ctx context.Context) ([]*models.Event, error)
	GetEventByID(ctx context.Context, id int) (*models.Event, error)
	DeleteEvent(ctx context.Context, id int) error
	GetEventsByCategory(ctx context.Context, categoryID int) ([]*models.Event, error)
}

// Методы для Categories
type CategoryStore interface {
	CreateCategory(ctx context.Context, category *models.Category) error
	GetCategories(ctx context.Context) ([]*models.Category, error)
	GetCategoryByID(ctx context.Context, id int) (*models.Category, error)
	UpdateCategory(ctx context.Context, category *models.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type Storage interface {
	UserStore
	EventStore
	CategoryStore
}

type PostgresStore struct {
	db DBTX
}

func NewPosgresStore(db DBTX) *PostgresStore {
	return &PostgresStore{
		db: db,
	}
}

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
