package main

import "unsafe"

type void struct{}

// AsReadOnly converts a slice of bidirectional channels to read only ones.
func AsReadOnly[T any](chs []chan T) []<-chan T {
	outpChs := make([]<-chan T, len(chs))
	for i := range outpChs {
		outpChs[i] = chs[i]
	}

	return outpChs
}

// AsReadOnlyWithUnsafe converts a slice of bidirectional channels to read only ones.
func AsReadOnlyWithUnsafe[T any](chs []chan T) []<-chan T {
	return *(*[]<-chan T)(unsafe.Pointer(&chs))
}

// MakeSliceOfChs makes a slice of bidirectional channels and initializes each of them.
//
// Params:
//   - 'nChs' - amount of channels
//   - 'bufSz' - buffer size (use '0' if no buffer need)
//
// Returns:
//   - slice of initialized bidirectional channels
//   - function of closing all returned channels
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
