package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	tasksCh := makeTasksChanel(tasks)
	errCh := make(chan error, n)
	interruptCh, errCountChan := observe(errCh, m)

	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			worker(tasksCh, errCh, interruptCh)
			wg.Done()
		}()
	}

	wg.Wait()

	close(errCh)

	errCount := <-errCountChan
	if m != 0 && errCount >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func worker(tasks <-chan Task, errors chan<- error, interrupt <-chan struct{}) {
	for {
		select {
		case <-interrupt:
			return
		default:
		}
		select {
		case task, ok := <-tasks:
			if !ok {
				return
			}
			errors <- task()
		default:
		}
	}
}

func observe(errCh <-chan error, m int) (<-chan struct{}, <-chan int) {
	interruptCh := make(chan struct{})
	errCountCh := make(chan int, 1)
	go func() {
		var count int
		for err := range errCh {
			if err != nil {
				count++
			}
			if count == m && m != 0 {
				errCountCh <- count
				break
			}
		}
		close(interruptCh)
		close(errCountCh)
	}()

	return interruptCh, errCountCh
}

func makeTasksChanel(tasks []Task) <-chan Task {
	tasksCh := make(chan Task, len(tasks))
	for _, task := range tasks {
		tasksCh <- task
	}
	close(tasksCh)
	return tasksCh
}
