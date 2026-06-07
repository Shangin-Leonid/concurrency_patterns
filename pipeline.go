package main

import "fmt"

// I am using funcs as values (like closures) to name them clearly but not to have conflicts with other files.
func run_pipeline() {

	// Generates and owns the data flow, data emmiter.
	generator := func(done <-chan void, data ...int) <-chan int {
		dataFlowCh := make(chan int)

		go func() {
			defer close(dataFlowCh)

			for _, v := range data {
				select {
				case <-done:
					return
				case dataFlowCh <- v:
				}
			}
		}()

		return dataFlowCh
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
	done := make(chan void)

	// Run the pipeline
	// "-2 -4 -6 -8 -10" expected
	for v := range stage(done, signChanger, stage(done, twicer, generator(done, data...))) {
		fmt.Println(v)
	}

}
