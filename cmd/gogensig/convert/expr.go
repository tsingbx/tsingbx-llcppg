package convert

import (
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
		return strconv.Atoi(v.Value)
	}
	return 0, errs.NewCantConvertError(p.e, "int")
}

func (p *ExprWrap) ToFloat(bitSize int) (float64, error) {
	v, ok := p.e.(*ast.BasicLit)
	if ok && v.Kind == ast.FloatLit {
		return strconv.ParseFloat(v.Value, bitSize)
	}
	return 0, errs.NewCantConvertError(v, "float")
}

func (p *ExprWrap) ToString() (string, error) {
	v, ok := p.e.(*ast.BasicLit)
	if ok && v.Kind == ast.StringLit {
		return v.Value, nil
	}
	return "", errs.NewCantConvertError(v, "string")
}

func (p *ExprWrap) ToChar() (int8, error) {
	v, ok := p.e.(*ast.BasicLit)
	if ok && v.Kind == ast.CharLit {
		iV, err := strconv.Atoi(v.Value)
		if err == nil {
			return int8(iV), nil
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
