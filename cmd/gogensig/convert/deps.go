package convert

import (
	"fmt"
	"go/token"
	"go/types"
	"log"
	"path/filepath"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/_xtool/llcppsymg/args"
	"github.com/goplus/llcppg/_xtool/llcppsymg/config/cfgparse"
	cfg "github.com/goplus/llcppg/cmd/gogensig/config"
	"github.com/goplus/llcppg/cmd/gogensig/errs"
	cppgtypes "github.com/goplus/llcppg/types"
	"github.com/goplus/llgo/xtool/env"
	"github.com/goplus/mod/gopmod"
)

type Module = gopmod.Module

type PkgDepLoader struct {
	module   *gopmod.Module
	pkg      *gogen.Package
	pkgCache map[string]*PkgInfo // pkgPath -> *PkgInfo
	regCache map[string]struct{} // pkgPath
}

func NewPkgDepLoader(mod *gopmod.Module, pkg *gogen.Package) *PkgDepLoader {
	return &PkgDepLoader{
		module:   mod,
		pkg:      pkg,
		pkgCache: make(map[string]*PkgInfo),
		regCache: make(map[string]struct{}),
	}
}

// for current package & dependent packages
type PkgInfo struct {
	PkgBase
	Deps     []*PkgInfo
	Dir      string   // absolute local path of the package
	includes []string // abs header path
}

type PkgBase struct {
	PkgPath  string            // package path, e.g. github.com/goplus/llgo/cjson
	CppgConf *cppgtypes.Config // llcppg.cfg
	Pubs     map[string]string // llcppg.pub
}

func NewPkgInfo(pkgPath string, pkgDir string, conf *cppgtypes.Config, pubs map[string]string) *PkgInfo {
	return &PkgInfo{
		PkgBase: PkgBase{PkgPath: pkgPath, Pubs: pubs, CppgConf: conf},
		Dir:     pkgDir,
	}
}

// LoadDeps loads direct dependencies of the current package and recursively loads their
// dependencies, to get the complete dependency.
func (pm *PkgDepLoader) LoadDeps(p *PkgInfo) ([]*PkgInfo, error) {
	deps, err := pm.Imports(p.CppgConf.Deps)
	if err != nil {
		return nil, err
	}
	return deps, nil
}

func (pm *PkgDepLoader) Imports(pkgPaths []string) (pkgs []*PkgInfo, err error) {
	pkgs = make([]*PkgInfo, len(pkgPaths))
	for i, pkgPath := range pkgPaths {
		pkgs[i], err = pm.Import(pkgPath)
		if err != nil {
			return nil, err
		}
	}
	return
}

func (pm *PkgDepLoader) Import(pkgPath string) (*PkgInfo, error) {
	if pm.module == nil {
		return nil, errs.NewModNotFoundError()
	}
	if pkg, exist := pm.pkgCache[pkgPath]; exist {
		return pkg, nil
	}
	pkg, err := pm.module.Lookup(pkgPath)
	if err != nil {
		return nil, err
	}

	pkgDir, err := filepath.Abs(pkg.Dir)
	if err != nil {
		return nil, err
	}
	pubs, err := cfg.ReadPubFile(filepath.Join(pkgDir, args.LLCPPG_PUB))
	if err != nil {
		return nil, err
	}
	cfg, err := cfg.GetCppgCfgFromPath(filepath.Join(pkgDir, args.LLCPPG_CFG))
	if err != nil {
		return nil, err
	}
	newPkg := NewPkgInfo(pkgPath, pkgDir, cfg, pubs)
	pm.pkgCache[pkgPath] = newPkg

	if len(cfg.Deps) > 0 {
		deps, err := pm.LoadDeps(newPkg)
		newPkg.Deps = deps
		if err != nil {
			return nil, fmt.Errorf("failed to get deps for package %s: %w", pkgPath, err)
		}
	}
	return newPkg, nil
}

func (pm *PkgDepLoader) InitDeps(p *PkgInfo) error {
	deps, err := pm.LoadDeps(p)
	p.Deps = deps
	if err != nil {
		return err
	}
	pm.RegisterDeps(p)
	return nil
}

// RegisterDeps registers types from dependent packages into the current conversion project's scope
func (pm *PkgDepLoader) RegisterDeps(p *PkgInfo) {
	for _, dep := range p.Deps {
		pm.RegisterDep(dep)
	}
}

func (pm *PkgDepLoader) RegisterDep(dep *PkgInfo) {
	if _, ok := pm.regCache[dep.PkgPath]; ok {
		return
	}
	pm.regCache[dep.PkgPath] = struct{}{}
	genPkg := pm.pkg
	scope := genPkg.Types.Scope()
	depPkg := genPkg.Import(dep.PkgPath)
	pm.RegisterDeps(dep)
	for cName, pubGoName := range dep.Pubs {
		if pubGoName == "" {
			pubGoName = cName
		}
		if obj := depPkg.TryRef(pubGoName); obj != nil {
			var preObj types.Object
			if pubGoName == cName {
				preObj = obj
			} else {
				preObj = gogen.NewSubst(token.NoPos, genPkg.Types, cName, obj)
			}
			if old := scope.Insert(preObj); old != nil {
				log.Printf("conflicted name `%v` in %v, previous definition is %v\n", pubGoName, dep.PkgPath, old)
			}
		}
	}
}

func (p *PkgInfo) GetIncPaths() ([]string, []string, error) {
	if p.includes != nil {
		return p.includes, nil, nil
	}
	expandedIncFlags := env.ExpandEnv(p.CppgConf.CFlags)
	cflags := cfgparse.ParseCFlags(expandedIncFlags)
	incPaths, notFounds, err := cflags.GenHeaderFilePaths(p.CppgConf.Include)
	p.includes = incPaths
	return incPaths, notFounds, err
}
