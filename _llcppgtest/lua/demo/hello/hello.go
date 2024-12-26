package main

import (
	"_llcppgtest/lua/out/lua"

	"github.com/goplus/llgo/c"
)

func main() {
	L := lua.Newstate__1()
	defer L.Close()
	L.Openlibs()
	// 0 -> lua.OK
	if res := L.Loadstring(c.Str("print('hello world')")); res != 0 {
		panic("error")
	}
	// MULTRET -> -1
	if res := L.Pcallk(c.Int(0), c.Int(-1), c.Int(0), 0, nil); res != 0 {
		panic("error")
	}
}
