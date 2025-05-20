package ncimpl

import (
	"log"
	"strings"

	"github.com/goplus/llcppg/_xtool/llcppsymg/tool/name"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cl/nc"
	llconfig "github.com/goplus/llcppg/config"
)

type ThirdTypeLoc struct {
	locMap map[string]string // type name from third package -> define location
}

/* TODO(xsw): remove this
func NewThirdTypeLoc() *ThirdTypeLoc {
	return &ThirdTypeLoc{
		locMap: make(map[string]string),
	}
}
*/

func (p *ThirdTypeLoc) Add(ident *ast.Ident, loc *ast.Location) {
	if p.locMap == nil {
		p.locMap = make(map[string]string)
	}
	if _, ok := p.locMap[ident.Name]; ok {
		// a third ident in multiple location is permit
		return
	}
	p.locMap[ident.Name] = loc.File
}

func (p *ThirdTypeLoc) Lookup(name string) (string, bool) {
	loc, ok := p.locMap[name]
	return loc, ok
}

type Converter struct {
	PkgName string

	TypeMap map[string]string // llcppg.pub
	FileMap map[string]*llconfig.FileInfo
	ConvSym func(name *ast.Object, mangleName string) (goName string, err error)

	locMap ThirdTypeLoc // record third type's location

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

// typedecl,enumdecl,funcdecl,funcdecl
// true determine continue execute the type gen
// if this type is in a third header,skip the type gen & collect the type info
func (p *Package) handleType(ident *ast.Ident, loc *ast.Location) (skip bool) {
	if curPkg := p.curFile.InCurPkg(); curPkg {
		return false
	}
	if _, ok := p.locMap.Lookup(ident.Name); ok {
		// a third ident in multiple location is permit
		return true
	}
	p.locMap.Add(ident, loc)
	return true
}
*/

func (p *Converter) convFile(file string, obj *ast.Object) (goFile string, ok bool) {
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
	if obj != nil && hf.FileType == llconfig.Third {
		p.locMap.Add(obj.Name, obj.Loc)
	}
	return hf.ToGoFileName(p.PkgName), hf.InCurPkg()
}

func (p *Converter) ConvDecl(file string, decl ast.Decl) (goName, goFile string, err error) {
	obj := ast.ObjectOf(decl)
	goFile, ok := p.convFile(file, obj)
	if !ok {
		err = nc.ErrSkip
		return
	}
	switch decl := decl.(type) {
	case *ast.FuncDecl:
		goName, err = p.ConvSym(obj, decl.MangledName)
		// only have error when symbol not found,current keep only log this error
		if err != nil {
			log.Printf("ConvDecl: %s not found in symbolmap: %s", decl.MangledName, err.Error())
			err = nc.ErrSkip
			return
		}
	case *ast.EnumTypeDecl:
		// support anonymous enum with empty name
		if obj.Name != nil {
			goName = p.declName(obj.Name.Name)
		}
	default:
		goName = p.declName(obj.Name.Name)
	}
	return
}

func (p *Converter) ConvMacro(file string, macro *ast.Macro) (goName, goFile string, err error) {
	goFile, ok := p.convFile(file, nil)
	if !ok {
		err = nc.ErrSkip
		return
	}
	goName = p.constName(macro.Name)
	return
}

func (p *Converter) ConvEnumItem(decl *ast.EnumTypeDecl, item *ast.EnumItem) (goName string, err error) {
	goName = p.constName(item.Name.Name)
	return
}

func (p *Converter) ConvTagExpr(cname string) string {
	return p.declName(cname)
}

func (p *Converter) Lookup(name string) (locFile string, ok bool) {
	return p.locMap.Lookup(name)
}

func (p *Converter) IsPublic(cname string) bool {
	return p.KeepUnderScore || rune(cname[0]) != '_'
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
