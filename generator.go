package main

import "fmt"

func run_generator() {

	// Consumer has now response for recieved channel. Only reads and processes as he wants.
	consumer := func() {
		ch := generator()

		for i := range ch {
			fmt.Println(i)
		}
	}

	consumer()
}

// Generator is a channel owner. It should:
// * instantiate the channel
// * close the channel
// * perform writes, or pass ownership to another goroutine
// * return to consumer a channel to read from
func generator() <-chan int {

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
