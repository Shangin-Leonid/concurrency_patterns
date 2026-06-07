package main

import (
	"errors"
	"fmt"
	"os"
)

type void struct{}

// Sentinel errors
var (
	ErrInvalidCommandLineArg = errors.New("invalid command line argument. Pass the only argument - name of pattern (see the list in Readme.md).")
)

// Get name of function to run using command line arg.
var funcByPatternName = map[string]func(){
	"lexical_confinement": run_lexical_confinement,
	"generator":           run_generator,
	"or_channel":          run_or_channel,
	"any_signal":          run_or_channel,
	"pipeline":            run_pipeline,

	"ct_select_priority": run_ct_select_priority,
}

func main() {

	defer func() {
		if p := recover(); p != nil {
			fmt.Println("Panic: ", p)
			return
		}
	}()

	err := parseCommandLineArgs(os.Args)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Print("Run the '", os.Args[1], "' pattern...\n\n")
	funcByPatternName[os.Args[1]]()
	fmt.Println()
}

func parseCommandLineArgs(args []string) error {
	if len(args) != 2 {
		return ErrInvalidCommandLineArg
	}

	if _, ok := funcByPatternName[args[1]]; !ok {
		return ErrInvalidCommandLineArg
	}

	return nil
}
