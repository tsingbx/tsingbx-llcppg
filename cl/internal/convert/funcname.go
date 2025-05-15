package convert

import (
	"strings"
)

// GoFuncSpec parses and stores components of a Go function name
// Input examples from llcppg.symb.json's "go" field:
// 1. Simple function: "AddPatchToArray"
// 2. Method with pointer receiver: "(*Sqlite3Stmt).Sqlite3BindParameterIndex"
// 3. Method with value receiver: "CJSONBool.CJSONCreateBool"
type GoFuncSpec struct {
	GoSymbName string // original full name from input
	FnName     string // function name without receiver
	IsMethod   bool   // if the function is a method
	RecvName   string // receiver name
	PtrRecv    bool   // if the receiver is a pointer
}

// - "AddPatchToArray" -> {goSymbolName: "AddPatchToArray", funcName: "AddPatchToArray"}
// - "(*Sqlite3Stmt).Sqlite3BindParameterIndex" -> {goSymbolName: "...", recvName: "Sqlite3Stmt", funcName: "Sqlite3BindParameterIndex", ptrRecv: true}
// - "CJSONBool.CJSONCreateBool" -> {goSymbolName: "...", recvName: "CJSONBool", funcName: "CJSONCreateBool", ptrRecv: false}
func NewGoFuncSpec(name string) *GoFuncSpec {
	l := strings.Split(name, ".")
	if len(l) < 2 {
		return &GoFuncSpec{GoSymbName: name, FnName: name, IsMethod: false}
	}
	recvName := l[0]
	ptrRecv := false
	if strings.HasPrefix(recvName, "(*") {
		ptrRecv = true
		recvName = strings.TrimPrefix(recvName, "(*")
		recvName = strings.TrimSuffix(recvName, ")")
	}
	return &GoFuncSpec{GoSymbName: name, RecvName: recvName, FnName: l[1], PtrRecv: ptrRecv, IsMethod: true}
}

func (g *GoFuncSpec) IsIgnore() bool {
	return g.GoSymbName == "-"
}
