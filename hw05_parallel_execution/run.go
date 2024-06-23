package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n int, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	errCount := int64(0)
	taskChan := make(chan Task)
	doneChan := make(chan struct{})
	wg := &sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go newWorker(m, taskChan, doneChan, &errCount, wg)
	}

	for _, task := range tasks {
		select {
		case <-doneChan:
			wg.Wait()
			return ErrErrorsLimitExceeded
		case taskChan <- task:
		}
	}

	close(taskChan)
	wg.Wait()

	select {
	case <-doneChan:
		return ErrErrorsLimitExceeded
	default:
	}

	return nil
}

func newWorker(m int, taskChan chan Task, doneChan chan struct{}, errCount *int64, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-doneChan:
			return
		case task, ok := <-taskChan:
			if !ok {
				return
			}

			if err := task(); err != nil {
				if atomic.AddInt64(errCount, 1) == int64(m) {
					close(doneChan)
				}
			}
		}
	}
}
