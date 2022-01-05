package memorystorage

import (
	"context"
	"strconv"
	"testing"
	"time"

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
		cases := []Case{
			{
				name: "success event create one",
				event: storage.Event{
					Title:  "event one",
					Desc:   "event one",
					Begin:  time.Now(),
					End:    time.Now().Add(time.Hour),
					UserID: 1,
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
					UserID: 1,
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
					UserID: 2,
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
					UserID: 1,
				},
				err: storage.ErrInvalidEventTitle,
			},
			{
				name: "invalid time",
				event: storage.Event{
					Title:  "event one",
					Desc:   "event desc",
					End:    time.Now().Add(time.Hour),
					UserID: 3,
				},
				err: storage.ErrInvalidEventTime,
			},
			{
				name: "invalid duration",
				event: storage.Event{
					Title:  "event one",
					Desc:   "event desc",
					Begin:  time.Now(),
					UserID: 4,
				},
				err: storage.ErrInvalidEventTime,
			},
			{
				name: "invalid user id",
				event: storage.Event{
					Title: "event one",
					Desc:  "event desc",
					Begin: time.Now(),
					End:   time.Now().Add(time.Hour),
				},
				err: storage.ErrInvalidEventUserID,
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
		st := New(config.Memory{})
		err := st.CreateEvent(context.Background(), storage.Event{
			Title:  "event one",
			Desc:   "event one",
			Begin:  time.Date(2021, 12, 28, 16, 0, 0, 0, time.Local),
			End:    time.Date(2021, 12, 28, 17, 0, 0, 0, time.Local),
			UserID: 1,
		})
		require.NoError(t, err)

		err = st.CreateEvent(context.Background(), storage.Event{
			Title:  "event two",
			Desc:   "event two",
			Begin:  time.Date(2021, 12, 28, 16, 30, 0, 0, time.Local),
			End:    time.Date(2021, 12, 28, 17, 30, 0, 0, time.Local),
			UserID: 1,
		})
		require.ErrorIs(t, err, ErrEventTimeBusy)

		err = st.CreateEvent(context.Background(), storage.Event{
			Title:  "event three",
			Desc:   "event three",
			Begin:  time.Date(2021, 12, 28, 15, 30, 0, 0, time.Local),
			End:    time.Date(2021, 12, 28, 16, 30, 0, 0, time.Local),
			UserID: 1,
		})
		require.ErrorIs(t, err, ErrEventTimeBusy)
	})
}

func TestGetAllEventsOfUser(t *testing.T) {
	st := New(config.Memory{})
	total := 100

	count := make(map[int]int)
	for i := 1; i <= total; i++ {
		userID := (i % 5) + 1
		err := st.CreateEvent(context.Background(), storage.Event{
			Title:  "event " + strconv.Itoa(i),
			Desc:   "event desc of " + strconv.Itoa(i),
			Begin:  time.Now(),
			End:    time.Now().Add(time.Nanosecond),
			UserID: int64(userID),
		})
		count[userID]++
		require.NoError(t, err)
	}

	for userID, c := range count {
		events, err := st.GetAllEventsOfUser(context.Background(), int64(userID))
		require.NoError(t, err)
		require.Len(t, events, c)
	}
}

func TestDeleteEvent(t *testing.T) {
	st := New(config.Memory{})
	require.ErrorIs(t, st.DeleteEvent(context.Background(), 1), ErrEventNotFound)
	err := st.CreateEvent(context.Background(), storage.Event{
		Title:  "event 1",
		Desc:   "event desc of 1",
		Begin:  time.Now(),
		End:    time.Now().Add(time.Hour),
		UserID: 1,
	})
	require.NoError(t, err)
	require.NoError(t, st.DeleteEvent(context.Background(), 1))
	require.ErrorIs(t, st.DeleteEvent(context.Background(), 1), ErrEventNotFound)
}

func TestUpdateEvent(t *testing.T) {
	st := New(config.Memory{})
	err := st.CreateEvent(context.Background(), storage.Event{
		Title:  "event 1",
		Desc:   "event desc of 1",
		Begin:  time.Now(),
		End:    time.Now().Add(time.Hour),
		UserID: 1,
	})
	require.NoError(t, err)
	err = st.CreateEvent(context.Background(), storage.Event{
		Title:  "event 2",
		Desc:   "event desc of 2",
		Begin:  time.Now().Add(time.Hour),
		End:    time.Now().Add(2 * time.Hour),
		UserID: 1,
	})
	require.NoError(t, err)

	event, err := st.GetEvent(context.Background(), 2)
	require.NoError(t, err)
	require.Equal(t, true, event.UpdatedAt.IsZero())

	eventTwoCreatedAt := event.CreatedAt
	eventTwoBegin := time.Now().Add(2 * time.Hour)
	eventTwoEnd := time.Now().Add(3 * time.Hour)
	err = st.UpdateEvent(context.Background(), storage.Event{
		ID:     2,
		Title:  "event 2 updated",
		Desc:   "event desc of 2 updated",
		Begin:  eventTwoBegin,
		End:    eventTwoEnd,
		UserID: 1,
	})
	require.NoError(t, err)

	event, err = st.GetEvent(context.Background(), 2)
	require.NoError(t, err)
	require.Equal(t, eventTwoCreatedAt, event.CreatedAt)
	require.Equal(t, false, event.UpdatedAt.IsZero())
	require.Equal(t, "event 2 updated", event.Title)
	require.Equal(t, "event desc of 2 updated", event.Desc)
	require.Equal(t, eventTwoBegin, event.Begin)
	require.Equal(t, eventTwoEnd, event.End)
	require.Equal(t, int64(1), event.UserID)

	err = st.UpdateEvent(context.Background(), storage.Event{
		ID:     2,
		Title:  "event 2 updated",
		Desc:   "event desc of 2 updated",
		Begin:  time.Now(),
		End:    time.Now().Add(time.Hour),
		UserID: 1,
	})
	require.ErrorIs(t, err, ErrEventTimeBusy)
}
