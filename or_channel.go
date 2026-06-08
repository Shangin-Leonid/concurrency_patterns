package main

import (
	"fmt"
	"time"
)

func run_or_channel() {

	startTS := time.Now()

	<-AnySignal(
		time.After(5*time.Second),
		time.After(4*time.Minute),
		time.After(1*time.Second), // This will finish first
		time.After(3*time.Hour),
		time.After(2*time.Minute),
	)

	// 1 second is expected
	fmt.Println("Done after", time.Since(startTS))
}

// AnySignal
// The pattern 'or channel' close to 'or done', but I find these names wrong,
// because the pattern doesn't actually join input channels.
// It can just signal about closing or reading from one of them.
// I'd better name this as 'any signal' or 'any read' (any channel have been read).
//
// You can also modify the pattern to resend 'signal' value. I haven't done it yet to keep the code simple.
func AnySignal[T any](chs ...<-chan T) <-chan T {

	switch len(chs) {
	case 0:
		return nil
	case 1:
		return chs[0]
	}

	sigCh := make(chan T)
	go func() {
		defer close(sigCh)

		switch len(chs) {
		case 2:
			select {
			case <-chs[0]:
			case <-chs[1]:
			}
		// You can continue with optimization for common cases with len == 3, 4, ... to avoid recursion
		default:
			select {
			case <-chs[0]:
			case <-chs[1]:
			case <-chs[2]:
			case <-AnySignal(append(chs[3:], sigCh)...):
			}
		}

	}()

	return sigCh
}
