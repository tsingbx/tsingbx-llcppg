package main

import (
	"fmt"

	"github.com/goplus/llcppg/_xtool/llcppsymg/names"
	"github.com/goplus/llcppg/_xtool/llcppsymg/parse"
)

func main() {
	TestToGoName()
	TestNameMapper()
	TestPubName()
	TestExportName()
	TestHeaderFileToGo()
}

func TestToGoName() {
	fmt.Println("=== Test ToGoName ===")
	process1 := parse.NewSymbolProcessor([]string{}, []string{"lua_", "luaL_"})
	process2 := parse.NewSymbolProcessor([]string{}, []string{"sqlite3_", "sqlite3_"})
	process3 := parse.NewSymbolProcessor([]string{}, []string{"INI"})

	testCases := []struct {
		processor *parse.SymbolProcessor
		input     string
	}{
		{process1, "lua_closethread"},
		{process1, "luaL_checknumber"},
		{process2, "sqlite3_close_v2"},
		{process2, "sqlite3_callback"},
		{process3, "GetReal"},
		{process3, "GetBoolean"},
		{process3, "INIReader"},
	}

	for _, tc := range testCases {
		result := names.GoName(tc.input, tc.processor.Prefixes, true)
		fmt.Printf("Before: %s After: %s\n", tc.input, result)
	}
	fmt.Println()
}

func TestNameMapper() {
	fmt.Println("=== Test NameMapper ===")

	mapper := names.NewNameMapper()
	testCases := []struct {
		name         string
		trimPrefixes []string
		toCamel      bool
		expected     string
		expectChange bool
	}{
		{"lua_closethread", []string{"lua_", "luaL_"}, true, "Closethread", true},
		{"luaL_checknumber", []string{"lua_", "luaL_"}, true, "Checknumber", true},
		{"_gmp_err", []string{}, true, "X_gmpErr", true},
		{"fn_123illegal", []string{"fn_"}, true, "X123illegal", true},
		{"fts5_tokenizer", []string{}, true, "Fts5Tokenizer", true},
		{"Fts5Tokenizer", []string{}, true, "Fts5Tokenizer__1", true},
		{"normal_var", []string{}, false, "Normal_var", true},
		{"Cameled", []string{}, false, "Cameled", false},
	}

	fmt.Println("\nTesting GetUniqueGoName:")
	for _, tc := range testCases {
		result, changed := mapper.GetUniqueGoName(tc.name, tc.trimPrefixes, tc.toCamel)
		if result != tc.expected || changed != tc.expectChange {
			fmt.Printf("Input: %s, Expected: %s %t, Got: %s %t\n", tc.name, tc.expected, tc.expectChange, result, changed)
		} else {
			fmt.Printf("Input: %s, Output: %s %t\n", tc.name, result, changed)
		}
	}
}

func TestPubName() {
	fmt.Println("\n=== Test PubName ===")
	testCases := []struct {
		input    string
		expected string
	}{
		{"sqlite_file", "SqliteFile"},
		{"_gmp_err", "X_gmpErr"},
		{"123illegal", "X123illegal"},
		{"alreadyCamel", "AlreadyCamel"},
		{"_x__y", "X_xY"},
		{"_x_y", "X_xY"},
		{"_x___y", "X_xY"},
		{"x_y", "XY"},
		{"x__y", "XY"},
		{"x_y_", "XY"},
		{"x__y_", "XY"},
		{"_", "X_"},
		{"__", "X__"},
		{"___", "X___"},
	}

	for _, tc := range testCases {
		result := names.PubName(tc.input)
		if result != tc.expected {
			fmt.Printf("Input: %s, Expected: %s, Got: %s\n", tc.input, tc.expected, result)
		} else {
			fmt.Printf("Input: %s, Output: %s\n", tc.input, result)
		}
	}
}

func TestExportName() {
	fmt.Println("\n=== Test ExportName ===")
	testCases := []struct {
		input    string
		expected string
	}{
		{"sqlite_file", "Sqlite_file"},
		{"_sqlite_file", "X_sqlite_file"},
		{"123illegal", "X123illegal"},
		{"CODE_MASK", "CODE_MASK"},
		{"_CODE_MASK", "X_CODE_MASK"},
		{"_x__y", "X_x__y"},
		{"_x_y", "X_x_y"},
		{"_x___y", "X_x___y"},
		{"x_y", "X_y"},
		{"x__y", "X__y"},
		{"x_y_", "X_y_"},
		{"x__y_", "X__y_"},
		{"_", "X_"},
		{"__", "X__"},
		{"___", "X___"},
	}

	for _, tc := range testCases {
		result := names.ExportName(tc.input)
		if result != tc.expected {
			fmt.Printf("Input: %s, Expected: %s, Got: %s\n", tc.input, tc.expected, result)
		} else {
			fmt.Printf("Input: %s, Output: %s\n", tc.input, result)
		}
	}
}

func TestHeaderFileToGo() {
	fmt.Println("\n=== Test HeaderFileToGo ===")
	testCases := []struct {
		input    string
		expected string
	}{
		{"/path/to/foo.h", "foo.go"},
		{"/path/to/_intptr.h", "X_intptr.go"},
		{"header.h", "header.go"},
		{"_impl.h", "X_impl.go"},
	}

	for _, tc := range testCases {
		result := names.HeaderFileToGo(tc.input)
		if result != tc.expected {
			fmt.Printf("Input: %s, Expected: %s, Got: %s\n", tc.input, tc.expected, result)
		} else {
			fmt.Printf("Input: %s, Output: %s\n", tc.input, result)
		}
	}
}
