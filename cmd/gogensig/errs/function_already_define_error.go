package errs

import "fmt"

type FuncAlreadyDefinedError struct {
	goSymbolName string
}

func (p *FuncAlreadyDefinedError) Error() string {
	return fmt.Sprintf("function %s already defined", p.goSymbolName)
}

func NewFuncAlreadyDefinedError(goSymbolName string) *FuncAlreadyDefinedError {
	return &FuncAlreadyDefinedError{goSymbolName: goSymbolName}
}
