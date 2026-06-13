package main

import (
	"context"
	"fmt"
	"sync"
)

func run_tee_channel() {

	inpCh := make(chan int)
	ctx := context.Background()

	go func() {
		defer close(inpCh)

		inpCh <- 111
		inpCh <- 88
		inpCh <- 0
		inpCh <- 0
		inpCh <- 88
		inpCh <- 111
	}()

	outpChs := TeeChannel(ctx, 4, inpCh)
	readersOrderedOutp := make([][]int, len(outpChs))

	wg := &sync.WaitGroup{}
	reader := func(outp *[]int, inpCh <-chan int) {
		defer wg.Done()

		for v := range inpCh {
			*outp = append(*outp, v)
		}
	}

	wg.Add(len(outpChs))
	for i := range len(outpChs) {
		readersOrderedOutp[i] = []int{}
		go reader(&readersOrderedOutp[i], outpChs[i])
	}
	wg.Wait()

	// resCh := FanIn(ctx, outpChs)

	for i, rOut := range readersOrderedOutp {
		fmt.Println("reader", i, ":", rOut)
	}

}

// TeeChannel is an implementation of the eponymous concurrency pattern.
// The pattern is used for replicating values from one input channel to multiple output channel.
// Value passing is synchronous: a new value can't be written before previous one is not read from all output channels.
//
// ATTENTION. Do you actually need to fork your input channel by this pattern?
// Maybe you can just replicate a value from this channel and push it to several processing functions or goroutines?
//
// Recieves context, number of output channels and input read-only channel.
// Returns a slice of output channel.
//
// TODO:
//   - optimization: maybe add a channel buffer to have non-blocking writing
func TeeChannel[T any](ctx context.Context, nChs int, inpCh <-chan T) []<-chan T {

	// Corner cases for 'nChs'
	switch {
	case nChs <= 0:
		return []<-chan T{}
	case nChs == 1:
		return []<-chan T{inpCh}
	}

	resendFromBufToOutp := func(bufCh <-chan T, outpCh chan T) {
		defer close(outpCh)

		var value T
		var ok bool

		for {
			// Read a value from buffer channel
			select {
			case <-ctx.Done():
				return
			case value, ok = <-bufCh:
				if !ok {
					return
				}
			}

			// Resend the value to calling code (to output channel)
			select {
			case <-ctx.Done():
				return
			case outpCh <- value:
			}
		}
	}

	outpChs, _ := MakeSliceOfChs[T](nChs, 0)          // Destination channels
	bufChs, closeBufChs := MakeSliceOfChs[T](nChs, 0) // Channels for data transit throw support goroutines

	// Run resenders from buffers to destination channel
	for i := range nChs {
		go resendFromBufToOutp(bufChs[i], outpChs[i])
	}

	// Run a goroutine forking (resendig) values to buffer channels
	go func() {
		defer closeBufChs()

		var value T
		var ok bool

		for {
			// Read value from input channel
			select {
			case <-ctx.Done():
				return
			case value, ok = <-inpCh:
				if !ok {
					return
				}
			}

			// Fill the buffer channels with current value
			for _, bufCh := range bufChs {
				select {
				case <-ctx.Done():
					return
				case bufCh <- value:
				}
			}
		}
	}()

	return AsReadOnlyWithUnsafe(outpChs)
}

// DoublingTeeChannel is the same as TeeChannel, but with 2 output channels.
func DoublingTeeChannel[T any](ctx context.Context, in <-chan T) (<-chan T, <-chan T) {
	out1 := make(chan T)
	out2 := make(chan T)
	closeAllOut := func() {
		close(out1)
		close(out2)
	}

	go func() {
		defer closeAllOut()

		var value T
		var ok bool

		for {
			// Read from 'in'
			select {
			case <-ctx.Done():
				return
			case value, ok = <-in:
				if !ok {
					return
				}
			}

			// Replicate to 'out1' and 'out2'
			select {

			case <-ctx.Done():
				return

			case out1 <- value:
				select {
				case <-ctx.Done():
					return
				case out2 <- value:
				}

			case out2 <- value:
				select {
				case <-ctx.Done():
					return
				case out1 <- value:
				}
			}

		}
	}()

	return out1, out2
}
