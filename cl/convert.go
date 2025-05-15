package cl

import (
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
	SymbFile  string // llcppg.symb.json
	CfgFile   string // llcppg.cfg
	PubFile   string // llcppg.pub
	OutputDir string

	Pkg *llconfig.Pkg
}

func Convert(config *ConvConfig) (err error) {
	cvt, err := convert.NewConverter(&convert.Config{
		PkgName:  config.PkgName,
		SymbFile: config.SymbFile,
		CfgFile:  config.CfgFile,
		PubFile:  config.PubFile,
		Pkg:      config.Pkg,
	})
	if err != nil {
		return
	}
	cvt.Convert()
	return nil
}
