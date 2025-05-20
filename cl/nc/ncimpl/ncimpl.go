package ncimpl

import (
	"log"
	"strings"

	"github.com/goplus/llcppg/_xtool/llcppsymg/tool/name"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/nc"
	llconfig "github.com/goplus/llcppg/config"
)

type Converter struct {
	PkgName string

	TypeMap map[string]string // llcppg.pub
	FileMap map[string]*llconfig.FileInfo
	ConvSym func(name *ast.Object, mangleName string) (goName string, err error)

	// CfgFile   string // llcppg.cfg
	TrimPrefixes   []string
	KeepUnderScore bool
}

/* TODO(xsw): remove this
func (p *Converter) setCurFile(file string) {
	info, exist := p.FileMap[file]
	if !exist {
		var availableFiles []string
		for f := range p.FileMap {
			availableFiles = append(availableFiles, f)
		}
		log.Panicf("File %q not found in FileMap. Available files:\n%s",
			file, strings.Join(availableFiles, "\n"))
	}
	p.GenPkg.SetCurFile(NewHeaderFile(file, info.FileType))
}

func (p *Package) SetCurFile(hfile *HeaderFile) {
	var curFile *HeaderFile
	for _, f := range p.files {
		if f.File == hfile.File {
			curFile = f
			break
		}
	}

	if curFile == nil {
		curFile = hfile
		p.files = append(p.files, curFile)
	}

	p.curFile = curFile
	// for third hfile not register to gogen.Package
	if curFile.FileType != llcppg.Third {
		fileName := curFile.ToGoFileName(p.conf.Name)
		if debugLog {
			log.Printf("SetCurFile: %s File in Current Package: %v\n", fileName, curFile.FileType)
		}
		p.p.SetCurFile(fileName, true)
		p.p.Unsafe().MarkForceUsed(p.p)
	}
}
*/

func (p *Converter) convFile(file string) (goFile string, ok bool) {
	info, exist := p.FileMap[file]
	if !exist {
		var availableFiles []string
		for f := range p.FileMap {
			availableFiles = append(availableFiles, f)
		}
		log.Panicf("File %q not found in FileMap. Available files:\n%s",
			file, strings.Join(availableFiles, "\n"))
	}
	hf := NewHeaderFile(file, info.FileType)
	return hf.ToGoFileName(p.PkgName), hf.InCurPkg()
}

func (p *Converter) ConvDecl(file string, decl ast.Decl) (goName, goFile string, err error) {
	goFile, ok := p.convFile(file)
	if !ok {
		err = nc.ErrSkip
		return
	}
	obj := ast.ObjectOf(decl)
	if fn, ok := decl.(*ast.FuncDecl); ok {
		goName, err = p.ConvSym(obj, fn.MangledName)
	} else {
		goName = p.declName(obj.Name.Name)
	}
	return
}

func (p *Converter) ConvMacro(file string, macro *ast.Macro) (goName, goFile string, err error) {
	panic("todo")
}

func (p *Converter) ConvEnumItem(decl *ast.EnumTypeDecl, item *ast.EnumItem) (goName, goFile string, err error) {
	goName = p.constName(item.Name.Name)
	return
}

// which is define in llcppg.cfg/typeMap
func (p *Converter) definedName(name string) (string, bool) {
	definedName, ok := p.TypeMap[name]
	if ok {
		if definedName == "" {
			return name, true
		}
		return definedName, true
	}
	return name, false
}

type NameMethod func(name string) string

// transformName handles identifier name conversion following these rules:
// 1. First checks if the name exists in predefined mapping (in typeMap of llcppg.cfg)
// 2. If not in predefined mapping, applies the transform function
// 3. Before applying the transform function, removes specified prefixes (obtained via trimPrefixes)
//
// Parameters:
//   - name: Original C/C++ identifier name
//   - transform: Name transformation function (like names.PubName or names.ExportName)
//
// Returns:
//   - Transformed identifier name
func (p *Converter) transformName(cname string, transform NameMethod) string {
	if definedName, ok := p.definedName(cname); ok {
		return definedName
	}
	return transform(name.RemovePrefixedName(cname, p.trimPrefixes()))
}

func (p *Converter) declName(cname string) string {
	return p.transformName(cname, name.PubName)
}

func (p *Converter) constName(cname string) string {
	return p.transformName(cname, name.ExportName)
}

func (p *Converter) trimPrefixes() []string {
	return p.TrimPrefixes
}
