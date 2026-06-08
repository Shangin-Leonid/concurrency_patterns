package main

import "fmt"

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
	generator := func(done <-chan void, data ...int) <-chan int {
		dataCh := make(chan int)

		go func() {
			defer close(dataCh)

			for _, v := range data {
				select {
				case <-done:
					return
				case dataCh <- v:
				}
			}
		}()

		return dataCh
	}

	// Makes and runs a stage of pipeline, using recieved processing func.
	// Data flows through it.
	stage := func(done <-chan void, processor func(int) int, inp <-chan int) <-chan int {
		outp := make(chan int)

		go func() {
			defer close(outp)

			for {
				select {
				case <-done:
					return
				case v, ok := <-inp:
					if !ok {
						return
					}
					v = processor(v)
					/*
						// If need. Depends on programm logic.
						// Because 'done' status may be changed since 'processor' started.
						select{
						case <-done:
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

	// Prepare demo data
	data := []int{1, 2, 3, 4, 5}
	twicer := func(v int) int {
		return v << 1
	}
	signChanger := func(v int) int {
		return -v
	}

	// Run the pipeline
	done := make(chan void)
	defer close(done)
	dataInpCh := generator(done, data...)
	pipelineOutpCh := stage(done, signChanger, stage(done, twicer, dataInpCh))

	// "-2 -4 -6 -8 -10" expected
	for v := range pipelineOutpCh {
		fmt.Println(v)
	}

}
