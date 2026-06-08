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
