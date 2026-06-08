package main

import (
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

	done := make(chan void)
	defer close(done)

	// Let's time it
	startTs := time.Now()

	// Run the pattern
	foDataChs := FanOut(done, NForks, dataProcessor, dataCh)
	// Here is the place for your 'fan-in' pattern
	for _, ch := range foDataChs {
		fmt.Println(<-ch)
	}

	// Results and output
	fmt.Println(NForks, "forks used")
	fmt.Println("It takes", time.Since(startTs), "instead of", (NForks * processingTime).Seconds(), "sec")

}

// FanOut is used to fork a stage of pipeline for concurrent executing.
// The stage need to satisfy some criteria:
// * the forked stage is a tight place in programm (too computationally expensive)
// * no matter the order of input processing and output reading
func FanOut[T any](done <-chan void, nForks int, processor func(T) T, inpCh <-chan T) []<-chan T {
	if nForks <= 0 {
		return []<-chan T{}
	}

	outpChs := make([]<-chan T, nForks)

	for i := range nForks {
		ch := make(chan T)
		outpChs[i] = (<-chan T)(ch)

		go func() {
			defer close(ch)

			select {
			case <-done:
				return
			default:
				v, ok := <-inpCh
				if !ok {
					return
				}
				ch <- processor(v)
			}
		}()
	}

	return outpChs
}
