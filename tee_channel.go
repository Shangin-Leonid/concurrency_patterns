package main

import (
	"context"
	"sync"
)

func TeeChannel[T any](ctx context.Context, nChs int, inpCh <-chan T) []<-chan T {

	// Corner cases
	switch {
	case nChs <= 0:
		return []<-chan T{}
	case nChs == 1:
		return []<-chan T{inpCh} // Maybe unsafe
	}

	outpChs, closeOutpChs := MakeSliceOfChs(nChs, 0)
	buffChs, _ := MakeSliceOfChs(nChs, 0)

	wg := &sync.WaitGroup{}

	// Fill the buffer channels
	fillBuffers := func(value T) (needExit bool) {
		for _, bCh := range buffChs {
			select {
			case <-ctx.Done():
				return true
			default:
			}
			wg.Add(1)
			bCh <- value
		}

		return false
	}

	resender := func(bufCh <-chan T, outpCh chan T) {
		defer close(outpCh)

		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-bufCh:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case outpCh <- v:
					wg.Done()
				}
			}
		}
	}

	go func() {
		defer closeOutpChs()

		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-inpCh:
				if !ok {
					return
				}

				needExit := fillBuffers(v)
				if needExit {
					return
				}
				wg.Wait()
			}
		}
	}()

	for i := range nChs {
		go resender(buffChs[i], outpChs[i])
	}

	return AsReadOnly(outpChs)
}
