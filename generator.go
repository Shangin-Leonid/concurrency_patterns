package main

import "fmt"

func run_generator() {

	ch := make(chan int)
	go func() {
		defer close(ch)

		for i := range 5 {
			ch <- i
		}
	}()

	for i := range ch {
		fmt.Println(i)
	}

}
