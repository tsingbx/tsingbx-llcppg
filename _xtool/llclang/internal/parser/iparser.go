package parser

import (
	"fmt"
	"os"
	"strconv"
)

type Mode int

func RunParseIntermediateFile(args []string) {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "RunParseIntermediateFile: lack of argument")
		os.Exit(1)
	}
	filename := args[0]
	var mode Mode
	if args[1] != "" {
		modeInt, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "RunParseIntermediateFile: invalid mode")
			os.Exit(1)
		}
		mode = Mode(modeInt)
	}
	err := parseIntermediateFile(filename, mode)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// parses an intermediate file (*.i) and output the corresponding AST to stdout.
func parseIntermediateFile(filename string, mode Mode) error {
	return fmt.Errorf("parseIntermediateFile:not implemented")
}
