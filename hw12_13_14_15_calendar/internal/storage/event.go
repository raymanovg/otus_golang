package storage

import "time"

type Event struct {
	ID        string
	Title     string
	Desc      string
	Time      time.Time
	Duration  time.Duration
	UserID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
