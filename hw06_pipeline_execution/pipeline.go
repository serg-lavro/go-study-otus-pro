package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func wrapper(st Stage, in In, stop chan struct{}) Out {
	stOut := st(in)
	out := make(Bi)
	var stopped bool

	go func() {
		for {
			select {
			case v, ok := <-stOut:
				if !ok {
					if !stopped {
						close(out)
					}
					return
				}
				if !stopped {
					out <- v
				}
			case <-stop:
				if !stopped {
					close(out)
					stopped = true
				}
			}
		}
	}()

	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	stops := make([]chan struct{}, len(stages))
	for i, st := range stages {
		stops[i] = make(chan struct{})
		out = wrapper(st, out, stops[i])
	}

	buf := []interface{}{}
	finishPipeline := make(chan struct{})
	var canceled bool
	go func() {
		defer close(finishPipeline)
		for {
			select {
			case <-done:
				canceled = true
				for _, s := range stops {
					close(s)
				}
				return
			case v, ok := <-out:
				if !ok {
					return
				}
				buf = append(buf, v)
			}
		}
	}()

	res := make(Bi)
	<-finishPipeline
	if canceled {
		close(res)
		return res
	}

	go func() {
		defer close(res)
		for _, v := range buf {
			res <- v
		}
	}()

	return res
}
