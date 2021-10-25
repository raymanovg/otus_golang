package hw06pipelineexecution

import (
	"github.com/stretchr/testify/require"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

func TestPipeline(t *testing.T) {
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	goID := func() int64 {
		var buf [64]byte
		n := runtime.Stack(buf[:], false)
		idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
		id, _ := strconv.ParseInt(idField, 10, 64)

		return id
	}

	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, []string{"102", "104", "106", "108", "110"}, result)
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})

	t.Run("execution stage concurrently", func(t *testing.T) {
		in := make(Bi)

		go func() {
			for _, v := range []int{1, 2, 3, 4, 5} {
				in <- v
			}
			close(in)
		}()

		stagesCount := 5
		goroutinesCount := 5

		runGoCount := make(map[int64]struct{}, stagesCount)

		mu := sync.Mutex{}

		stageFunc := func(v interface{}) interface{} {
			mu.Lock()
			runGoCount[goID()] = struct{}{}
			mu.Unlock()
			return v
		}

		var stages []Stage
		for i := 1; i <= stagesCount; i++ {
			stages = append(stages, g(strconv.Itoa(i), stageFunc))
		}

		for range ExecutePipeline(in, nil, stages...) {
		}

		require.Len(t, runGoCount, goroutinesCount, "not every stage in a separate goroutine")
	})

	t.Run("empty stage", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)

		result := ExecutePipeline(in, done, []Stage{}...)
		require.Nil(t, result, "if stage empty must return nil")
	})
}
