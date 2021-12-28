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
	lastEventId int64
	events      []storage.Event
	mu          sync.RWMutex
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	if err := storage.Validate(event); err != nil {
		return err
	}
	if !s.IsEventTimeBusy(event) {
		return ErrEventTimeBusy
	}

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
		}
	}

	event.ID = s.lastEventId + 1
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	s.events = append(s.events, event)
	s.lastEventId++

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, event := range s.events {
		select {
		case <-ctx.Done():
			return nil
		default:
			if eventID == event.ID {
				s.events = append(s.events[:i], s.events[i+1:]...)
				return nil
			}
		}
	}
	return errors.New("not found")
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	if err := storage.Validate(event); err != nil {
		return err
	}
	if !s.IsEventTimeBusy(event) {
		return ErrEventTimeBusy
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for i, e := range s.events {
		select {
		case <-ctx.Done():
			return nil
		default:
			if event.ID == e.ID {
				s.events[i].Title = event.Title
				s.events[i].Desc = event.Desc
				s.events[i].Time = event.Time
				s.events[i].Duration = event.Duration
				s.events[i].UpdatedAt = time.Now()
				return nil
			}
		}
	}
	return errors.New("not found")
}

func (s *Storage) GetAllEventsOfUser(ctx context.Context, userID int64) ([]storage.Event, error) {
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

func (s *Storage) GetAllEvents(ctx context.Context) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.events, nil
}

func (s *Storage) IsEventTimeBusy(event storage.Event) bool {
	begin := event.Time
	end := event.Time.Add(event.Duration)
	for _, e := range s.events {
		if e.UserID != event.UserID {
			continue
		}
		if end.Before(e.Time) {
			continue
		}
		if !begin.After(e.Time.Add(e.Duration)) {
			return false
		}
	}
	return true
}

func New() *Storage {
	return &Storage{}
}
