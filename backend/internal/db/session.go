package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rx3lixir/agg-api/internal/models"
)

func (s *PostgresStore) CreateSession(parentContext context.Context, session *models.Session) (*models.Session, error) {
	ctx, cancel := context.WithTimeout(parentContext, time.Second*3)
	defer cancel()

	query := `
	INSERT INTO sessions (user_email, refresh_token, is_revoked, expires_at)
	VALUES ($1, $2, $3, $4)`

	err := s.db.QueryRow(
		ctx,
		query,
		session.UserEmail,
		session.RefreshToken,
		session.IsRevoked,
		session.ExpiresAt,
	).Scan(&session.Id)

	if err != nil {
		return nil, fmt.Errorf("failed to create session for %v: %d", session.UserEmail, err)
	}

	return session, nil
}

func (s *PostgresStore) GetSession(parentContext context.Context, id string) (*models.Session, error) {
	ctx, cancel := context.WithTimeout(parentContext, time.Second*3)
	defer cancel()

	query := `SELECT id, user_email, refresh_token, is_revoked, created_at, expires_at FROM sessions WHERE id=$`

	row := s.db.QueryRow(ctx, query, id)

	session := new(models.Session)

	err := row.Scan(
		session.Id,
		session.UserEmail,
		session.RefreshToken,
		session.IsRevoked,
		session.CreatedAt,
		session.ExpiresAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("session %s not found", id)
		}
		return nil, fmt.Errorf("failed to get session by id %s: %w", id, err)
	}

	return session, nil
}

func (s *PostgresStore) RevokeSession(parentCtx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*3)
	defer cancel()

	var exists bool

	// Проверка существования ивента в базе
	err := s.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM sessions WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("session with Id %v not found", id)
	}

	query := `UPDATE sessions SET is_revoked=1 WHERE id=$1`

	row, err := s.db.Query(
		ctx,
		query,
		id,
	)
	defer row.Close()

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("session %v not found", id)
		}
		return fmt.Errorf("failed to get session by id %v: %w", id, err)
	}

	return nil
}

func (s *PostgresStore) DeleteSession(parentCtx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*3)
	defer cancel()

	row, err := s.db.Query(ctx, "DELETE FROM sessions WHERE id=$1", id)
	defer row.Close()
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("session %v not found", id)
		}
		return fmt.Errorf("failed to get session by id %v: %w", id, err)
	}

	return nil
}
