package convert

import (
	"github.com/goplus/llcppg/_xtool/llcppsymg/names"
	"github.com/goplus/llcppg/llcppg"
)

type HeaderFile struct {
	File         string
	IncPath      string
	IsHeaderFile bool
	FileType     llcppg.FileType
	IsSys        bool
}

func (p *HeaderFile) ToGoFileName(pkgName string) string {
	if p.IsHeaderFile {
		switch p.FileType {
		case llcppg.Inter:
			return names.HeaderFileToGo(p.File)
		case llcppg.Impl:
			return pkgName + "_autogen.go"
		case llcppg.Third:
			// todo(zzy):ignore third file when dependency refactored
			if p.IsSys {
				return names.HeaderFileToGo(p.File)
			} else {
				// todo(zzy):temp gen third type to libname_autogen.go
				return pkgName + "_autogen.go"
			}
		default:
			panic("unkown FileType")
		}
	}
	// package name as the default file
	return p.File + ".go"
}

func (p *HeaderFile) InCurPkg() bool {
	return p.FileType == llcppg.Inter || p.FileType == llcppg.Impl
}

func NewHeaderFile(file string, incPath string, isHeaderFile bool, fileType llcppg.FileType, isSys bool) *HeaderFile {
	return &HeaderFile{
		File:         file,
		IncPath:      incPath,
		IsHeaderFile: isHeaderFile,
		FileType:     fileType,
		IsSys:        isSys,
	}
}
