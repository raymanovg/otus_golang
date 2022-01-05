package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	conf config.SQLStorage
	db   *sqlx.DB
}

func New(config config.SQLStorage) *Storage {
	return &Storage{
		conf: config,
	}
}

func (s *Storage) Connect(ctx context.Context) (err error) {
	conf, err := pgx.ParseConfig(s.conf.DSN)
	if err != nil {
		return err
	}
	s.db, err = sqlx.Open("pgx", stdlib.RegisterConnConfig(conf))
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}
	s.db.SetMaxIdleConns(s.conf.MaxIdleConns)
	s.db.SetMaxOpenConns(s.conf.MaxOpenConns)

	return s.db.PingContext(ctx)
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	if err := storage.ValidateFull(event); err != nil {
		return fmt.Errorf("invalid event: %w", err)
	}
	if busy, _ := s.IsEventTimeBusy(ctx, event); busy {
		return errors.New("event time is busy")
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO events ("title", "description", "begin", "end", "user_id", "created_at") 
				values ($1, $2, $3, $4, $5, $6)`,
		event.Title,
		event.Desc,
		event.Begin,
		event.End,
		event.UserID,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to create: %w", err)
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID int64) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM events WHERE id=$1", eventID)
	return err
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	if err := storage.ValidateTitle(event.Title); err != nil {
		return fmt.Errorf("invalid event: %w", err)
	}
	if busy, _ := s.IsEventTimeBusy(ctx, event); busy {
		return errors.New("event time is busy")
	}

	_, err := s.db.ExecContext(ctx,
		`UPDATE events SET "title"=$1, "description"=$2, "begin"=$3, "end"=$4, "updated_at"=$5 WHERE id=$6`,
		event.Title,
		event.Desc,
		event.Begin,
		event.End,
		time.Now(),
		event.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}
	return nil
}

func (s *Storage) GetAllEventsOfUser(ctx context.Context, userID int64) ([]storage.Event, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT "id", "title", "description", "begin", "end", "user_id", "created_at", "updated_at"
		FROM events WHERE "user_id" = $1`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve events of user: %w", err)
	}
	defer rows.Close()

	return scanEvents(rows)
}

func (s *Storage) IsEventTimeBusy(ctx context.Context, event storage.Event) (bool, error) {
	row := s.db.QueryRowContext(ctx, `SELECT * FROM "events" WHERE
      "user_id" = $1 AND "begin" <= $2 AND "end" >= $3`, event.UserID, event.End, event.Begin)

	var id int64
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, err
}

func scanEvents(rows *sql.Rows) ([]storage.Event, error) {
	var events []storage.Event
	for rows.Next() {
		var event storage.Event
		var updatedAt sql.NullTime
		var createdAt sql.NullTime

		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Desc,
			&event.Begin,
			&event.End,
			&event.UserID,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("cannot scan events: %w", err)
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
