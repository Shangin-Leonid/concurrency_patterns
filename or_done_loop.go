package main

import (
	"context"
	"fmt"
	"time"
)

func run_ct_or_done_loop() {

	// Instead of using this:
	// (it is wrong because of block during waiting a value in 'channel', while context can be canceled)
	/*
		loop:
			for v := range channel {
				process(v)

				select {
				case <-ctx.Done():
					break loop // or return
				default:
				}
			}
	*/
	// Or instead of using this:
	/*
		loop:
			for {
				select {
				case <-ctx.Done():
					break loop // or return
				case v, ok := <-channel:
					if !ok {
						break loop // or return
					}
					// Do something with val
				}
			}
	*/
	// Just wrap the 'channel' into 'OrDone':

	// Write to 'ch' 5 times
	ch := make(chan int)
	uncontrollableWriter := func() {
		defer close(ch)

		timeout := time.After(5 * time.Second)
		for i := 1; i <= 100; i++ {
			select {
			case <-timeout:
				return
			case ch <- i:
			}
		}
	}

	go uncontrollableWriter()

	// Read from 'ch' until its closing or context canceling
	ctx, cancelCtx := context.WithCancel(context.Background())
	counter := 0
	for v := range OrDoneRange(ctx, ch) {
		fmt.Println(v)

		counter++
		if counter == 4 {
			cancelCtx()
		}
	}
	// "1 2 3 4" expected

}

func OrDoneRange[T any](ctx context.Context, ch <-chan T) <-chan T {
	outpCh := make(chan T, 1)

	go func() {
		defer close(outpCh)

		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-ch:
				if !ok {
					return
				}
				outpCh <- v
			}
		}
	}()

	return outpCh
}
