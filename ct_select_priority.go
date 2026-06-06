package main

import (
	"fmt"
	"time"
)

func run_ct_select_priority() {
	const maxReads = 1000
	lowPriorityCh := makeChWriter(maxReads)
	highPriorityCh := makeChWriter(maxReads)

	lowCounter := 0
	highCounter := 0
	for range maxReads {
		select {
		case <-highPriorityCh:
			highCounter++
			// Proccess data read from channel
			time.Sleep(10 * time.Millisecond)
		case <-highPriorityCh:
			highCounter++
			// Proccess data read from channel
			time.Sleep(10 * time.Millisecond)
		case <-lowPriorityCh:
			lowCounter++
			// Proccess data read from channel
			time.Sleep(10 * time.Millisecond)
		}
	}

	fmt.Println("Reads distribution:")
	fmt.Println("High priority channel - ", highCounter)
	fmt.Println("Low priority channel - ", lowCounter)
}

func makeChWriter(nWrites int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)

		for i := range nWrites {
			ch <- i
		}
	}()

	return ch
}
