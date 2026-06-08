package main

import (
	"fmt"
)

func run_channel_filter() {

	isEven := func(n int) bool {
		return n%2 == 0
	}

	// To filter channel means to ignore 'false' values and to resend 'true' ones.
	for v := range FilterFn(isEven, randIntGenerator(10, 50)) {
		fmt.Println(v)
	}

}

// To filter channel means to ignore 'false' values and to resend 'true' ones.
// 1) Reads from input
// 2) Check value by predicate
// 3) Resend to output if true
func FilterFn[T any](predicate func(T) bool, inpCh <-chan T) <-chan T {
	outpCh := make(chan T)

	go func() {
		defer close(outpCh)

		for v := range inpCh {
			if predicate(v) {
				outpCh <- v
			}
		}
	}()

	return outpCh
}
