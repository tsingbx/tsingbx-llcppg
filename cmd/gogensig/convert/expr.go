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
		v, _, err := LitToInt(v.Value)
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

func LitToInt(lit string) (uint64, IntType, error) {
	var val int64
	var uval uint64
	var err error
	if uval, err = strconv.ParseUint(lit, 0, 64); err == nil {
		switch {
		case uval <= math.MaxInt32:
			return uval, TypeInt, nil
		case uval <= math.MaxUint32:
			return uval, TypeUint, nil
		default:
			return uval, TypeUlong, nil
		}
	}

	// handle negative numbers
	if val, err = strconv.ParseInt(lit, 0, 64); err == nil {
		if val < 0 {
			if val >= math.MinInt32 {
				// For negative numbers, preserve only the lower 32 bits
				return uint64(uint32(int32(val))), TypeInt, nil
			}
			return uint64(val), TypeLong, nil
		}
	}
	return 0, TypeInt, err
}

func litToFloat(lit string, bitSize int) (float64, error) {
	return strconv.ParseFloat(lit, bitSize)
}

func litToString(lit string) (string, error) {
	return strconv.Unquote(lit)
}
