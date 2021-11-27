package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProgressBar(t *testing.T) {
	t.Run("calculating percent", func(t *testing.T) {
		cases := []struct {
			name    string
			total   int64
			current int64
			percent int
		}{
			{"percent 500 of 1000", 1000, 0, 0},
			{"percent 10 of 1000", 1000, 10, 1},
			{"percent 500 of 1000", 1000, 500, 50},
			{"percent 500 of 1000", 1000, 1000, 100},
			{"percent 500 of 1000", 1000, 10000, 1000},
			{"percent 500 of 1000", 1000, 10000, 1000},
		}
		pb := NewProgressBar(0)
		for _, c := range cases {
			pb.SetTotal(c.total)
			pb.SetCurrent(c.current)
			require.Equal(t, c.percent, pb.getPercent(), "wrong percent")
		}
	})

	t.Run("templating output", func(t *testing.T) {
		cases := []struct {
			out     string
			current int64
			total   int64
		}{
			{"[%-100s] 10%%     10/1 bite", 1, 10},
			{"[%-100s] 20%%     10/2 bite", 2, 10},
			{"[%-100s] 30%%     10/3 bite", 3, 10},
			{"[%-100s] 40%%     10/4 bite", 4, 10},
			{"[%-100s] 50%%     10/5 bite", 5, 10},
			{"[%-100s] 60%%     10/6 bite", 6, 10},
			{"[%-100s] 70%%     10/7 bite", 7, 10},
			{"[%-100s] 80%%     10/8 bite", 8, 10},
			{"[%-100s] 90%%     10/9 bite", 9, 10},
			{"[%-100s]100%%     10/10 bite", 10, 10},
		}

		bf := bytes.Buffer{}
		pb := NewProgressBar(0)
		pb.SetOutput(&bf)

		for _, c := range cases {
			pb.SetTotal(c.total)
			pb.Start()
			pb.Add(c.current)
			require.Equal(
				t,
				fmt.Sprintf(c.out, strings.Repeat(">", pb.percent)),
				strings.TrimSpace(bf.String()),
				"pb output is not correct",
			)
			bf.Reset()
		}
	})

	t.Run("reset pb data", func(t *testing.T) {
		pb := NewProgressBar(1000)

		pb.Add(100)
		require.Equal(t, int64(100), pb.current, "current is not correct after 'add' func call")
		require.Equal(t, 10, pb.percent, "percent is not correct after 'add' func call")
		require.Equal(t, ">>>>>>>>>>", pb.rate, "rate is not correct for 10 percent")

		pb.Reset()
		require.Equal(t, int64(0), pb.current, "current must be 0 after reset")
		require.Equal(t, 0, pb.percent, "percent must be 0 after reset")
		require.Equal(t, "", pb.rate, "rate must be empty string after reset")

		pb.Add(10)
		pb.Start()
		require.Equal(t, int64(0), pb.current, "current data must be reset after start")
		require.Equal(t, 0, pb.percent, "percent must be reset start")
		require.Equal(t, "", pb.rate, "rate must be empty string after start")
	})
}
