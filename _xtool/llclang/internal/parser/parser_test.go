package parser_test

import (
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/goplus/lib/c"
	"github.com/goplus/llcppg/_xtool/llclang/internal/parser"
	"github.com/goplus/llpkg/cjson"
)

func TestParser(t *testing.T) {
	cases := []string{"class", "comment"}
	// https://github.com/goplus/llgo/issues/1114
	// todo(zzy):use os.ReadDir
	for _, folder := range cases {
		testFrom(t, filepath.Join("testdata", folder), "temp.h", false)
	}
}

func testFrom(t *testing.T, dir string, filename string, gen bool) {
	var expect string
	var err error
	if !gen {
		json, err := os.ReadFile(filepath.Join(dir, "expect.json"))
		if err != nil {
			t.Fatal("ReadExpectFile failed:", err)
		}
		expect = string(json)
	}
	ast, err := parser.Do(&parser.ConverterConfig{
		File:  filepath.Join(dir, filename),
		IsCpp: true,
		Args:  []string{"-fparse-all-comments"},
	})
	if err != nil {
		t.Fatal("Do failed:", err)
	}
	// https://github.com/goplus/llgo/issues/1116
	// astJson, err := json.MarshalIndent(ast, "", "  ")
	// todo(zzy):use json.Marshal
	if err != nil {
		t.Fatal("MarshalIndent failed:", err)
	}
	json := parser.MarshalASTFile(ast)
	output := json.Print()
	actual := c.GoString(output)
	defer cjson.FreeCStr(unsafe.Pointer(output))
	defer json.Delete()

	if gen {
		err = os.WriteFile(filepath.Join(dir, "expect.json"), []byte(actual), os.ModePerm)
		if err != nil {
			t.Fatal("WriteFile failed:", err)
		}
	} else if expect != actual {
		t.Fatal("expect != actual")
	}
}
