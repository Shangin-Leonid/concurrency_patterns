package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func run_fan_in() {

	// Preparing
	const NForks = 20

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	processingTime := 2 * time.Second
	dataProcessor := func(v int) int {
		time.Sleep(processingTime)
		return -v
	}

	// Let's time it
	ts := time.Now()

	// Run the pattern
	for v := range FanIn(ctx, FanOut(ctx, NForks, dataProcessor, randIntGenerator(NForks, 50))) {
		fmt.Println(v)
	}

	// Results and output
	fmt.Println(NForks, "forks used")
	fmt.Println("It takes", time.Since(ts), "instead of", (NForks * processingTime).Seconds(), "sec")
}

// FanIn is used to join forked (for example by FanOut) processes (channels).
// The output order is not guaranteed.
//
// You may add 'context' to have an opportunity to stop before all of input channels will be closed.
//
// To use it in 'graceful shutdown' scenario (each input value must be processed):
// * close all input channels
// * don't cancel passed context until output channel will be closed
func FanIn[T any](ctx context.Context, inpChs []<-chan T) <-chan T {
	outpCh := make(chan T)

	wg := &sync.WaitGroup{}
	wg.Add(len(inpChs))

	// Listens and resends safely from input to output channel.
	resender := func(ch <-chan T) {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-ch:
				if !ok {
					return
				}

				select {
				case <-ctx.Done():
					return
				case outpCh <- v:
				}
			}
		}
	}

	// Run resender for each input channel.
	for _, ch := range inpChs {
		go resender(ch)
	}

	// Wait for draining all input channels and close the output channel.
	go func() {
		defer close(outpCh)
		wg.Wait()
	}()

	return outpCh
}
