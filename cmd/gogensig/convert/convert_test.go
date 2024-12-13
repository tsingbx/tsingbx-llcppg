package convert_test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/goplus/llcppg/_xtool/llcppsymg/args"
	"github.com/goplus/llcppg/ast"
	"github.com/goplus/llcppg/cmd/gogensig/config"
	"github.com/goplus/llcppg/cmd/gogensig/convert"
	"github.com/goplus/llcppg/cmd/gogensig/convert/basic"
	"github.com/goplus/llcppg/cmd/gogensig/unmarshal"
	cppgtypes "github.com/goplus/llcppg/types"
	"github.com/goplus/llgo/xtool/env"
)

func init() {
	convert.SetDebug(convert.DbgFlagAll)
}

func TestFromTestdata(t *testing.T) {
	testFromDir(t, "./_testdata", false, "openssl")
}

// test sys type in stdinclude to package
func TestSysToPkg(t *testing.T) {
	name := "_systopkg"
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	testFrom(t, name, path.Join(dir, "_testdata", name), false, func(t *testing.T, pkg *convert.Package) {
		typConv := pkg.GetTypeConv()
		if typConv.SysTypeLoc == nil {
			t.Fatal("sysTypeLoc is nil")
		}
		pkgIncTypes := make(map[string]map[string][]string)

		// full type in all std lib
		for name, info := range typConv.SysTypeLoc {
			targetPkg, isDefault := convert.IncPathToPkg(info.IncPath)
			if isDefault {
				targetPkg = "github.com/goplus/llgo/c [default]"
			}
			if pkgIncTypes[targetPkg] == nil {
				pkgIncTypes[targetPkg] = make(map[string][]string, 0)
			}
			if pkgIncTypes[targetPkg][info.IncPath] == nil {
				pkgIncTypes[targetPkg][info.IncPath] = make([]string, 0)
			}
			pkgIncTypes[targetPkg][info.IncPath] = append(pkgIncTypes[targetPkg][info.IncPath], name)
		}

		for pkg, incTypes := range pkgIncTypes {
			t.Logf("\x1b[1;32m %s \x1b[0m Package contains inc types:", pkg)
			for incPath, types := range incTypes {
				t.Logf("\x1b[1;33m  - %s\x1b[0m (%s):", incPath, pkg)
				sort.Strings(types)
				t.Logf("    - %s", strings.Join(types, " "))
			}
		}

		// check referd type in std lib
		// Expected type to package mappings
		expected := map[string]string{
			"mbstate_t":   "github.com/goplus/llgo/c",
			"wint_t":      "github.com/goplus/llgo/c",
			"ptrdiff_t":   "github.com/goplus/llgo/c",
			"int8_t":      "github.com/goplus/llgo/c",
			"max_align_t": "github.com/goplus/llgo/c",
			"FILE":        "github.com/goplus/llgo/c",
			"tm":          "github.com/goplus/llgo/c/time",
			"time_t":      "github.com/goplus/llgo/c/time",
			"clock_t":     "github.com/goplus/llgo/c/time",
			"fenv_t":      "github.com/goplus/llgo/c/math",
			"size_t":      "github.com/goplus/llgo/c",
		}

		for name, exp := range expected {
			if _, ok := typConv.SysTypePkg[name]; ok {
				if typConv.SysTypePkg[name].PkgPath != exp {
					t.Errorf("type [%s]: expected package [%s], got [%s] in header [%s]", name, exp, typConv.SysTypePkg[name].PkgPath, typConv.SysTypePkg[name].Header.IncPath)
				} else {
					t.Logf("refer type [%s] expected package [%s] from header [%s]", name, exp, typConv.SysTypePkg[name].Header.IncPath)
				}
			} else {
				t.Logf("missing expected type %s (package: %s)", name, exp)
			}
		}
	}, []string{})
}

func TestDepPkg(t *testing.T) {
	name := "_depcjson"
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	testFrom(t, name, path.Join(dir, "_testdata", name), false, nil, []string{})
}

func testFromDir(t *testing.T, relDir string, gen bool, debugName string) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Getwd failed:", err)
	}
	dir = path.Join(dir, relDir)
	fis, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal("ReadDir failed:", err)
	}
	for _, fi := range fis {
		name := fi.Name()
		if strings.HasPrefix(name, "_") || strings.HasPrefix(name, ".") {
			continue
		}

		if len(debugName) > 0 && name != debugName {
			continue
		}

		t.Run(name, func(t *testing.T) {
			if len(debugName) > 0 {
				testFrom(t, name, dir+"/"+name, gen, nil, []string{"-v"})
			} else {
				testFrom(t, name, dir+"/"+name, gen, nil, []string{})
			}
		})
	}
}

func testFrom(t *testing.T, name, dir string, gen bool, validateFunc func(t *testing.T, pkg *convert.Package), llcppsigfetchArgs []string) {
	confPath := filepath.Join(dir, "conf")
	cfgPath := filepath.Join(confPath, "llcppg.cfg")
	symbPath := filepath.Join(confPath, "llcppg.symb.json")
	pubPath := filepath.Join(confPath, "llcppg.pub")
	expect := filepath.Join(dir, "gogensig.expect")
	var expectContent []byte
	if !gen {
		var err error
		expectContent, err = os.ReadFile(expect)
		if err != nil {
			t.Fatal(expectContent)
		}
	}

	cfg, err := config.GetCppgCfgFromPath(cfgPath)
	if err != nil {
		t.Fatal(err)
	}

	// origin cflags + test deps folder cflags,because the test deps 's cflags is depend on machine
	if cfg.CFlags != "" {
		cfg.CFlags = env.ExpandEnv(cfg.CFlags)
	}

	cfg.CFlags += " -I" + filepath.Join(dir, "hfile")
	flagedCfgPath, err := config.CreateTmpJSONFile(args.LLCPPG_CFG, cfg)
	defer os.Remove(flagedCfgPath)

	if err != nil {
		t.Fatal(err)
	}

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Fatal(err)
		}
	}()
	outputDir, err := ModInit(name)
	defer os.RemoveAll(outputDir)

	// patch the test file's cflags
	preprocess := func(p *convert.Package) {
		var patchFlags func(pkg *convert.PkgInfo)
		patchFlags = func(pkg *convert.PkgInfo) {
			if pkg.PkgPath != "." {
				incFlags := " -I" + filepath.Join(pkg.Dir, "hfile")
				pkg.CppgConf.CFlags += incFlags
				cfg.CFlags += incFlags
			}

			for _, dep := range pkg.Deps {
				patchFlags(dep)
				if err != nil {
					t.Fatal(err)
				}
			}
		}
		patchFlags(p.PkgInfo)
		err = config.CreateJSONFile(flagedCfgPath, cfg)
		if err != nil {
			t.Fatal(err)
		}
	}

	p, pkg, err := basic.ConvertProcesser(&basic.Config{
		PkgPreprocessor: preprocess,
		AstConvertConfig: convert.AstConvertConfig{
			PkgName:   name,
			SymbFile:  symbPath,
			CfgFile:   flagedCfgPath,
			OutputDir: outputDir,
			PubFile:   pubPath,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := config.SigfetchConfig(flagedCfgPath, confPath, llcppsigfetchArgs)
	if err != nil {
		t.Fatal(err)
	}

	inputdata, err := unmarshal.UnmarshalFileSet(bytes)
	if err != nil {
		t.Fatal(err)
	}

	err = p.ProcessFileSet(inputdata)
	if err != nil {
		t.Fatal(err)
	}

	var res strings.Builder

	outDir, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, fi := range outDir {
		if strings.HasSuffix(fi.Name(), "go.mod") || strings.HasSuffix(fi.Name(), "go.sum") || strings.HasSuffix(fi.Name(), "llcppg.pub") {
			continue
		} else {
			content, err := os.ReadFile(filepath.Join(outputDir, fi.Name()))
			if err != nil {
				t.Fatal(err)
			}
			res.WriteString(fmt.Sprintf("===== %s =====\n", fi.Name()))
			res.Write(content)
			res.WriteString("\n")
		}
	}

	pub, err := os.ReadFile(filepath.Join(outputDir, "llcppg.pub"))
	if err == nil {
		res.WriteString("===== llcppg.pub =====\n")
		res.Write(pub)
	}

	if gen {
		if err := os.WriteFile(expect, []byte(res.String()), 0644); err != nil {
			t.Fatal(err)
		}
	} else {
		expect := string(expectContent)
		got := res.String()
		if strings.TrimSpace(expect) != strings.TrimSpace(got) {
			t.Errorf("does not match expected.\nExpected:\n%s\nGot:\n%s", expect, got)
		}
	}

	if validateFunc != nil {
		validateFunc(t, pkg)
	}
}

// ===========================error
func TestNewAstConvert(t *testing.T) {
	_, err := convert.NewAstConvert(&convert.AstConvertConfig{
		PkgName:  "test",
		SymbFile: "",
		CfgFile:  "",
	})
	if err != nil {
		t.Fatal("NewAstConvert Fail")
	}
}

func TestNewAstConvertFail(t *testing.T) {
	_, err := convert.NewAstConvert(nil)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestVisitDone(t *testing.T) {
	pkg, err := convert.NewAstConvert(&convert.AstConvertConfig{
		PkgName:  "test",
		SymbFile: "",
		CfgFile:  "",
	})
	if err != nil {
		t.Fatal("NewAstConvert Fail")
	}
	pkg.SetVisitDone(func(pkg *convert.Package, incPath string) {
		if incPath != "test.h" {
			t.Fatal("doc path error")
		}
	})
	pkg.VisitDone("test.h")
}

func TestVisitFail(t *testing.T) {
	converter, err := convert.NewAstConvert(&convert.AstConvertConfig{
		PkgName:  "test",
		SymbFile: "",
		CfgFile:  "",
	})
	if err != nil {
		t.Fatal("NewAstConvert Fail")
	}

	// expect type
	converter.VisitTypedefDecl(&ast.TypedefDecl{
		Name: &ast.Ident{Name: "NormalType"},
		Type: &ast.BuiltinType{Kind: ast.Int},
	})

	// not appear in output,because expect error
	converter.VisitTypedefDecl(&ast.TypedefDecl{
		Name: &ast.Ident{Name: "Foo"},
		Type: nil,
	})

	errRecordType := &ast.RecordType{
		Tag: ast.Struct,
		Fields: &ast.FieldList{
			List: []*ast.Field{
				{Type: &ast.BuiltinType{Kind: ast.Int, Flags: ast.Double}},
			},
		},
	}
	// error field type for struct
	converter.VisitStruct(&ast.Ident{Name: "Foo"}, nil, &ast.TypeDecl{
		Name: &ast.Ident{Name: "Foo"},
		Type: errRecordType,
	})

	// error field type for anonymous struct
	converter.VisitStruct(&ast.Ident{Name: "Foo"}, nil, &ast.TypeDecl{
		Name: nil,
		Type: errRecordType,
	})

	converter.VisitEnumTypeDecl(&ast.EnumTypeDecl{
		Name: &ast.Ident{Name: "NormalType"},
		Type: &ast.EnumType{},
	})

	// error enum item for anonymous enum
	converter.VisitEnumTypeDecl(&ast.EnumTypeDecl{
		Name: nil,
		Type: &ast.EnumType{
			Items: []*ast.EnumItem{
				{Name: &ast.Ident{Name: "Item1"}},
			},
		},
	})

	converter.VisitFuncDecl(&ast.FuncDecl{
		Name: &ast.Ident{Name: "Foo"},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.BuiltinType{Kind: ast.Int, Flags: ast.Double}},
				},
			},
		},
	})
	// not appear in output

	buf, err := converter.Pkg.WriteDefaultFileToBuffer()
	if err != nil {
		t.Fatalf("WriteTo failed: %v", err)
	}

	expectedOutput :=
		`
package test

import (
	"github.com/goplus/llgo/c"
	_ "unsafe"
)

type NormalType c.Int

type Foo struct {
	Unused [8]uint8
}
`
	if strings.TrimSpace(expectedOutput) != strings.TrimSpace(buf.String()) {
		t.Errorf("does not match expected.\nExpected:\n%s\nGot:\n%s", expectedOutput, buf.String())
	}
}

func TestGetIncPathFail(t *testing.T) {
	cfg, err := config.CreateTmpJSONFile("llcppg.cfg", &cppgtypes.Config{
		Include: []string{"unexist.h"},
	})
	if err != nil {
		t.Fatal(err)
	}
	converter, err := convert.NewAstConvert(&convert.AstConvertConfig{
		PkgName:  "test",
		SymbFile: "",
		CfgFile:  cfg,
	})
	if err != nil {
		t.Fatal("NewAstConvert Fail")
	}
	converter.VisitStart("test.h", "", false)
}

func ModInit(name string) (string, error) {
	tempDir, err := os.MkdirTemp("", "gogensig-test")
	if err != nil {
		return "", err
	}
	outputDir := filepath.Join(tempDir, name)
	err = os.MkdirAll(outputDir, 0744)
	if err != nil {
		return "", err
	}
	projectRoot, err := filepath.Abs("../../../")
	if err != nil {
		return "", err
	}
	if err := os.Chdir(outputDir); err != nil {
		return "", err
	}

	err = config.RunCommand(outputDir, "go", "mod", "init", name)
	if err != nil {
		return "", err
	}
	err = config.RunCommand(outputDir, "go", "get", "github.com/goplus/llgo@main")
	if err != nil {
		return "", err
	}
	err = config.RunCommand(outputDir, "go", "get", "github.com/goplus/llcppg")
	if err != nil {
		return "", err
	}
	err = config.RunCommand(outputDir, "go", "mod", "edit", "-replace", "github.com/goplus/llcppg="+projectRoot)
	if err != nil {
		return "", err
	}
	return outputDir, nil
}
