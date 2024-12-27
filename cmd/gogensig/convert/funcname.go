package convert

import "strings"

type GoFuncName struct {
	goSymbolName string
	recvName     string
	funcName     string
	ptrRecv      bool // if the receiver is a pointer
}

func NewGoFuncName(name string) *GoFuncName {
	l := strings.Split(name, ".")
	if len(l) < 2 {
		return &GoFuncName{goSymbolName: name, funcName: name}
	}
	recvName := l[0]
	ptrRecv := false
	if strings.HasPrefix(recvName, "(*") {
		ptrRecv = true
		recvName = strings.TrimPrefix(recvName, "(*")
		recvName = strings.TrimSuffix(recvName, ")")
	}
	return &GoFuncName{goSymbolName: name, recvName: recvName, funcName: l[1], ptrRecv: ptrRecv}
}

func (p *GoFuncName) HasReceiver() bool {
	return len(p.recvName) > 0
}

func (p *GoFuncName) OriginGoSymbolName() string {
	return p.goSymbolName
}
