package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	goID := func() int64 {
		var buf [64]byte
		n := runtime.Stack(buf[:], false)
		idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
		id, _ := strconv.ParseInt(idField, 10, 64)
		return id
	}

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 7
		maxErrorsCount := 1

		require.Eventually(t, func() bool {
			err := Run(tasks, workersCount, maxErrorsCount)
			require.NoError(t, err)

			return int32(tasksCount) == runTasksCount
		}, time.Second, time.Millisecond)
	})

	t.Run("working concurrently", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		workersCount := 5
		maxErrorsCount := 0

		runGoList := make(map[int64]struct{}, workersCount)
		mu := sync.Mutex{}

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				mu.Lock()
				runGoList[goID()] = struct{}{}
				mu.Unlock()
				return nil
			})
		}

		require.Eventually(t, func() bool {
			err := Run(tasks, workersCount, maxErrorsCount)
			require.NoError(t, err)

			return runTasksCount == int32(tasksCount) && len(runGoList) == workersCount
		}, time.Second, time.Millisecond)
	})

	t.Run("zero or negative max errors count", func(t *testing.T) {
		cases := []struct {
			name        string
			maxErrCount int
			err         error
		}{
			{"zero max errors count", 0, nil},
			{"-1 max errors count", -1, nil},
			{"-10 max errors count", -10, nil},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				tasksCount := 10
				workersCount := 5

				errorTasksCount := rand.Intn(tasksCount)

				tasks := make([]Task, 0, tasksCount)

				for i := 0; i < errorTasksCount; i++ {
					tasks = append(tasks, func() error {
						return errors.New("error from task")
					})
				}
				for i := 0; i < tasksCount-errorTasksCount; i++ {
					tasks = append(tasks, func() error {
						return nil
					})
				}

				rand.Seed(time.Now().UnixNano())
				rand.Shuffle(len(tasks), func(i, j int) {
					tasks[i], tasks[j] = tasks[j], tasks[i]
				})

				require.Equal(t, Run(tasks, workersCount, c.maxErrCount), c.err)
			})
		}
	})
}
