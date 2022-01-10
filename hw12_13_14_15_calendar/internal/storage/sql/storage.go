package sqlstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrEventTimeBusy = errors.New("event time is busy")
	ErrEventNotFound = errors.New("event not found")
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
	busy, err := s.checkEventTime(ctx, event)
	if err != nil {
		return fmt.Errorf("failed to check event time: %w", err)
	}
	if busy {
		return ErrEventTimeBusy
	}

	query := `INSERT INTO events ("id", "title", "desc", "begin", "end", "userID") 
			  VALUES (:id, :title, :desc, :begin, :end, :userID)
	`
	args := map[string]interface{}{
		"id":     event.ID.String(),
		"title":  event.Title,
		"desc":   event.Desc,
		"begin":  event.Begin,
		"end":    event.End,
		"userID": event.UserID,
	}

	if _, err := s.db.NamedExecContext(ctx, query, args); err != nil {
		return fmt.Errorf("failed to create: %w", err)
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	query := `DELETE FROM events WHERE id = :id`
	args := map[string]interface{}{"id": eventID}
	_, err := s.db.NamedExecContext(ctx, query, args)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}

	return err
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	if err := storage.ValidateTitle(event.Title); err != nil {
		return fmt.Errorf("invalid event: %w", err)
	}
	if err := storage.ValidateEventTime(event.Begin, event.End); err != nil {
		return err
	}
	if busy, _ := s.checkEventTime(ctx, event); busy {
		return ErrEventTimeBusy
	}

	query := `UPDATE events SET "title" = :title, "desc" = :desc, "begin" = :begin, "end" = :end WHERE id = :id`
	args := map[string]interface{}{
		"title": event.Title,
		"desc":  event.Desc,
		"begin": event.Begin,
		"end":   event.End,
		"id":    event.ID,
	}

	res, err := s.db.NamedExecContext(ctx, query, args)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}
	if affected == 0 {
		return ErrEventNotFound
	}

	return nil
}

func (s *Storage) GetAllEventsOfUser(ctx context.Context, userID uuid.UUID) ([]storage.Event, error) {
	query := `SELECT * FROM events WHERE "userID" = :userID`
	args := map[string]interface{}{"userID": userID}

	rows, err := s.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve events of user %d: %w", userID, err)
	}
	defer rows.Close()

	return scanEvents(rows)
}

func (s *Storage) checkEventTime(ctx context.Context, event storage.Event) (bool, error) {
	query := `SELECT "id" FROM events WHERE "userID" = :userID AND "begin" <= :end AND "end" >= :begin LIMIT 1`
	args := map[string]interface{}{
		"userID": event.UserID.String(),
		"end":    event.End,
		"begin":  event.Begin,
	}

	rows, err := s.db.NamedQueryContext(ctx, query, args)
	if err != nil {
		return true, err
	}

	defer rows.Close()

	return rows.Next(), nil
}

func scanEvents(rows *sqlx.Rows) ([]storage.Event, error) {
	var events []storage.Event
	for rows.Next() {
		var event storage.Event
		if err := rows.StructScan(&event); err != nil {
			return nil, fmt.Errorf("cannot scan events: %w", err)
		}
		events = append(events, event)
	}
	return events, nil
}
