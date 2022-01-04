package app

import "time"

type Event struct {
	ID     int64
	Title  string
	Desc   string
	UserID int64
	Begin  time.Time
	End    time.Time
}
