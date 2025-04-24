package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rx3lixir/agg-api/internal/models"
)

type Storage interface {
	CreateAccount(*models.Account) error
	UpdateAccount(*models.Account) error
	GetAccountByID(id int) (*models.Account, error)
	DeleteAccount(id int) error
}

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPosgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{
		pool: pool,
	}
}

func (s *PostgresStore) CreateAccount(*models.Account) error {
	return nil
}

func (s *PostgresStore) UpdateAccount(*models.Account) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*models.Account, error) {
	return nil, nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}
