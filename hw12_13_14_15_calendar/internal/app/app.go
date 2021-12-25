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
}

func New(logger Logger, storage Storage) *App {
	return &App{
		log:     logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	return a.storage.CreateEvent(ctx, storage.Event{ID: id, Title: title})
}

// TODO
