package parser

import (
	"fmt"
	"strconv"
)

type Mode int

func RunParseIntermediateFile(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("RunParseIntermediateFile: lack of argument")
	}
	filename := args[0]
	var mode Mode
	if modeStr := args[1]; modeStr != "" {
		modeInt, err := strconv.Atoi(modeStr)
		if err != nil {
			return fmt.Errorf("RunParseIntermediateFile: invalid mode %s: %w", modeStr, err)
		}
		mode = Mode(modeInt)
	}
	err := parseIntermediateFile(filename, mode)
	if err != nil {
		return err
	}
	return nil
}

// parses an intermediate file (*.i) and output the corresponding AST to stdout.
func parseIntermediateFile(filename string, mode Mode) error {
	return fmt.Errorf("parseIntermediateFile:not implemented")
}
