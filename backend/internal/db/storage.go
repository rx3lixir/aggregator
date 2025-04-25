package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rx3lixir/agg-api/internal/models"
)

type Storage interface {
	CreateAccount(*models.Account) error
	UpdateAccount(*models.Account) error
	GetAccounts() ([]*models.Account, error)
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

func (s *PostgresStore) GetAccounts() ([]*models.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx, "SELECT * FROM account")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []*models.Account{}

	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PostgresStore) GetAccountByID(id int) (*models.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx, "SELECT * FROM account where id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d", id)
}

func (s *PostgresStore) DeleteAccount(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var exists bool

	err := s.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM account WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("account with ID %d not found", id)
	}

	_, err = s.pool.Exec(ctx, "DELETE FROM account WHERE id = $1", id)
	return err
}

func scanIntoAccount(rows pgx.Rows) (*models.Account, error) {
	account := new(models.Account)

	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.PasswordHash,
		&account.Email,
		&account.Events,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return account, err
}
