package storage

import "time"

type Event struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"userID"`
	Title     string    `db:"title"`
	Desc      string    `db:"desc"`
	Begin     time.Time `db:"begin"`
	End       time.Time `db:"end"`
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`
}
