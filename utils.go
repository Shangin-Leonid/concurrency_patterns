package main

// To filter channel means to ignore 'false' values and to resend 'true' ones.
// 1) Reads from input
// 2) Check value by predicate
// 3) Resend to output if true
func filterFn[T any](predicate func(T) bool, inpCh <-chan T) <-chan T {
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

// repeat is an infinite data generator (repeater)
func repeat[T any](done <-chan void, dataGenerator func() T) <-chan T {
	outpCh := make(chan T)

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

// takeN is a broker, its only mission is taking and resending a finite sequence of data objects (messages)
func takeN[T any](n int, inpCh <-chan T) <-chan T {
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
