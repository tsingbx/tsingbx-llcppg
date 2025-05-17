package cltest

import (
	cfg "github.com/goplus/llcppg/cmd/gogensig/config"
)

func NewConvSym(syms ...cfg.SymbolEntry) func(mangleName string) (goName string, err error) {
	return fromSymbTable(cfg.CreateSymbolTable(syms))
}

func GetConvSym(symbFile string) func(mangleName string) (goName string, err error) {
	var symbTable *cfg.SymbolTable
	if symbFile == "" {
		symbTable = cfg.CreateSymbolTable([]cfg.SymbolEntry{})
	} else if tab, err := cfg.NewSymbolTable(symbFile); err == nil {
		symbTable = tab
	} else {
		panic(err)
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
