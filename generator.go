package main

import "fmt"

func run_generator() {

	ch := generateChWriter()

	for i := range ch {
		fmt.Println(i)
	}

}

func generateChWriter() <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)

		for i := range 5 {
			ch <- i
		}
	}()

	return ch
}
