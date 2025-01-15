package convert_test

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cmd/gogensig/convert"
)

func TestBasicLitFail(t *testing.T) {
	t.Parallel()
	type CaseType[T any] struct {
		name string
		expr ast.Expr
		want T
	}
	type CaseTypeSlice[T any] []CaseType[T]
	testCases := CaseTypeSlice[any]{
		{
			name: "ToInt",
			expr: &ast.BasicLit{Kind: ast.IntLit, Value: "ABC"},
			want: 123,
		},
		{
			name: "ToInt",
			expr: &ast.TagExpr{Tag: ast.Class, Name: &ast.Ident{Name: "Foo"}},
			want: 123,
		},
		{
			name: "ToFloat",
			expr: &ast.TagExpr{Tag: ast.Class, Name: &ast.Ident{Name: "Foo"}},
			want: 123.123,
		},
		{
			name: "ToString",
			expr: &ast.TagExpr{Tag: ast.Class, Name: &ast.Ident{Name: "Foo"}},
			want: "abcd",
		},
		{
			name: "ToChar",
			expr: &ast.TagExpr{Tag: ast.Class, Name: &ast.Ident{Name: "Foo"}},
			want: (int8)(98),
		},
	}

	for _, tc := range testCases {
		t.Run("convert "+tc.name, func(t *testing.T) {
			if tc.name == "ToInt" {
				_, err := convert.Expr(tc.expr).ToInt()
				expectError(t, err)
			} else if tc.name == "ToFloat" {
				_, err := convert.Expr(tc.expr).ToFloat(64)
				expectError(t, err)
			} else if tc.name == "ToChar" {
				_, err := convert.Expr(tc.expr).ToChar()
				expectError(t, err)
			} else if tc.name == "ToString" {
				_, err := convert.Expr(tc.expr).ToString()
				expectError(t, err)
			}
		})
	}
}

func TestBasicLitOK(t *testing.T) {
	t.Parallel()
	type CaseType[T any] struct {
		name string
		expr ast.Expr
		want T
	}
	type CaseTypeSlice[T any] []CaseType[T]
	testCases := CaseTypeSlice[any]{
		{
			name: "ToInt",
			expr: &ast.BasicLit{Kind: ast.IntLit, Value: "123"},
			want: 123,
		},
		{
			name: "ToFloat",
			expr: &ast.BasicLit{Kind: ast.FloatLit, Value: "123.123"},
			want: 123.123,
		},
		{
			name: "ToString",
			expr: &ast.BasicLit{Kind: ast.StringLit, Value: "\"abcd\""},
			want: "abcd",
		},
		{
			name: "ToChar",
			expr: &ast.BasicLit{Kind: ast.CharLit, Value: "98"},
			want: (int8)(98),
		},
	}

	for _, tc := range testCases {
		t.Run("convert "+tc.name, func(t *testing.T) {
			if tc.name == "ToInt" {
				result, err := convert.Expr(tc.expr).ToInt()
				checkResult(t, result, err, tc.want)
			} else if tc.name == "ToFloat" {
				result, err := convert.Expr(tc.expr).ToFloat(64)
				checkResult(t, result, err, tc.want)
			} else if tc.name == "ToChar" {
				result, err := convert.Expr(tc.expr).ToChar()
				checkResult(t, result, err, tc.want)
			} else if tc.name == "ToString" {
				result, err := convert.Expr(tc.expr).ToString()
				checkResult(t, result, err, tc.want)
			}
		})
	}
}

func TestLitToInt(t *testing.T) {
	type CaseType struct {
		lit  string
		want uint64
		typ  convert.IntType
	}

	int32Min := uint64(1) << 31
	int64Min := uint64(1) << 63
	testCases := []CaseType{
		{lit: "123", want: 123, typ: convert.TypeInt},
		{lit: "0xDEEDBEAF", want: 0xDEEDBEAF, typ: convert.TypeUint},                  // DEEDBEAF
		{lit: "0x80000000", want: int32Min, typ: convert.TypeUint},                    // INT32_MIN
		{lit: "0x7FFFFFFF", want: math.MaxInt32, typ: convert.TypeInt},                // INT32_MAX
		{lit: "0xFFFFFFFF", want: math.MaxUint32, typ: convert.TypeUint},              // UINT32_MAX
		{lit: "0xFFFFFFFFFFFFFFFF", want: 0xFFFFFFFFFFFFFFFF, typ: convert.TypeUlong}, // UINT64_MAX
		{lit: "-2147483648", want: int32Min, typ: convert.TypeInt},                    // INT32_MIN
		{lit: "2147483647", want: math.MaxInt32, typ: convert.TypeInt},                // INT32_MAX
		{lit: "-9223372036854775808", want: int64Min, typ: convert.TypeLong},          // INT64_MIN
		{lit: "9223372036854775807", want: math.MaxInt64, typ: convert.TypeUlong},     // INT64_MAX
	}
	for _, tc := range testCases {
		result, typ, err := convert.LitToInt(tc.lit)
		if err != nil {
			t.Error(err)
		}
		if tc.typ != typ {
			t.Error(tc.lit, "type mismatch want", tc.typ, "got", typ)
		}
		if result^tc.want != 0 || result != tc.want {
			t.Error(tc.lit, "result mismatch want", tc.want, "got", result, typ)
		}
	}
}

func expectError(t *testing.T, err error) {
	if err == nil {
		t.Error("expect error")
	}
}

func checkResult(t *testing.T, result any, err error, want any) {
	t.Helper()
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(result, want) {
		t.Error(cmp.Diff(result, want))
	}
}
