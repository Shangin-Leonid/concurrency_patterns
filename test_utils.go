package main

import (
	"context"
	"math/rand/v2"
)

func randIntGenerator(nInts int, upperBound int) <-chan int {
	return TakeN(nInts, Repeat(context.Background(), func() int { return rand.IntN(upperBound) }))
}

func constIntGenerator(nInts int, val int) <-chan int {
	return TakeN(nInts, Repeat(context.Background(), func() int { return val }))
}
