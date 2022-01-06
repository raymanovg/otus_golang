package storage

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"userID"`
	Title     string    `db:"title"`
	Desc      string    `db:"desc"`
	Begin     time.Time `db:"begin"`
	End       time.Time `db:"end"`
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`
}
