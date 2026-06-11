package main

import (
	"context"
	"sync"
)

func run_variable_tee_channel() {

}

/*
Корректна ли данная реализация паттерна 'tee channel'? Есть ли здесь потенциальные deadlock, race condition, утечки горутин или другие критические ошибки? Какие есть недостатки у этого кода?
*/

func TeeChannel[T any](ctx context.Context, nChs int, inpCh <-chan T) []<-chan T {

	// Corner cases
	switch {
	case nChs <= 0:
		return []<-chan T{}
	case nChs == 1:
		return []<-chan T{inpCh} // Maybe unsafe
	}

	outpChs, closeOutpChs := MakeSliceOfChs[T](nChs, 0) // Destination channels
	bufChs, _ := MakeSliceOfChs[T](nChs, 0)             // Channels for data transit

	wg := &sync.WaitGroup{} // Synchronization of 'resenders'and buffer filling

	//
	resendFromBufToOutp := func(bufCh <-chan T, outpCh chan T) {
		defer close(outpCh)

		var value T
		var ok bool

		for {
			select {
			case <-ctx.Done():
				wg.Done()
				return
			case value, ok = <-bufCh:
				if !ok {
					wg.Done()
					return
				}
				select {
				case <-ctx.Done():
					wg.Done()
					return
				case outpCh <- value:
					wg.Done()
				}
			}
		}
	}

	// Run resenders from buffers to destination channel
	for i := range nChs {
		go resendFromBufToOutp(bufChs[i], outpChs[i])
	}

	// Fill the buffer channels once
	replicateToBufsOnce := func(value T) (ok bool) {
		for _, bufCh := range bufChs {
			select {
			case <-ctx.Done():
				return false
			default:
			}
			bufCh <- value
		}

		return true
	}

	// Run a goroutine forking (resendig) values to buffer channels
	go func() {
		defer closeOutpChs()

		var value T
		var ok bool

		for {
			wg.Add(nChs)

			select {
			case <-ctx.Done():
				return
			case value, ok = <-inpCh:
				if !ok {
					return
				}
			}

			ok = replicateToBufsOnce(value)
			if !ok {
				return
			}
			wg.Wait()
		}
	}()

	return AsReadOnly(outpChs)
}
