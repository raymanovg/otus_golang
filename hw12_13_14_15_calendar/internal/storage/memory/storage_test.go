package memorystorage

import (
	"context"
	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStorage(t *testing.T) {
	t.Run("create events", func(t *testing.T) {
		st := New()
		err := st.CreateEvent(context.Background(), storage.Event{})
		require.NoError(t, err)
	})
}
