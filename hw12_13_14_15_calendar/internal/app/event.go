package app

import "time"

type Event struct {
	ID       string
	Title    string
	Desc     string
	Time     time.Time
	Duration time.Duration
	UserID   string
}
