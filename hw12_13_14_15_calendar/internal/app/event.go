package app

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Title  string
	Desc   string
	Begin  time.Time
	End    time.Time
}
