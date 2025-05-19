package nc

import (
	"errors"

	"github.com/goplus/llcppg/ast"
)

var (
	// ErrSkip is used to skip the node
	ErrSkip = errors.New("skip this node")
)

type NodeConverter interface {
	ConvDecl(decl ast.Decl) (goName, goFile string, err error) // ErrSkip
	ConvEnumItem(decl *ast.EnumTypeDecl, item *ast.EnumItem) (goName, goFile string, err error)
	ConvMacro(macro *ast.Macro) (goName, goFile string, err error)
}
