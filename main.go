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

const (
	RunAllPatternsUserCommand = "all"
)

// Get name of function to run using command line arg.
var funcByPatternName = map[string]func(){
	"lexical_confinement": run_lexical_confinement,
	"generator":           run_generator,
	"repeat_take":         run_repeat_take,
	"or_channel":          run_or_channel,
	"any_signal":          run_or_channel,
	"pipeline":            run_pipeline,
	"channel_filter":      run_channel_filter,
	"fan_out":             run_fan_out,
	"fan_in":              run_fan_in,
	"fan_out_fan_in":      run_fan_out_fan_in,

	"ct_select_priority": run_ct_select_priority,
	"ct_or_done_loop":    run_ct_or_done_loop,
}

func main() {

	defer func() {
		if p := recover(); p != nil {
			fmt.Println("> Panic: ", p)
			return
		}
	}()

	err := isUserCommandValid(os.Args)
	if err != nil {
		fmt.Println("> Error: ", err)
		printValidCommands()
		return
	}

	if os.Args[1] == RunAllPatternsUserCommand {
		for pn := range funcByPatternName {
			runPattern(pn)
		}
	} else {
		runPattern(os.Args[1])
	}

}

func isUserCommandValid(args []string) error {
	if len(args) != 2 {
		return ErrInvalidCommandLineArg
	}

	if args[1] == "all" {
		return nil
	}

	if _, ok := funcByPatternName[args[1]]; !ok {
		return ErrInvalidCommandLineArg
	}

	return nil
}

func printValidCommands() {
	fmt.Println()
	fmt.Println("> List of valid commands (command line arguments):")
	fmt.Println("  * all")
	for k := range funcByPatternName {
		fmt.Println("  *", k)
	}
	fmt.Println()
}

func runPattern(patternName string) {
	fmt.Println()
	fmt.Print("> '", patternName, "' pattern is starting...\n")

	funcByPatternName[patternName]()

	fmt.Println("> Finished")
	fmt.Println()
}
