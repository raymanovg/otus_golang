package memorystorage

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage"
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

		st := New()
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
				name: "invalid desc",
				event: storage.Event{
					Title:  "event one",
					Desc:   "",
					Begin:  time.Now(),
					End:    time.Now().Add(time.Hour),
					UserID: 2,
				},
				err: storage.ErrInvalidEventDesc,
			},
			{
				name: "invalid time",
				event: storage.Event{
					Title:  "event one",
					Desc:   "event desc",
					End:    time.Now().Add(time.Hour),
					UserID: 3,
				},
				err: storage.ErrInvalidEventBeginTime,
			},
			{
				name: "invalid duration",
				event: storage.Event{
					Title:  "event one",
					Desc:   "event desc",
					Begin:  time.Now(),
					UserID: 4,
				},
				err: storage.ErrInvalidEventEndTime,
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
		st := New()
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				err := st.CreateEvent(context.Background(), c.event)
				require.ErrorIs(t, err, c.err)
			})
		}
	})

	t.Run("event time busy", func(t *testing.T) {
		st := New()
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

func TestGetAllEvents(t *testing.T) {
	st := New()
	total := 100
	for i := 1; i <= total; i++ {
		err := st.CreateEvent(context.Background(), storage.Event{
			Title:  "event title",
			Desc:   "event desc",
			Begin:  time.Now(),
			End:    time.Now().Add(1 * time.Hour),
			UserID: int64(i),
		})
		require.NoError(t, err)
	}

	events, err := st.GetAllEvents(context.Background())
	require.NoError(t, err)
	require.Len(t, events, total)
	require.Equal(t, int64(1), events[0].ID, "first event id must be 1")
	require.Equal(t, int64(total), events[len(events)-1].ID, "last event id must be "+strconv.Itoa(total))
	require.Equal(t, int64(1), events[0].UserID, "first event user id must be 1")
	require.Equal(t, int64(total), events[len(events)-1].UserID, "last event user id must be "+strconv.Itoa(total))
}
