package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrEmptyTasksList      = errors.New("list of tasks is empty")
	ErrMinGoroutinesNumber = errors.New("number of goroutines must be posiive")
	ErrMinErrorsNumber     = errors.New("number of errors can not be negative")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return ErrEmptyTasksList
	}
	if n <= 0 {
		return ErrMinGoroutinesNumber
	}
	if m < 0 {
		return ErrMinErrorsNumber
	}

	tasksCh := make(chan Task, len(tasks))
	resultCh := make(chan error)
	stopCh := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(n)

	for _, task := range tasks {
		tasksCh <- task
	}
	close(tasksCh)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			worker(tasksCh, resultCh, stopCh)
		}()
	}

	go func() {
		defer close(resultCh)
		wg.Wait()
	}()

	errorsCounter := 0
	for err := range resultCh {
		if err != nil {
			errorsCounter++
		}
		if errorsCounter == m {
			close(stopCh)
		}
	}

	if errorsCounter >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func worker(tasks <-chan Task, results chan<- error, stop <-chan struct{}) {
	for task := range tasks {
		select {
		case <-stop:
			return
		default:
			results <- task()
		}
	}
}
