package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func wrapper(st Stage, in In, stop In) Out {
	stOut := st(in)
	out := make(Bi)

	go func() {
		for {
			select {
			case v, ok := <-stOut:
				if !ok {
					close(out)
					return
				}
				out <- v
			case <-stop:
				close(out)
				for range stOut {} //nolint
				return
			}
		}
	}()

	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	firstIn := make(Bi)
	out := In(firstIn)
	for _, st := range stages {
		out = wrapper(st, out, done)
	}

	go func() {
		defer close(firstIn)
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				firstIn <- v
			case <-done:
				return
			}
		}
	}()

	return out
}
