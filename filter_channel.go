package main

import (
	"fmt"
	"math/rand/v2"
)

func run_channel_filter() {

	randGenerator := func(n int) <-chan int {
		outpCh := make(chan int)

		go func() {
			defer close(outpCh)

			for range n {
				outpCh <- rand.IntN(100)
			}
		}()

		return outpCh
	}

	isEven := func(n int) bool {
		return n%2 == 0
	}

	// To filter channel means to ignore 'false' values and to resend 'true' ones.
	for v := range filterFn(isEven, randGenerator(10)) {
		fmt.Print(v, " ")
	}

}
