package main

import (
	"errors"
	"fmt"
	"time"
)

func run_worker_pool() {
	requests := []int{1, 2, 3, 4, 5, 6}
	reqHandler := func(req int) error {
		time.Sleep(2 * time.Second)
		fmt.Print(req, " ")
		return nil
	}

	wp, err := NewWorkerPool(10, reqHandler)
	if err != nil {
		fmt.Println(err)
		return
	}

	for r := range requests {
		wp.Handle(r)
	}
	wp.Wait()
	fmt.Println()
}

var (
	ErrInvalidArgsInWPInit = errors.New("invalid args in 'WorkerPool' initialization")
)

// TODO add errors handling
// TODO add context
// TODO add docs
type WorkerPool[Req any] struct {
	handler func(Req) error
	pool    chan void
}

func NewWorkerPool[Req any](nMaxWorkers int, requestHandler func(Req) error) (*WorkerPool[Req], error) {
	if nMaxWorkers < 1 {
		return &WorkerPool[Req]{}, ErrInvalidArgsInWPInit
	}

	wp := WorkerPool[Req]{
		handler: requestHandler,
		pool:    make(chan void, nMaxWorkers),
	}

	for range cap(wp.pool) {
		wp.freeWorker()
	}

	return &wp, nil
}

func (wp *WorkerPool[Req]) freeWorker() {
	wp.pool <- void{}

}

func (wp *WorkerPool[Req]) takeWorker() {
	<-wp.pool

}

// Handle
func (wp *WorkerPool[Req]) Handle(request Req) {
	wp.takeWorker()

	go func() {
		defer wp.freeWorker()
		_ = wp.handler(request)
	}()
}

// Wait all workers finishing
func (wp *WorkerPool[Req]) Wait() {
	for range cap(wp.pool) {
		wp.takeWorker()
	}
}
