package hw06pipelineexecution

import "sync/atomic"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var canceled atomic.Bool
	interimIn := make(Bi)

	// goroutine for managing input: receiving and closing
	go func() {
		defer close(interimIn)
		for {
			i, ok := <-in
			if !ok || canceled.Load() {
				return
			}
			interimIn <- i
		}
	}()

	out := Out(interimIn)
	for _, st := range stages {
		out = st(out)
	}

	// goroutine for collecting output and managing shutdown
	buf := []interface{}{}
	finishPipeline := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				canceled.Store(true)
				// flushing output so the close of input would propagate through
				// upstream stages; closing finish channel prior flush for perf
				close(finishPipeline)
				for {
					_, ok := <-out
					if !ok {
						break
					}
				}
				return
			case v, ok := <-out:
				if !ok {
					close(finishPipeline)
					return
				}
				buf = append(buf, v)
			}
		}
	}()

	res := make(Bi)

	<-finishPipeline

	if canceled.Load() {
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
