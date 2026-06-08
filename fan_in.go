package main

import (
	"fmt"
	"sync"
	"time"
)

func run_fan_in() {

	// Preparing
	const NForks = 20

	done := make(chan void)
	defer close(done)

	processingTime := 2 * time.Second
	dataProcessor := func(v int) int {
		time.Sleep(processingTime)
		return -v
	}

	// Let's time it
	startTs := time.Now()

	// Run the pattern
	for v := range FanIn(FanOut(done, NForks, dataProcessor, randIntGenerator(NForks, 50))...) {
		fmt.Println(v)
	}

	// Results and output
	fmt.Println(NForks, "forks used")
	fmt.Println("It takes", time.Since(startTs), "instead of", (NForks * processingTime).Seconds(), "sec")
}

// FanIn is used to join forked (for example in FanOut) processes (channels).
// The output order is not guaranteed
//
// You may add 'done' channel to have an opportunity to stop before all of input channels will be closed.
func FanIn[T any](inpChs ...<-chan T) <-chan T {
	outpCh := make(chan T)

	wg := &sync.WaitGroup{}
	wg.Add(len(inpChs))

	// Run listener-resender for each of input channels
	for _, ch := range inpChs {
		go func() {
			defer wg.Done()
			for v := range ch {
				outpCh <- v
			}
		}()
	}

	// Wait for draining all input channels and close the output channel.
	go func() {
		defer close(outpCh)
		wg.Wait()
	}()

	return outpCh
}
