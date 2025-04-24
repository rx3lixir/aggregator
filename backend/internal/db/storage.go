package db

import (
	"context"
	"fmt"
	"time"

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

func (s *PostgresStore) CreateAccount(account *models.Account) error {
	query := `
		INSERT INTO account (first_name, last_name, email, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	// Контекст для регулировки продолжительности операций
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.pool.QueryRow(
		ctx,
		query,
		account.FirstName,
		account.LastName,
		account.Email,
		account.PasswordHash,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		return fmt.Errorf("Failed to create account: %w", err)
	}

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
