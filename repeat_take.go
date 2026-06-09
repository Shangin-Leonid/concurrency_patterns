package main

import (
	"context"
	"fmt"
)

// 'repeat-take' pattern includes 3 instances:
// * 'repeat' - an infinite data generator (repeater)
// * 'takeN' - a broker, its only mission is taking and resending a finite sequence of data objects (messages)
// * 'consumer' - a final data receiver and processer (I use a printing loop as implicit 'consumer')
//
// Usage by 'takeN(n, repeat(...))' means only (n+1) executions of repeater,
// so the pattern is efficient enough.
//
// I am using funcs as values (like closures) to name them clearly but not to have conflicts with other files.
func run_repeat_take() {

	counter := 0
	dataGenerator := func() int {
		counter++
		return counter
	}

	// Run 'repeat-take' pattern.
	// "1 2 3 4 5" expected
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()
	dataSource := Repeat(ctx, dataGenerator)
	for v := range TakeN(5, dataSource) {
		fmt.Println(v)
	}
}

// Repeat is an infinite data generator (repeater)
func Repeat[T any](ctx context.Context, dataGenerator func() T) <-chan T {
	outpCh := make(chan T)

	go func() {
		defer close(outpCh)

		for {
			select {
			case <-ctx.Done():
				return
			case outpCh <- dataGenerator():
			}
		}
	}()

	return outpCh
}

// TakeN is a broker, its only mission is taking and resending a finite sequence of data objects (messages)
func TakeN[T any](n int, inpCh <-chan T) <-chan T {
	takenCh := make(chan T)

	go func() {
		defer close(takenCh)

		for range n {
			v := <-inpCh
			takenCh <- v
		}
	}()

	return takenCh
}
