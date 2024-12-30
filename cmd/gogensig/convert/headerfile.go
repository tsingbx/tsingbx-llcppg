package convert

import (
	"github.com/goplus/llcppg/cmd/gogensig/convert/names"
)

type HeaderFile struct {
	File         string
	IncPath      string
	IsHeaderFile bool
	InCurPkg     bool
	IsSys        bool
}

func (p *HeaderFile) ToGoFileName() string {
	var fileName string
	if p.IsHeaderFile {
		// path to go filename
		fileName = names.HeaderFileToGo(p.File)
	} else {
		// package name as the default file
		fileName = p.File + ".go"
	}
	return fileName
}

func NewHeaderFile(file string, incPath string, isHeaderFile bool, inCurPkg bool, isSys bool) *HeaderFile {
	return &HeaderFile{
		File:         file,
		IncPath:      incPath,
		IsHeaderFile: isHeaderFile,
		InCurPkg:     inCurPkg,
		IsSys:        isSys,
	}
}
