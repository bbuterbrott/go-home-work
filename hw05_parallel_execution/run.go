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
	errorCh := make(chan struct{})
	consumerStopCh := make(chan struct{})
	wg := sync.WaitGroup{}
	startConsumers(n, &wg, queueCh, errorCh, consumerStopCh)

	var errorExit bool
	queueAdderStopCh := make(chan struct{})
	go startErrorCounter(m, errorCh, queueAdderStopCh, consumerStopCh, &errorExit)
	go startQueueAdder(tasks, queueCh, queueAdderStopCh)

	waitForCompletion(&wg, errorCh)
	fmt.Println("Finished processing all tasks")

	if errorExit {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func startConsumers(n int, wg *sync.WaitGroup, queueCh <-chan Task, errorCh chan<- struct{}, stopCh <-chan struct{}) {
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
					select {
					case <-stopCh:
						fmt.Printf("Consumer %v stopped by receiving stop signal\n", i)
						break
					default:
					}

					select {
					case <-stopCh:
						fmt.Printf("Consumer %v stopped by receiving stop signal\n", i)
						break
					case errorCh <- struct{}{}:
						fmt.Printf("Consumer %v finished job with error\n", i)
					}
				} else {
					fmt.Printf("Consumer %v finished job without errors\n", i)
				}
			}
		}(i)
	}
}

func startErrorCounter(m int, errorCh <-chan struct{}, queueAdderStopCh chan<- struct{}, consumerStopCh chan<- struct{}, errorExit *bool) {
	var errCount int
	for {
		_, ok := <-errorCh
		if !ok {
			break
		}
		errCount++
		fmt.Printf("Current error count: %v\n", errCount)
		if m > 0 && errCount >= m {
			fmt.Printf("Got %v errors. Stopping\n", errCount)
			*errorExit = true
			defer close(queueAdderStopCh)
			defer close(consumerStopCh)
			break
		}
	}
}

func startQueueAdder(tasks []Task, queueCh chan<- Task, stopCh <-chan struct{}) {
	for _, task := range tasks {
		select {
		case <-stopCh:
			fmt.Println("Stopped adding new tasks to queue")
			break
		default:
		}

		select {
		case <-stopCh:
			fmt.Println("Stopped adding new tasks to queue")
			break
		case queueCh <- task:
			fmt.Println("Added task")
		}
	}
	defer close(queueCh)
}

func waitForCompletion(wg *sync.WaitGroup, errorCh chan struct{}) {
	wg.Wait()
	go func() {
		for {
			_, ok := <-errorCh
			if !ok {
				break
			}
		}
	}()
	defer close(errorCh)
}
