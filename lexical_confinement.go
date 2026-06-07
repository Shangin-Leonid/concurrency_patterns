package main

import (
	"fmt"
	"sync"
)

func run_lexical_confinement() {
	// Pass data by value, not by closure capture.
	// Each instance of function processes its own chunk of data.
	// This prevents concurrent accessing to 'allData' without any sync primitives.
	processDataChunk := func(wg *sync.WaitGroup, dataChunk []int) {
		defer wg.Done()

		fmt.Println("Processor instance runs at: ", dataChunk)
	}

	allData := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	wg := sync.WaitGroup{}
	wg.Add(3)
	// Decompose and process 'allData' concurrently.
	go processDataChunk(&wg, allData[:3])
	go processDataChunk(&wg, allData[3:7])
	go processDataChunk(&wg, allData[7:])
	wg.Wait()
}
