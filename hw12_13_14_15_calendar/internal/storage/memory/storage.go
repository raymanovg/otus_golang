package memorystorage

import (
	"context"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"sync"
)

type Storage struct {
	// TODO
	mu sync.RWMutex
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	return nil
}

func New() *Storage {
	return &Storage{}
}

// TODO
