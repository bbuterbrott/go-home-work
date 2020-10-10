package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when it has received M errors from tasks. If M <= 0 that means that errors will be ignored.
func Run(tasks []Task, n int, m int) error {
	fmt.Printf("Task size: %v, n=%v, m=%v\n", len(tasks), n, m)

	if n <= 0 {
		return fmt.Errorf("n should be a positive number")
	}
	ignoreErrors := m <= 0

	queueCh := make(chan Task)
	var errCountValue int32
	errCount := &errCountValue
	go func() {
		for _, task := range tasks {
			fmt.Printf("Current error count: %v\n", *errCount)
			if int(*errCount) >= m && !ignoreErrors {
				fmt.Printf("Got %v errors. Sending done signal\n", *errCount)
				break
			}
			queueCh <- task
		}
		close(queueCh)
	}()

	fmt.Println("Started processing tasks")

	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			fmt.Printf("Consumer %v started\n", i)
			defer wg.Done()

			for {
				task, ok := <-queueCh

				if !ok {
					fmt.Printf("Consumer %v finished working\n", i)
					break
				}

				err := task()

				if err != nil {
					atomic.AddInt32(errCount, 1)
					fmt.Printf("Consumer %v finished job with error. Total error count: %v\n", i, *errCount)
				} else {
					fmt.Printf("Consumer %v finished job without errors\n", i)
				}
			}
		}(i)
	}

	wg.Wait()

	fmt.Println("Finished processing all tasks")

	if int(*errCount) >= m && !ignoreErrors {
		return ErrErrorsLimitExceeded
	}

	return nil
}
