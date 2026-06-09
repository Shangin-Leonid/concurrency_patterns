package main

import (
	"context"
	"fmt"
	"time"
)

func run_fan_out_fan_in() {

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
	for v := range FanOutFanIn(ctx, NForks, dataProcessor, randIntGenerator(NForks, 50)) {
		fmt.Println(v)
	}

	// Results and output
	fmt.Println(NForks, "forks used")
	fmt.Println("It takes", time.Since(ts), "instead of", (NForks * processingTime).Seconds(), "sec")
}

// FanOutFanIn is a combination of fanOut and fanIn functions.
// I find it vital to encapsulate both patterns to exclude unsafe usage of any one.
//
// Read more in 'FanOut' and 'FanIn' docs.
//
// To use it in 'graceful shutdown' scenario (each input value must be processed):
// * close the input channel
// * don't cancel passed context until output channel will be closed
func FanOutFanIn[T any](ctx context.Context, nForks int, processor func(T) T, inpCh <-chan T) <-chan T {
	return FanIn(ctx, FanOut(ctx, nForks, processor, inpCh))
}
