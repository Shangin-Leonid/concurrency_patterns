package main

import "fmt"

func run_repeat_take() {

	// Infinite generator.
	// Usage by 'takeN(n, repeat(...))' means (n+1) executions of repeater.
	repeat := func(done <-chan void, dataGenerator func() int) <-chan int {
		outpCh := make(chan int)

		go func() {
			defer close(outpCh)

			for {
				select {
				case <-done:
					return
				case outpCh <- dataGenerator():
				}
			}
		}()

		return outpCh
	}

	takeN := func(n int, inpCh <-chan int) <-chan int {
		takenCh := make(chan int)

		go func() {
			defer close(takenCh)

			for range n {
				v := <-inpCh
				takenCh <- v
			}
		}()

		return takenCh
	}

	counter := 0
	dataGenerator := func() int {
		counter++
		return counter
	}

	// Run 'repeat-take' pattern.
	// "1 2 3 4 5" expected
	done := make(chan void)
	defer close(done)
	dataSource := repeat(done, dataGenerator)
	for v := range takeN(5, dataSource) {
		fmt.Println(v)
	}
}
