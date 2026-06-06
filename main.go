package main

import (
	"errors"
	"fmt"
	"os"
)

// Sentinel errors
var (
	ErrInvalidCommandLineArg = errors.New("invalid command line argument. Pass the only argument - name of pattern (see the list in Readme.md).")
)

// Get name of function to run using command line arg.
var funcByPatternName = map[string]func(){
	"generator": run_generator,

	"ct_select_priority":     run_ct_select_priority,
	"ct_lexical_confinement": run_ct_lexical_confinement,
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

	fmt.Println()
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
