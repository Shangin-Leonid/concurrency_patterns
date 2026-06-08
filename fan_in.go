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
	for v := range FanIn(FanOut(done, NForks, dataProcessor, randIntGenerator(NForks, 50))) {
		fmt.Println(v)
	}

	// Results and output
	fmt.Println(NForks, "forks used")
	fmt.Println("It takes", time.Since(startTs), "instead of", (NForks * processingTime).Seconds(), "sec")
}

func FanIn[T any](inpChs []<-chan T) <-chan T {
	outpCh := make(chan T)

	go func() {
		defer close(outpCh)

		// Run listener-resender for each of input channels
		wg := &sync.WaitGroup{}
		for _, ch := range inpChs {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for v := range ch {
					outpCh <- v
				}
			}()

		}
		wg.Wait()
	}()

	return outpCh
}
