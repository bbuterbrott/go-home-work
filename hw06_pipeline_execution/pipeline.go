package hw06_pipeline_execution //nolint:golint,stylecheck

import (
	"log"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil {
		log.Print("in channel must not be nil")
		return nil
	}

	if len(stages) == 0 {
		log.Print("you must specify at least one stage")
		return nil
	}

	stageCh := in
	for _, stage := range stages {
		stageCh = stage(withDone(stageCh, done))
	}
	return stageCh
}

func withDone(in In, done In) In {
	newIn := make(chan interface{})
	go func() {
		defer close(newIn)
	f:
		for {
			select {
			case data, ok := <-in:
				if !ok {
					break f
				}
				newIn <- data
			case <-done:
				break f
			}
		}
	}()
	return newIn
}
