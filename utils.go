package main

import "unsafe"

type void struct{}

func AsReadOnly[T any](chs []chan T) []<-chan T {
	outpChs := make([]<-chan T, len(chs))
	for i := range outpChs {
		outpChs[i] = chs[i]
	}

	return outpChs
}

func AsReadOnlyWithUnsafe[T any](chs []chan T) []<-chan T {
	return *(*[]<-chan T)(unsafe.Pointer(&chs))
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
