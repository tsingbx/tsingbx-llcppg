package convert

import (
	"testing"

	"github.com/goplus/gogen"
	"github.com/goplus/llcppg/ast"
	cfg "github.com/goplus/llcppg/cmd/gogensig/config"
	cppgtypes "github.com/goplus/llcppg/types"
)

func TestTypeRefIncompleteFail(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected panic, got nil")
		}
	}()
	pkg := NewPackage(&PackageConfig{
		PkgBase: PkgBase{
			PkgPath:  ".",
			CppgConf: &cppgtypes.Config{},
			Pubs:     make(map[string]string),
		},
		Name:        "testpkg",
		GenConf:     &gogen.Config{},
		OutputDir:   "",
		SymbolTable: cfg.CreateSymbolTable([]cfg.SymbolEntry{}),
	})
	pkg.cvt.SysTypeLoc["Bar"] = &HeaderInfo{
		IncPath: "Bar",
		Path:    "Bar",
	}
	pkg.incomplete["Bar"] = &gogen.TypeDecl{}
	err := pkg.NewTypedefDecl(&ast.TypedefDecl{
		Name: &ast.Ident{Name: "Foo"},
		Type: &ast.TagExpr{
			Name: &ast.Ident{Name: "Bar"},
		},
	})
	if err != nil {
		t.Fatal("NewTypedefDecl failed:", err)
	}
	delete(pkg.incomplete, "Bar")

	_, err = pkg.WriteToBuffer("testpkg")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	pkg.handleTyperefIncomplete(&ast.TagExpr{
		Tag: 0,
		Name: &ast.ScopingExpr{
			X: &ast.Ident{Name: "Bar"},
		},
	}, nil)
}
