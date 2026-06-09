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
	startTs := time.Now()

	// Run the pattern
	for v := range FanIn(ctx, FanOut(ctx, NForks, dataProcessor, randIntGenerator(NForks, 50))...) {
		fmt.Println(v)
	}

	// Results and output
	fmt.Println(NForks, "forks used")
	fmt.Println("It takes", time.Since(startTs), "instead of", (NForks * processingTime).Seconds(), "sec")
}

// FanIn is used to join forked (for example in FanOut) processes (channels).
// The output order is not guaranteed
//
// You may add 'context' to have an opportunity to stop before all of input channels will be closed.
func FanIn[T any](ctx context.Context, inpChs ...<-chan T) <-chan T {
	outpCh := make(chan T)

	wg := &sync.WaitGroup{}
	wg.Add(len(inpChs))

	// Run listener-resender for each of input channels
	for _, ch := range inpChs {
		go func(ch <-chan T) {
			defer wg.Done()
			for v := range OrDoneRange(ctx, ch) {
				select {
				case <-ctx.Done():
					return
				case outpCh <- v:
				}
			}
		}(ch)
	}

	// Wait for draining all input channels and close the output channel.
	go func() {
		defer close(outpCh)
		wg.Wait()
	}()

	return outpCh
}
