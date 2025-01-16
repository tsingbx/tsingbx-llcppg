package convert

import (
	"fmt"
	"math"
	"strconv"

	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cmd/gogensig/errs"
)

type ExprWrap struct {
	e ast.Expr
}

func Expr(e ast.Expr) *ExprWrap {
	return &ExprWrap{e: e}
}

func (p *ExprWrap) ToInt() (int, error) {
	v, ok := p.e.(*ast.BasicLit)
	if ok && v.Kind == ast.IntLit {
		v, _, err := litToInt(v.Value)
		if err != nil {
			return 0, err
		}
		return int(v), nil
	}
	return 0, errs.NewCantConvertError(p.e, "int")
}

func (p *ExprWrap) ToFloat(bitSize int) (float64, error) {
	v, ok := p.e.(*ast.BasicLit)
	if ok && v.Kind == ast.FloatLit {
		return litToFloat(v.Value, bitSize)
	}
	return 0, errs.NewCantConvertError(v, "float")
}

func (p *ExprWrap) ToString() (string, error) {
	v, ok := p.e.(*ast.BasicLit)
	if ok && v.Kind == ast.StringLit {
		return litToString(v.Value)
	}
	return "", errs.NewCantConvertError(v, "string")
}

func (p *ExprWrap) ToChar() (int8, error) {
	v, ok := p.e.(*ast.BasicLit)
	if ok && v.Kind == ast.CharLit {
		var iV int8
		_, err := fmt.Sscan(v.Value, &iV)
		if err == nil {
			return iV, nil
		}
	}
	return 0, errs.NewCantConvertError(p.e, "char")
}

func (p *ExprWrap) IsVoid() bool {
	retType, ok := p.e.(*ast.BuiltinType)
	if ok && retType.Kind == ast.Void {
		return true
	}
	return false
}

type IntType string

const (
	TypeInt   IntType = "Int"
	TypeUint  IntType = "Uint"
	TypeLong  IntType = "Long"
	TypeUlong IntType = "Ulong"
)

// litToInt parses a string literal into an int64 value.
// It returns the value and the smallest IntType that can hold the value.
func litToInt(lit string) (int64, IntType, error) {
	val, err := strconv.ParseInt(lit, 0, 64)
	if err != nil {
		return 0, TypeInt, err
	}

	if val < 0 {
		if val >= math.MinInt32 {
			return val, TypeInt, nil
		}
		return val, TypeLong, nil
	}
	return val, TypeInt, nil
}

// litToUint parses a string literal into a uint64 value.
// It returns the value and the smallest IntType that can hold the value.
func litToUint(lit string) (uint64, IntType, error) {
	uval, err := strconv.ParseUint(lit, 0, 64)
	if err != nil {
		return 0, TypeInt, err
	}
	switch {
	case uval <= math.MaxInt32:
		return uval, TypeInt, nil
	case uval <= math.MaxUint32:
		return uval, TypeUint, nil
	default:
		return uval, TypeUlong, nil
	}
}

func litToFloat(lit string, bitSize int) (float64, error) {
	return strconv.ParseFloat(lit, bitSize)
}

func litToString(lit string) (string, error) {
	return strconv.Unquote(lit)
}
