package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when it has received M errors from tasks. If M <= 0 that means that errors will be ignored.
func Run(tasks []Task, n int, m int) error {
	fmt.Printf("Task size: %v, n=%v, m=%v\n", len(tasks), n, m)

	if n <= 0 {
		return fmt.Errorf("number of tasks should be a positive number, n=%v", n)
	}

	queueCh := make(chan Task)
	resultCh := make(chan error, n)
	wg := &sync.WaitGroup{}

	wg.Add(n)
	for i := 0; i < n; i++ {
		go startConsumer(i, wg, queueCh, resultCh)
	}

	errorExit := startQueueAdder(m, tasks, queueCh, resultCh)

	go func() {
		for range resultCh {
		}
	}()

	waitForCompletion(wg, resultCh)

	if errorExit {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func startConsumer(i int, wg *sync.WaitGroup, queueCh <-chan Task, resultCh chan<- error) {
	fmt.Printf("Consumer %v started\n", i)
	defer wg.Done()

	for task := range queueCh {
		resultCh <- task()
	}
}

func startQueueAdder(m int, tasks []Task, queueCh chan<- Task, resultCh <-chan error) bool {
	var errCount int
	defer close(queueCh)

	for _, task := range tasks {
		for {
			if len(resultCh) == 0 {
				break
			}

			err := <-resultCh

			if m > 0 && err != nil {
				errCount++
				fmt.Printf("Current error count: %v\n", errCount)
				if errCount == m {
					fmt.Printf("Got %v errors. Stopping\n", errCount)
					return true
				}
			}
		}

		queueCh <- task
		fmt.Println("Added new task to queue")
	}

	return false
}

func waitForCompletion(wg *sync.WaitGroup, resultCh chan error) {
	go func() {
		for range resultCh {
		}
	}()

	wg.Wait()
	defer close(resultCh)
}
