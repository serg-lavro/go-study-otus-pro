package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	channelRelay := func(stOut In) Out {
		out := make(Bi)

		go func() {
			defer func() { for range stOut {} }() //nolint
			defer close(out)
			for {
				select {
				case v, ok := <-stOut:
					if !ok {
						return
					}
					out <- v
				case <-done:
					return
				}
			}
		}()

		return out
	}

	out := channelRelay(in)
	for _, st := range stages {
		out = st(out)
		out = channelRelay(out)
	}

	return out
}
