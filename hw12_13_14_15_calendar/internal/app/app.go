package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type App struct {
	logger  Logger
	storage Storage
}

type Storage interface {
	Connect(ctx context.Context) error
	Close() error
	CreateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, eventID uuid.UUID) error
	UpdateEvent(ctx context.Context, event storage.Event) error
	GetAllEventsOfUser(ctx context.Context, userID uuid.UUID) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event Event) error {
	return a.storage.CreateEvent(ctx, storage.Event{
		ID:     event.ID,
		Title:  event.Title,
		Desc:   event.Desc,
		Begin:  event.Begin,
		End:    event.End,
		UserID: event.UserID,
	})
}

func (a *App) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	return a.storage.DeleteEvent(ctx, eventID)
}

func (a *App) UpdateEvent(ctx context.Context, event Event) error {
	return a.storage.UpdateEvent(ctx, storage.Event{
		ID:    event.ID,
		Title: event.Title,
		Desc:  event.Desc,
		Begin: event.Begin,
		End:   event.End,
	})
}

func (a *App) GetAllEventsOfUser(ctx context.Context, userID uuid.UUID) ([]Event, error) {
	eventsInStorage, err := a.storage.GetAllEventsOfUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	var events []Event
	for _, event := range eventsInStorage {
		select {
		case <-ctx.Done():
			return nil, nil
		default:
			events = append(events, Event{
				ID:    event.ID,
				Title: event.Title,
				Desc:  event.Desc,
				Begin: event.Begin,
				End:   event.End,
			})
		}
	}
	return events, nil
}
