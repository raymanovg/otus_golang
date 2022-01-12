package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrEventTimeBusy = errors.New("event time is busy")
	ErrEventNotFound = errors.New("event not found")
)

type Storage struct {
	config config.Memory
	events map[uuid.UUID]map[uuid.UUID]storage.Event
	mu     sync.RWMutex
}

func (s *Storage) Connect(ctx context.Context) error {
	return nil
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	if err := storage.ValidateFull(event); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	userEvents, ok := s.events[event.UserID]
	if !ok {
		s.events[event.UserID] = make(map[uuid.UUID]storage.Event)
	}

	event.ID = uuid.New()
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	for _, e := range userEvents {
		select {
		case <-ctx.Done():
			return nil
		default:
			if !s.IsEventTimeBusy(e, event) {
				return ErrEventTimeBusy
			}
		}
	}

	s.events[event.UserID][event.ID] = event

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, userEvents := range s.events {
		if _, ok := userEvents[eventID]; ok {
			select {
			case <-ctx.Done():
				return nil
			default:
				delete(userEvents, eventID)
				return nil
			}
		}
	}
	return ErrEventNotFound
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	if err := storage.ValidateTitle(event.Title); err != nil {
		return err
	}
	if err := storage.ValidateEventTime(event.Begin, event.End); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	userEvents, ok := s.events[event.UserID]
	if !ok {
		return ErrEventNotFound
	}
	savedEvent, ok := userEvents[event.ID]
	if !ok {
		return ErrEventNotFound
	}

	for _, e := range userEvents {
		select {
		case <-ctx.Done():
			return nil
		default:
			if e.ID != event.ID && !s.IsEventTimeBusy(e, event) {
				return ErrEventTimeBusy
			}
		}
	}

	savedEvent.Title = event.Title
	savedEvent.Desc = event.Desc
	savedEvent.Begin = event.Begin
	savedEvent.End = event.End
	savedEvent.UpdatedAt = time.Now()

	s.events[event.UserID][event.ID] = savedEvent

	return nil
}

func (s *Storage) GetAllEventsOfUser(ctx context.Context, userID uuid.UUID) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var events []storage.Event
	userEvents, ok := s.events[userID]
	if !ok {
		return events, nil
	}
	for _, e := range userEvents {
		select {
		case <-ctx.Done():
			return nil, nil
		default:
			events = append(events, e)
		}
	}
	return events, nil
}

func (s *Storage) GetEvent(ctx context.Context, eventID uuid.UUID) (storage.Event, error) {
	for _, userEvents := range s.events {
		select {
		case <-ctx.Done():
			return storage.Event{}, nil
		default:
			if event, ok := userEvents[eventID]; ok {
				return event, nil
			}
		}
	}
	return storage.Event{}, ErrEventNotFound
}

func (s *Storage) IsEventTimeBusy(savedEvent storage.Event, newEvent storage.Event) bool {
	if savedEvent.UserID != newEvent.UserID || newEvent.End.Before(savedEvent.Begin) {
		return true
	}
	if !newEvent.Begin.After(savedEvent.End) {
		return false
	}
	return true
}

func New(conf config.Memory) *Storage {
	return &Storage{
		config: conf,
		events: make(map[uuid.UUID]map[uuid.UUID]storage.Event),
	}
}
