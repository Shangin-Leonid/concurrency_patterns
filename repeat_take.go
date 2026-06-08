package main

import "fmt"

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
