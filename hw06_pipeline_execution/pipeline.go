package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	stageIn := in
	for _, stage := range stages {
		stageIn = stage(withDone(stageIn, done))
	}
	return stageIn
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
