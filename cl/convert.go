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
	OutputDir string

	Pkg *llconfig.Pkg

	// CfgFile   string // llcppg.cfg
	Pubs           map[string]string // llcppg.pub
	Deps           []string          // dependent packages
	TrimPrefixes   []string
	Libs           string
	KeepUnderScore bool
}

func Convert(config *ConvConfig) (pkg Package, err error) {
	cvt, err := convert.NewConverter(&convert.Config{
		PkgName:   config.PkgName,
		ConvSym:   config.ConvSym,
		OutputDir: config.OutputDir,
		Pkg:       config.Pkg,

		Pubs:           config.Pubs,
		Deps:           config.Deps,
		TrimPrefixes:   config.TrimPrefixes,
		Libs:           config.Libs,
		KeepUnderScore: config.KeepUnderScore,
	})
	if err != nil {
		return
	}
	cvt.Convert()
	gp := cvt.GenPkg
	return Package{gp.Pkg(), gp.PkgInfo}, nil
}
