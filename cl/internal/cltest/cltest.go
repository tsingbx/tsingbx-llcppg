package cltest

import (
	cfg "github.com/goplus/llcppg/cmd/gogensig/config"
)

func NewConvSym(syms ...cfg.SymbolEntry) func(mangleName string) (goName string, err error) {
	return fromSymbTable(cfg.CreateSymbolTable(syms))
}

func GetConvSym(symbFile string) func(mangleName string) (goName string, err error) {
	if symbFile == "" {
		panic("symbol file not set")
	}
	symbTable, err := cfg.NewSymbolTable(symbFile)
	if err != nil {
		// NOTE(xsw): not a good idea, but make sense in test cases
		return NewConvSym()
	}
	return fromSymbTable(symbTable)
}

func fromSymbTable(symbTable *cfg.SymbolTable) func(mangleName string) (goName string, err error) {
	return func(mangleName string) (goName string, err error) {
		item, err := symbTable.LookupSymbol(mangleName)
		if err != nil {
			return
		}
		return item.GoName, nil
	}
}
