package sqlstorage

import (
	"context"

	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct{}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
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
	return nil, nil
}

func (s *Storage) IsEventTimeBusy(ctx context.Context, event storage.Event) bool {
	return false
}
