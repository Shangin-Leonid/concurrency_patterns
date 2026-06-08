package main

import (
	"fmt"
	"time"
)

// How to use 'select' with a priority of one channel?
// The result is 2/3 of reading from high priority channel and 1/3 for low one.
// Works (helpful) when writers are faster than readers.
func run_ct_select_priority() {

	generator := func(nWrites int) <-chan int {
		ch := make(chan int)

		go func() {
			defer close(ch)

			for i := range nWrites {
				ch <- i
			}
		}()

		return ch
	}

	const maxReads = 1000

	lowPriorityCh := generator(maxReads)
	highPriorityCh := generator(maxReads)

	lowCounter := 0
	highCounter := 0
	for range maxReads {
		select {
		case <-highPriorityCh:
			highCounter++
			// If no processing exists, than writes will be slower than reads,
			// so reader will process all channels and will always starve.
			time.Sleep(10 * time.Millisecond) // Proccess data
		case <-highPriorityCh:
			highCounter++
			time.Sleep(10 * time.Millisecond) // Proccess data
		case <-lowPriorityCh:
			lowCounter++
			time.Sleep(10 * time.Millisecond) // Proccess data
		}
	}

	fmt.Println("Reads distribution:")
	fmt.Println("High priority channel - ", highCounter)
	fmt.Println("Low priority channel - ", lowCounter)
}
