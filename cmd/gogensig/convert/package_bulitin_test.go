package convert

import (
	"go/types"
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
	pkg.incomplete["Bar"] = &Incomplete{}
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

	err = pkg.WritePkgFiles()
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

func TestPubMethodName(t *testing.T) {
	name := types.NewTypeName(0, nil, "Foo", nil)
	named := types.NewNamed(name, nil, nil)
	ptrRecv := types.NewPointer(named)
	fnName := "Foo"
	pubName := pubMethodName(ptrRecv, &GoFuncSpec{GoSymbName: fnName, FnName: fnName, PtrRecv: true, IsMethod: true})
	if pubName != "(*Foo).Foo" {
		t.Fatal("Expected pubName to be '(*Foo).Foo', got", pubName)
	}
	valRecv := named
	pubName = pubMethodName(valRecv, &GoFuncSpec{GoSymbName: fnName, FnName: fnName, IsMethod: true})
	if pubName != "Foo.Foo" {
		t.Fatal("Expected pubName to be 'Foo.Foo', got", pubName)
	}

	unknownRecv := types.NewStruct(nil, []string{})
	pubName = pubMethodName(unknownRecv, &GoFuncSpec{GoSymbName: fnName, FnName: fnName, IsMethod: false})
	if pubName != "Foo" {
		t.Fatal("Expected pubName to be 'Foo', got", pubName)
	}
}

func TestGetNameType(t *testing.T) {
	named := types.NewNamed(types.NewTypeName(0, nil, "Foo", nil), nil, nil)
	ptrNamed := types.NewPointer(named)
	customSturct := types.NewStruct(nil, nil)

	namedRes := getNamedType(named)
	if namedRes != named {
		t.Fatal("Expected namedRes to be *types.Named, got", namedRes)
	}

	ptrNamedRes := getNamedType(ptrNamed)
	if ptrNamedRes != named {
		t.Fatal("Expected ptrNamedRes to be *types.Named, got", ptrNamedRes)
	}

	customRes := getNamedType(customSturct)
	if customRes != nil {
		t.Fatal("Expected nil, got", customRes)
	}
}
