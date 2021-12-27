package app

import "time"

type Event struct {
	ID       string
	Title    string
	Desc     string
	UserID   int64
	Time     time.Time
	Duration time.Duration
}
