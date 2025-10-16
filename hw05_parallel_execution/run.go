package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type MutexMaxErrorCounter struct {
	l      sync.Mutex
	errMax int
	errCnt int
}

func (mec *MutexMaxErrorCounter) inc() {
	mec.l.Lock()
	defer mec.l.Unlock()
	mec.errCnt++
}

func (mec *MutexMaxErrorCounter) exeeded() bool {
	mec.l.Lock()
	defer mec.l.Unlock()
	return mec.errCnt >= mec.errMax
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	wg.Add(n)
	tch := make(chan Task, len(tasks))
	errCnt := MutexMaxErrorCounter{errMax: m}

	// spawn goroutins
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for t := range tch {
				if errCnt.exeeded() {
					break
				}
				// do task
				err := t()
				if err != nil {
					errCnt.inc()
				}
			}
		}()
	}

	// write tasks to channel
	for _, t := range tasks {
		tch <- t
	}
	close(tch)

	// wg wait
	wg.Wait()

	if errCnt.exeeded() {
		return ErrErrorsLimitExceeded
	}

	return nil
}
