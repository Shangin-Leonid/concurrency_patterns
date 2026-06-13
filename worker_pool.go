package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func run_worker_pool() {

	const nVals = 20
	in := constIntGenerator(nVals, 7)
	dur := 3 * time.Second
	counter := 0
	procFunc := func(_ int) int {
		time.Sleep(dur)
		counter++
		return counter
	}

	const nWorkers = 10
	ts := time.Now()
	for v := range WorkerPool(context.Background(), nVals/2, procFunc, in) {
		fmt.Print(v, " ")
	}
	fmt.Println()
	fmt.Println("Used ", nWorkers, "workers, waste ", time.Since(ts), " instead of ", (nVals * dur).Seconds(), "s")

}

// WorkerPool
func WorkerPool[I, O any](ctx context.Context, nWorkers int, workerFunc func(I) O, in <-chan I) <-chan O {
	out := make(chan O, nWorkers)

	go func() {
		defer close(out)

		wg := &sync.WaitGroup{}
		wg.Add(nWorkers)

		worker := func() {
			defer wg.Done()

			var inValue I
			var ok bool
			var outValue O

			for {
				select {
				case <-ctx.Done():
					return
				case inValue, ok = <-in:
					if !ok {
						return
					}
				}

				outValue = workerFunc(inValue)

				select {
				case <-ctx.Done():
					return
				case out <- outValue:
				}
			}
		}

		for range nWorkers {
			go worker()
		}

		wg.Wait()
	}()

	return out
}
