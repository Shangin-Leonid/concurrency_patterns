package main

import (
	"fmt"
	"time"
)

func run_time_limiter() {

	longF := func() {
		time.Sleep(3 * time.Second)
	}

	quickF := func() {
		time.Sleep(200 * time.Millisecond)
	}

	timeLim := 1 * time.Second

	fmt.Println(WithTimeLimit(timeLim, longF))
	fmt.Println(WithTimeLimit(timeLim, quickF))

}

func WithTimeLimit(timeLimit time.Duration, f func()) (ok bool) {

	ok = true
	sigCh := make(chan void)

	go func() {
		defer close(sigCh)
		f()
	}()

	select {
	case <-time.After(timeLimit):
		return !ok
	case <-sigCh:
		return ok
	}
}
