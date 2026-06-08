package main

import "fmt"

// I am using funcs as values (like closures) to name them clearly but not to have conflicts with other files.
func run_generator() {

	// Generator is a channel owner. It should:
	// * instantiate the channel
	// * close the channel
	// * perform writes, or pass ownership to another goroutine
	// * return to consumer a channel to read from
	generator := func() <-chan int {

		ch := make(chan int)

		producer := func() {
			defer close(ch)

			for i := range 5 {
				ch <- i
			}
		}

		go producer()

		return ch
	}

	// Consumer has no response for recieved channel. Only reads and processes as he wants.
	consumer := func() {
		ch := generator()

		for i := range ch {
			fmt.Println(i)
		}
	}

	// Run the pattern.
	// "0 1 2 3 4" expected
	consumer()
}
