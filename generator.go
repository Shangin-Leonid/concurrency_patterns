package main

import "fmt"

func run_generator() {

	consumer := func() {
		ch := generator()

		for i := range ch {
			fmt.Println(i)
		}
	}

	consumer()
}

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
