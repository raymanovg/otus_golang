package storage

import "time"

type Event struct {
	ID        int64
	Title     string
	Desc      string
	Begin     time.Time
	End       time.Time
	UserID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
