package app

import "time"

type Event struct {
	ID        string
	Title     string
	Desc      string
	EventTime time.Time
	Duration  time.Duration
	UserID    string
}
