package main

import (
	"fmt"
	"time"
)

func run_ct_or_done_loop() {

	// Instead of using this:
	// (it is wrong because of block during waiting a value in 'channel', while 'done' can give a signal)
	/*
		loop:
			for v := range channel {
				process(v)

				select {
				case <-done:
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
				case <-done:
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

	// Read from 'ch' until its closing or 'done' signal
	done := make(chan void)
	counter := 0
	for v := range OrDoneRange(done, ch) {
		fmt.Println(v)

		counter++
		if counter == 4 {
			close(done)
		}
	}
	// "1 2 3 4" expected

}

func OrDoneRange[T any](done <-chan void, ch <-chan T) <-chan T {
	outpCh := make(chan T, 1)

	go func() {
		defer close(outpCh)

		for {
			select {
			case <-done:
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
