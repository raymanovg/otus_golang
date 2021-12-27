package app

import (
	"context"
	"go.uber.org/zap"

	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  *zap.Logger
	storage Storage
}

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, eventID int64) error
	UpdateEvent(ctx context.Context, event storage.Event) error
	GetAllEvents(ctx context.Context, userID int64) ([]storage.Event, error)
}

func New(logger *zap.Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event Event) error {
	return a.storage.CreateEvent(ctx, storage.Event{
		ID:       event.ID,
		Title:    event.Title,
		Desc:     event.Desc,
		Time:     event.Time,
		Duration: event.Duration,
		UserID:   event.UserID,
	})
}

func (a *App) DeleteEvent(ctx context.Context, eventID int64) error {
	return nil
}

func (a *App) UpdateEvent(ctx context.Context, event Event) error {
	return nil
}

func (a *App) GetAllEvents(ctx context.Context, userID int64) ([]Event, error) {
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
