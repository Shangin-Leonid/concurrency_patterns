package main

import "math/rand/v2"

func randIntGenerator(nInts int, upperBound int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for range nInts {
			out <- rand.IntN(upperBound)
		}
	}()

	return out
}

func constIntGenerator(nInts int, val int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for range nInts {
			out <- val
		}
	}()

	return out
}
