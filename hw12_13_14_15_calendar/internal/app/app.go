package app

import (
	"context"

	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	log     Logger
	storage Storage
}

type Logger interface {
	Info(string)
	Warn(string)
	Error(string)
	Debug(string)
}

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	GetAllEvents(ctx context.Context, userID string) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		log:     logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event Event) error {
	return a.storage.CreateEvent(ctx, storage.Event{
		ID:        event.ID,
		Title:     event.Title,
		Desc:      event.Desc,
		EventTime: event.EventTime,
		Duration:  event.Duration,
		UserID:    event.UserID,
	})
}

func (a *App) GetAllEvents(ctx context.Context, userID string) ([]Event, error) {
	eventsInStorage, err := a.storage.GetAllEvents(ctx, userID)
	if err != nil {
		return nil, err
	}
	var events []Event
	for _, event := range eventsInStorage {
		select {
		case <-ctx.Done():
			return nil, nil
		default:
			events = append(events, Event{ID: event.ID})
		}
	}
	return events, nil
}
