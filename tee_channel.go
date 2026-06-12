package main

import (
	"context"
	"fmt"
)

func run_tee_channel() {

	inpCh := make(chan int)
	ctx := context.Background()

	go func() {
		defer close(inpCh)

		inpCh <- 10
		inpCh <- 20
		inpCh <- 30
	}()

	outpChs := TeeChannel(ctx, 3, inpCh)

	resCh := FanIn(ctx, outpChs)

	for v := range resCh {
		fmt.Print(v, " ")
	}
	fmt.Println()

}

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
			// TODO: maybe add a channel buffer to have non-blocking writing
			for _, bufCh := range bufChs {
				select {
				case <-ctx.Done():
					return
				case bufCh <- value:
				}
			}
		}
	}()

	return AsReadOnly(outpChs)
}
