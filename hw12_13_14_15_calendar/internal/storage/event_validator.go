package storage

import (
	"errors"
	"time"
)

var (
	ErrInvalidEventTitle  = errors.New("invalid event title")
	ErrInvalidEventTime   = errors.New("invalid event time")
	ErrInvalidEventUserID = errors.New("invalid event user id")
)

func ValidateEventTime(eventBegin time.Time, eventEnd time.Time) error {
	if eventBegin.IsZero() || eventEnd.IsZero() {
		return ErrInvalidEventTime
	}
	if eventBegin.After(eventEnd) || eventBegin.Equal(eventEnd) {
		return ErrInvalidEventTime
	}
	return nil
}

func ValidateTitle(eventTitle string) error {
	if eventTitle == "" {
		return ErrInvalidEventTitle
	}
	return nil
}

func ValidateUserID(eventUserID int64) error {
	if eventUserID == 0 {
		return ErrInvalidEventUserID
	}
	return nil
}

func ValidateFull(event Event) error {
	if err := ValidateTitle(event.Title); err != nil {
		return err
	}
	if err := ValidateUserID(event.UserID); err != nil {
		return err
	}
	return ValidateEventTime(event.Begin, event.End)
}
