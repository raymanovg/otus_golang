package memorystorage

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/config"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestCreateEvent(t *testing.T) {
	type Case struct {
		name  string
		event storage.Event
		err   error
	}
	t.Run("success event create", func(t *testing.T) {
		userID := uuid.New()
		fmt.Println(userID.String())
		cases := []Case{
			{
				name: "success event create one",
				event: storage.Event{
					Title:  "event one",
					Desc:   "event one",
					Begin:  time.Now(),
					End:    time.Now().Add(time.Hour),
					UserID: userID,
				},
				err: nil,
			},
			{
				name: "success event create two",
				event: storage.Event{
					Title:  "event two",
					Desc:   "event two",
					Begin:  time.Now().Add(2 * time.Hour),
					End:    time.Now().Add(3 * time.Hour),
					UserID: userID,
				},
				err: nil,
			},
			{
				name: "success event create three",
				event: storage.Event{
					Title:  "event three",
					Desc:   "event three",
					Begin:  time.Now(),
					End:    time.Now().Add(time.Hour),
					UserID: uuid.New(),
				},
				err: nil,
			},
		}

		st := New(config.Memory{})
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				err := st.CreateEvent(context.Background(), c.event)
				require.ErrorIs(t, err, c.err)
			})
		}
	})

	t.Run("invalid event data", func(t *testing.T) {
		cases := []Case{
			{
				name: "invalid title",
				event: storage.Event{
					Title:  "",
					Desc:   "event one",
					Begin:  time.Now(),
					End:    time.Now().Add(time.Hour),
					UserID: uuid.New(),
				},
				err: storage.ErrInvalidEventTitle,
			},
			{
				name: "empty begin time of event",
				event: storage.Event{
					Title:  "event one",
					Desc:   "event desc",
					End:    time.Now().Add(time.Hour),
					UserID: uuid.New(),
				},
				err: storage.ErrInvalidEventTime,
			},
			{
				name: "empty end time of event",
				event: storage.Event{
					Title:  "event one",
					Desc:   "event desc",
					Begin:  time.Now(),
					UserID: uuid.New(),
				},
				err: storage.ErrInvalidEventTime,
			},
		}
		st := New(config.Memory{})
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				err := st.CreateEvent(context.Background(), c.event)
				require.ErrorIs(t, err, c.err)
			})
		}
	})

	t.Run("event time busy", func(t *testing.T) {
		userID := uuid.New()
		st := New(config.Memory{})
		err := st.CreateEvent(context.Background(), storage.Event{
			Title:  "event one",
			Desc:   "event one",
			Begin:  time.Date(2021, 12, 28, 16, 0, 0, 0, time.Local),
			End:    time.Date(2021, 12, 28, 17, 0, 0, 0, time.Local),
			UserID: userID,
		})
		require.NoError(t, err)

		err = st.CreateEvent(context.Background(), storage.Event{
			Title:  "event two",
			Desc:   "event two",
			Begin:  time.Date(2021, 12, 28, 16, 30, 0, 0, time.Local),
			End:    time.Date(2021, 12, 28, 17, 30, 0, 0, time.Local),
			UserID: userID,
		})
		require.ErrorIs(t, err, ErrEventTimeBusy)

		err = st.CreateEvent(context.Background(), storage.Event{
			Title:  "event three",
			Desc:   "event three",
			Begin:  time.Date(2021, 12, 28, 15, 30, 0, 0, time.Local),
			End:    time.Date(2021, 12, 28, 16, 30, 0, 0, time.Local),
			UserID: userID,
		})
		require.ErrorIs(t, err, ErrEventTimeBusy)
	})
}

func TestGetAllEventsOfUser(t *testing.T) {
	st := New(config.Memory{})
	total := 100

	userIDs := []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.New(),
	}
	count := make(map[uuid.UUID]int)
	for i := 1; i <= total; i++ {
		userIndex := i % 5
		err := st.CreateEvent(context.Background(), storage.Event{
			Title:  "event " + strconv.Itoa(i),
			Desc:   "event desc of " + strconv.Itoa(i),
			Begin:  time.Now(),
			End:    time.Now().Add(time.Nanosecond),
			UserID: userIDs[userIndex],
		})
		count[userIDs[userIndex]]++
		require.NoError(t, err)
	}

	for userID, c := range count {
		events, err := st.GetAllEventsOfUser(context.Background(), userID)
		require.NoError(t, err)
		require.Len(t, events, c)
	}
}

func TestDeleteEvent(t *testing.T) {
	st := New(config.Memory{})
	require.ErrorIs(t, st.DeleteEvent(context.Background(), uuid.New()), ErrEventNotFound)

	userID := uuid.New()
	err := st.CreateEvent(context.Background(), storage.Event{
		Title:  "event 1",
		Desc:   "event desc of 1",
		Begin:  time.Now(),
		End:    time.Now().Add(time.Hour),
		UserID: userID,
	})
	require.NoError(t, err)

	events, err := st.GetAllEventsOfUser(context.Background(), userID)
	require.NoError(t, err)
	require.Len(t, events, 1)

	require.NoError(t, st.DeleteEvent(context.Background(), events[0].ID))
	require.ErrorIs(t, st.DeleteEvent(context.Background(), events[0].ID), ErrEventNotFound)
}

func TestUpdateEvent(t *testing.T) {
	userID := uuid.New()
	st := New(config.Memory{})
	err := st.CreateEvent(context.Background(), storage.Event{
		Title:  "event 1",
		Desc:   "event desc of 1",
		Begin:  time.Now(),
		End:    time.Now().Add(time.Hour),
		UserID: userID,
	})
	require.NoError(t, err)
	err = st.CreateEvent(context.Background(), storage.Event{
		Title:  "event 2",
		Desc:   "event desc of 2",
		Begin:  time.Now().Add(time.Hour),
		End:    time.Now().Add(2 * time.Hour),
		UserID: userID,
	})
	require.NoError(t, err)

	events, err := st.GetAllEventsOfUser(context.Background(), userID)
	require.NoError(t, err)
	require.Len(t, events, 2)

	for _, event := range events {
		require.Equal(t, false, event.UpdatedAt.IsZero())
	}

	eventTwoCreatedAt := events[1].CreatedAt
	eventTwoBegin := time.Now().Add(2 * time.Hour)
	eventTwoEnd := time.Now().Add(3 * time.Hour)
	err = st.UpdateEvent(context.Background(), storage.Event{
		ID:     events[1].ID,
		Title:  "event 2 updated",
		Desc:   "event desc of 2 updated",
		Begin:  eventTwoBegin,
		End:    eventTwoEnd,
		UserID: userID,
	})
	require.NoError(t, err)

	event, err := st.GetEvent(context.Background(), events[1].ID)
	require.NoError(t, err)
	require.Equal(t, eventTwoCreatedAt, event.CreatedAt)
	require.Equal(t, false, event.UpdatedAt.IsZero())
	require.Equal(t, "event 2 updated", event.Title)
	require.Equal(t, "event desc of 2 updated", event.Desc)
	require.Equal(t, eventTwoBegin, event.Begin)
	require.Equal(t, eventTwoEnd, event.End)
	require.Equal(t, userID, event.UserID)

	err = st.UpdateEvent(context.Background(), storage.Event{
		ID:     events[1].ID,
		Title:  "event 2 updated",
		Desc:   "event desc of 2 updated",
		Begin:  time.Now(),
		End:    time.Now().Add(time.Hour),
		UserID: userID,
	})
	require.ErrorIs(t, err, ErrEventTimeBusy)
}
