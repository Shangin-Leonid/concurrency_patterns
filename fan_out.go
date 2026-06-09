package main

import (
	"context"
	"fmt"
	"time"
)

func run_fan_out() {

	// Preparing
	const NForks = 20
	dataCh := randIntGenerator(NForks, 100)

	processingTime := 2 * time.Second
	dataProcessor := func(v int) int {
		time.Sleep(processingTime)
		return -v
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	// Let's time it
	ts := time.Now()

	// Run the pattern
	foDataChs := FanOut(ctx, NForks, dataProcessor, dataCh)
	// Here is the place for your 'fan-in' pattern
	for _, ch := range foDataChs {
		fmt.Println(<-ch)
	}

	// Results and output
	fmt.Println(NForks, "forks used")
	fmt.Println("It takes", time.Since(ts), "instead of", (NForks * processingTime).Seconds(), "sec")

}

// FanOut is used to fork a stage of pipeline for concurrent executing.
//
// In fact, the stage need to satisfy some criteria:
// * the forked stage is a tight place in programm (too computationally expensive)
// * no matter the order of input processing and output reading
//
// To use it in 'graceful shutdown' scenario (each input value must be processed):
// * close the input channel
// * don't cancel passed context until output channel will be closed
func FanOut[T any](ctx context.Context, nForks int, processor func(T) T, inpCh <-chan T) []<-chan T {
	if nForks <= 0 {
		return []<-chan T{}
	}

	worker := func(outpCh chan T) {
		defer close(outpCh)

		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-inpCh:
				if !ok {
					return
				}

				v = processor(v)

				select {
				case <-ctx.Done():
					return
				case outpCh <- v:
				}
			}
		}
	}

	outpChs := make([]<-chan T, nForks)
	for i := range outpChs {
		ch := make(chan T)
		outpChs[i] = (<-chan T)(ch)
		go worker(ch)
	}

	return outpChs
}
