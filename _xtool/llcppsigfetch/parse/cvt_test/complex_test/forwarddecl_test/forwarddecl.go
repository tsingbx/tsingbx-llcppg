package main

import (
	test "github.com/goplus/llcppg/_xtool/llcppsigfetch/parse/cvt_test"
	"github.com/goplus/llcppg/_xtool/llcppsymg/clangutils"
)

func main() {
	TestClassDecl()
}

func TestClassDecl() {
	test.RunTestWithConfig(&clangutils.Config{
		File:  "./hfile/forwarddecl.h",
		Temp:  false,
		IsCpp: false,
	})
}
