package main

import (
	"context"
	"fmt"
)

// The 'pipeline' pattern: there are
// * source of data
// * some stages (data processors)
// * data flowing from source through stages
// * consistency of these instances and processes
//
// I am using funcs as values (like closures) to name them clearly and not to have conflicts with other files.
// I mean to shadow them in 'run_pipeline()'
func run_pipeline() {

	// Generates and owns the data flow, data emmiter.
	generator := func(ctx context.Context, data ...int) <-chan int {
		dataCh := make(chan int)

		go func() {
			defer close(dataCh)

			for _, v := range data {
				select {
				case <-ctx.Done():
					return
				case dataCh <- v:
				}
			}
		}()

		return dataCh
	}

	// Prepare demo data
	data := []int{1, 2, 3, 4, 5}
	twicer := func(v int) int {
		return v << 1
	}
	signChanger := func(v int) int {
		return -v
	}

	// Run the pipeline
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()
	dataInpCh := generator(ctx, data...)
	pipelineOutpCh := Stage(ctx, signChanger, Stage(ctx, twicer, dataInpCh))

	// "-2 -4 -6 -8 -10" expected
	for v := range pipelineOutpCh {
		fmt.Println(v)
	}

}

// Stage makes and runs a stage of pipeline, using recieved processing func.
// Data flows through it.
func Stage[T any](ctx context.Context, processor func(T) T, inp <-chan T) <-chan T {
	outp := make(chan T)

	go func() {
		defer close(outp)

		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-inp:
				if !ok {
					return
				}
				v = processor(v)
				/*
					// If need. Depends on programm logic.
					// Because 'ctx.Done()' status may be changed since 'processor' started.
					select{
					case <-ctx.Done():
						return
					default:
					}
				*/
				outp <- v
			}
		}
	}()

	return outp
}
