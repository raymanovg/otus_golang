package storage

import "errors"

var (
	ErrInvalidEventTitle    = errors.New("invalid event title")
	ErrInvalidEventDesc     = errors.New("invalid event description")
	ErrInvalidEventTime     = errors.New("invalid event time")
	ErrInvalidEventDuration = errors.New("invalid event duration")
	ErrInvalidEventUserID   = errors.New("invalid event user id")
)

func Validate(event Event) error {
	if event.Title == "" {
		return ErrInvalidEventTitle
	}
	if event.Desc == "" {
		return ErrInvalidEventDesc
	}
	if event.Time.IsZero() {
		return ErrInvalidEventTime
	}
	if event.Duration == 0 {
		return ErrInvalidEventDuration
	}
	if event.UserID == 0 {
		return ErrInvalidEventUserID
	}

	return nil
}
