package main

import (
	"github.com/goplus/llcppg/_xtool/llcppsigfetch/parse"
	test "github.com/goplus/llcppg/_xtool/llcppsigfetch/parse/cvt_test"
	"github.com/goplus/llcppg/llcppg"
)

func main() {
	TestForwardDecl()
	TestForwardDeclCrossFile()
}

func TestForwardDecl() {
	test.RunTestWithConfig(&parse.ParseConfig{
		Conf: &llcppg.Config{
			Include: []string{"forwarddecl.h"},
			CFlags:  "-I./hfile/",
		},
	}, true)
}

func TestForwardDeclCrossFile() {
	test.RunTestWithConfig(&parse.ParseConfig{
		Conf: &llcppg.Config{
			Include: []string{"def.h"},
			CFlags:  "-I./hfile/",
		},
	}, true)
}
