package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrEventIDUsed   = errors.New("event with the id is already exist")
	ErrEventTimeBusy = errors.New("event time is busy")
)

type Storage struct {
	events []storage.Event
	mu     sync.RWMutex
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range s.events {
		select {
		case <-ctx.Done():
			return nil
		default:
			if e.ID == event.ID {
				return ErrEventIDUsed
			}
			if e.UserID == event.UserID && e.Time == event.Time {
				return ErrEventTimeBusy
			}
		}
	}

	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	s.events = append(s.events, event)

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID int64) error {
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	return nil
}

func (s *Storage) GetAllEvents(ctx context.Context, userID int64) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var events []storage.Event
	for _, e := range s.events {
		select {
		case <-ctx.Done():
			return nil, nil
		default:
			if userID == e.UserID {
				events = append(events, e)
			}
		}

	}
	return events, nil
}

func New() *Storage {
	return &Storage{}
}
