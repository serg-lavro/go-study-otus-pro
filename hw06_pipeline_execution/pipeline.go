package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func channelRelay(stOut In, stop In) Out {
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
	out := channelRelay(in, done)
	for _, st := range stages {
		out = st(out)
		out = channelRelay(out, done)
	}

	return out
}
