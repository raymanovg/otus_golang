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
			if !s.IsEventTimeBusy(e, event) {
				return ErrEventTimeBusy
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

	s.mu.Lock()
	defer s.mu.Unlock()

	for i, e := range s.events {
		select {
		case <-ctx.Done():
			return nil
		default:
			if event.ID != e.ID {
				continue
			}
			if !s.IsEventTimeBusy(e, event) {
				return ErrEventTimeBusy
			}

			s.events[i].Title = event.Title
			s.events[i].Desc = event.Desc
			s.events[i].Begin = event.Begin
			s.events[i].End = event.End
			s.events[i].UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("not found")
}

func (s *Storage) GetAllEventsOfUser(ctx context.Context, userID int64) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
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
	s.mu.RLock()
	defer s.mu.RUnlock()
	select {
	case <-ctx.Done():
		return nil, nil
	default:
		return s.events, nil
	}
}

func (s *Storage) IsEventTimeBusy(savedEvent storage.Event, newEvent storage.Event) bool {
	begin := newEvent.Begin
	end := newEvent.End

	if savedEvent.UserID != newEvent.UserID || end.Before(savedEvent.Begin) {
		return true
	}
	if !begin.After(savedEvent.End) {
		return false
	}
	return true
}

func New() *Storage {
	return &Storage{}
}
