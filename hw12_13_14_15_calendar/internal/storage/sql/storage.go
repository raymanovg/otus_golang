package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	db *sqlx.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) (err error) {
	config, err := pgx.ParseConfig("postgres://calendar:calendar@postgres:5432/calendar?sslmode=disable")
	if err != nil {
		return err
	}
	// создается пул соединений
	s.db, err = sqlx.Open("pgx", stdlib.RegisterConnConfig(config))
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}

	return s.db.PingContext(ctx)
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO events (
						"title",
						"description",
						"begin", 
                    	"end", 
                    	"user_id",
						"created_at",
						"updated_at"
                   ) values ($1, $2, $3, $4, $5, $6, $7)`,
		event.Title,
		event.Desc,
		event.Begin,
		event.End,
		event.UserID,
		event.CreatedAt,
		event.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID int64) error {
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	return nil
}

func (s *Storage) GetAllEventsOfUser(ctx context.Context, userID int64) ([]storage.Event, error) {
	return nil, nil
}

func (s *Storage) GetAllEvents(ctx context.Context) ([]storage.Event, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT 
       id, 
       title, description, time, user_id, (EXTRACT(EPOCH FROM duration) * 1000000000)::bigint as duration, created_at, updated_at
FROM events`)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()
	var events []storage.Event
	for rows.Next() {
		var event storage.Event

		var updatedAt sql.NullTime
		var createdAt sql.NullTime

		if err = rows.Scan(
			&event.ID,
			&event.Title,
			&event.Desc,
			&event.Begin,
			&event.End,
			&event.UserID,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("cannot scan: %w", err)
		}
		if createdAt.Valid {
			event.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			event.UpdatedAt = updatedAt.Time
		}

		events = append(events, event)
	}

	return events, nil
}
