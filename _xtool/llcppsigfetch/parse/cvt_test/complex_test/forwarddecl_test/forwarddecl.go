package main

import (
	test "github.com/goplus/llcppg/_xtool/llcppsigfetch/parse/cvt_test"
	"github.com/goplus/llcppg/_xtool/llcppsymg/clangutils"
)

func main() {
	TestForwardDecl()
	TestForwardDeclCrossFile()
}

func TestForwardDecl() {
	test.RunTestWithConfig(&clangutils.Config{
		File:  "./hfile/forwarddecl.h",
		Temp:  false,
		IsCpp: false,
	})
}

func TestForwardDeclCrossFile() {
	test.RunTestWithConfig(&clangutils.Config{
		File:  "./hfile/def.h",
		Temp:  false,
		IsCpp: false,
	})
}
