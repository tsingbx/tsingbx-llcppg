package ncimpl

import (
	"github.com/goplus/llcppg/ast"
	llconfig "github.com/goplus/llcppg/config"
)

type Converter struct {
	FileMap map[string]*llconfig.FileInfo
	ConvSym func(name *ast.Object, mangleName string) (goName string, err error)

	// CfgFile   string // llcppg.cfg
	TypeMap        map[string]string // llcppg.pub
	TrimPrefixes   []string
	KeepUnderScore bool
}
