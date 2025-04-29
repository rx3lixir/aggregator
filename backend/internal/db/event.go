package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rx3lixir/agg-api/internal/models"
)

func (s *PostgresStore) CreateEvent(parentCtx context.Context, event *models.Event) error {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*3)
	defer cancel()

	query := `
		INSERT INTO events (name, description, category_id, date, time, location, price, rating, image, source)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	err := s.db.QueryRow(
		ctx,
		query,
		event.Name,
		event.Description,
		event.CategoryId,
		event.Date,
		event.Time,
		event.Location,
		event.Price,
		event.Rating,
		event.Image,
		event.Source,
	).Scan(&event.Id, &event.CreatedAt, &event.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

func (s *PostgresStore) UpdateEvent(parentCtx context.Context, event *models.Event) error {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*3)
	defer cancel()

	var exists bool

	// Проверка существования ивента в базе
	err := s.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)", event.Id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("event with Id %d not found", event.Id)
	}

	query := `UPDATE events SET name = $1, description = $2, category_id = $3, date = $4, time = $5, location = $6, price = $7, rating = $8, image = $9, source = $10, updated_at = NOW() WHERE id = $11 RETURNING updated_at`

	err = s.db.QueryRow(
		ctx,
		query,
		event.Name,
		event.Description,
		event.CategoryId,
		event.Date,
		event.Time,
		event.Location,
		event.Price,
		event.Rating,
		event.Image,
		event.Source,
		event.Id,
	).Scan(&event.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update event %d: %w", event.Id, err)
	}

	return nil
}

func (s *PostgresStore) GetEvents(parentCtx context.Context) ([]*models.Event, error) {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*3)
	defer cancel()

	query := `SELECT id, name, description, category_id, date, time, location, price, rating, image, source, created_at, updated_at FROM events`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*models.Event{}

	for rows.Next() {
		event, err := scanIntoEvent(rows)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating event rows: %w", err)
	}

	return events, nil
}

func (s *PostgresStore) GetEventByID(parentCtx context.Context, id int) (*models.Event, error) {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*3)
	defer cancel()

	query := `SELECT id, name, description, category_id, date, time, location, price, rating, image, source, created_at, updated_at FROM events WHERE id = $1`

	row := s.db.QueryRow(ctx, query)

	event := new(models.Event)

	err := row.Scan(
		&event.Id,
		&event.Name,
		&event.Description,
		&event.CategoryId,
		&event.Date,
		&event.Time,
		&event.Location,
		&event.Price,
		&event.Rating,
		&event.Image,
		&event.Source,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("event %d not found", id)
		}
		return nil, fmt.Errorf("failed to get event by id %d: %w", id, err)
	}

	return event, nil
}

func (s *PostgresStore) DeleteEvent(parentCtx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*3)
	defer cancel()

	cmdTag, err := s.db.Exec(ctx, "DELETE FROM events WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete event %d: %w", id, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("event with ID %d not found for deletion", id)
	}

	return nil
}

func (s *PostgresStore) GetEventsByCategory(parentCtx context.Context, categoryID int) ([]*models.Event, error) {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*3)
	defer cancel()

	query := `SELECT id, name, description, category_id, date, time, location, price, rating, image, source, created_at, updated_at FROM events WHERE category_id = $1`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*models.Event{}

	for rows.Next() {
		event, err := scanIntoEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating event rows by category: %w", err)
	}

	return events, nil
}

func scanIntoEvent(rows pgx.Rows) (*models.Event, error) {
	event := new(models.Event)

	err := rows.Scan(
		&event.Id,
		&event.Name,
		&event.Description,
		&event.CategoryId,
		&event.Date,
		&event.Time,
		&event.Location,
		&event.Price,
		&event.Rating,
		&event.Image,
		&event.Source,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return event, nil
}
