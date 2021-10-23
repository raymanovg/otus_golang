package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	results := make(chan error)
	errCount, done := observe(results, m)
	taskStream := generator(tasks, done)

	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			worker(taskStream, results, done)
			wg.Done()
		}()
	}

	wg.Wait()

	close(results)

	c := <-errCount
	if m >= 1 && c >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func worker(tasks <-chan Task, results chan<- error, done <-chan struct{}) {
	for task := range tasks {
		select {
		case <-done:
			return
		case results <- task():
		}
	}
}

func observe(results <-chan error, m int) (<-chan int, <-chan struct{}) {
	done := make(chan struct{})
	errCount := make(chan int, 1)

	go func() {
		defer close(done)
		defer close(errCount)

		var c int
		for err := range results {
			if err != nil {
				c++
			}
			if c == m && m != 0 {
				errCount <- c
				break
			}
		}
	}()

	return errCount, done
}

func generator(tasks []Task, done <-chan struct{}) <-chan Task {
	tasksCh := make(chan Task)
	go func() {
		defer close(tasksCh)
		for _, task := range tasks {
			select {
			case <-done:
				return
			case tasksCh <- task:
			}
		}
	}()

	return tasksCh
}
