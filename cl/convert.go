package cl

import (
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/internal/convert"
	llconfig "github.com/goplus/llcppg/config"
)

const DbgFlagAll = convert.DbgFlagAll

func SetDebug(flag int) {
	convert.SetDebug(flag)
}

func ModInit(deps []string, outputDir string, modulePath string) error {
	return convert.ModInit(deps, outputDir, modulePath)
}

type ConvConfig struct {
	PkgName   string
	ConvSym   func(name *ast.Object, mangleName string) (goName string, err error)
	CfgFile   string // llcppg.cfg
	OutputDir string

	Pkg *llconfig.Pkg
}

func Convert(config *ConvConfig) (pkg Package, err error) {
	cvt, err := convert.NewConverter(&convert.Config{
		PkgName: config.PkgName,
		ConvSym: config.ConvSym,
		CfgFile: config.CfgFile,
		Pkg:     config.Pkg,
	})
	if err != nil {
		return
	}
	cvt.Convert()
	gp := cvt.GenPkg
	return Package{gp.Pkg(), gp.PkgInfo}, nil
}
