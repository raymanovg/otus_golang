package storage

import "errors"

var (
	ErrInvalidEventTitle     = errors.New("invalid event title")
	ErrInvalidEventDesc      = errors.New("invalid event description")
	ErrInvalidEventBeginTime = errors.New("invalid event begin time")
	ErrInvalidEventEndTime   = errors.New("invalid event end time")
	ErrInvalidEventUserID    = errors.New("invalid event user id")
)

func Validate(event Event) error {
	if event.Title == "" {
		return ErrInvalidEventTitle
	}
	if event.Desc == "" {
		return ErrInvalidEventDesc
	}
	if event.Begin.IsZero() {
		return ErrInvalidEventBeginTime
	}
	if event.End.IsZero() {
		return ErrInvalidEventEndTime
	}
	if event.UserID == 0 {
		return ErrInvalidEventUserID
	}

	return nil
}
