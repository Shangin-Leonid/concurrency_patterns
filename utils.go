package main

import "math/rand/v2"

func randIntGenerator(nInts int, upperBound int) <-chan int {
	outpCh := make(chan int)

	go func() {
		defer close(outpCh)

		for range nInts {
			outpCh <- rand.IntN(upperBound)
		}
	}()

	return outpCh
}
