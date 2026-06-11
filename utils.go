package main

import "math/rand/v2"

type void struct{}

func randIntGenerator(nInts int, upperBound int) <-chan int {
	outpCh := make(chan int)

	go func() {
		defer close(outpCh)

		for range nInts {
			outpCh <- rand.IntN(upperBound)
		}
	}()

	return outpCh
}

func AsReadOnly[T any](chs []chan T) []<-chan T {
	outpChs := make([]<-chan T, len(chs))
	for i := range outpChs {
		outpChs[i] = chs[i]
	}

	return outpChs
}

func MakeSliceOfChs[T any](nChs int, bufSz int) (_ []chan T, closeFunc func()) {

	if nChs <= 0 {
		return []chan T{}, nil
	}

	chs := make([]chan T, nChs)
	for i := range chs {
		chs[i] = make(chan T, bufSz)
	}

	closeFunc = func() {
		for i := range chs {
			close(chs[i])
		}
	}

	return chs, closeFunc

}
