package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rx3lixir/agg-api/internal/models"
)

// Интерфейс для абстракции методов базы данных от pgxpool
type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Storage interface {
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	GetUsers(ctx context.Context) ([]*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type PostgresStore struct {
	db DBTX
}

func NewPosgresStore(db DBTX) *PostgresStore {
	return &PostgresStore{
		db: db,
	}
}

func (s *PostgresStore) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (name, email, password, is_admin)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
		user.IsAdmin,
	).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("Failed to create user: %w", err)
	}

	return nil
}

func (s *PostgresStore) UpdateUser(ctx context.Context, user *models.User) error {
	var exists bool

	// Проверка есть ли запрашиваемый пользователь
	err := s.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", user.Id).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("user with ID %d not found", user.Id)
	}

	query := `
		UPDATE users
		SET name = $1, email = $2, password = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at
	`
	err = s.db.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
		user.Id).Scan(&user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update user %d: %w", user.Id, err)
	}

	return err
}

func (s *PostgresStore) GetUsers(ctx context.Context) ([]*models.User, error) {
	rows, err := s.db.Query(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*models.User{}

	for rows.Next() {
		user, err := scanIntoUser(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}

func (s *PostgresStore) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	row := s.db.QueryRow(ctx, "SELECT id, name, email, password, is_admin, created_at, updated_at FROM users WHERE id = $1", id)

	user := new(models.User)
	err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user %d not found", id)
		}
		return nil, fmt.Errorf("failed to get user by id %d: %w", id, err)
	}

	return user, nil
}

func (s *PostgresStore) DeleteUser(ctx context.Context, id int) error {
	cmdTag, err := s.db.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failedt to delete user %d: %w", id, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user with ID %d not found for deletion", id)
	}

	return nil
}

func scanIntoUser(rows pgx.Rows) (*models.User, error) {
	user := new(models.User)

	err := rows.Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, err
}
