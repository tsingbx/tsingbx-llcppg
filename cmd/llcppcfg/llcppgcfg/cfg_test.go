package llcppgcfg

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"
)

func Test_emptyStringError_Error(t *testing.T) {
	tests := []struct {
		name string
		p    *emptyStringError
		want string
	}{
		{
			"newEmptyStringError",
			newEmptyStringError("param"),
			"param can't be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.p.Error(); got != tt.want {
				t.Errorf("emptyStringError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newEmptyStringError(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *emptyStringError
	}{
		{
			"newEmptyStringError",
			args{
				"param",
			},
			&emptyStringError{"param"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.name, func(t *testing.T) {
			if got := newEmptyStringError(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEmptyStringError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isExcludeDir(t *testing.T) {
	type args struct {
		relPath        string
		excludeSubdirs []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"isExcludeDir",
			args{
				filepath.Join("a", "b", "c", "d"),
				[]string{"a"},
			},
			true,
		},
		{
			"isExcludeDir",
			args{
				filepath.Join("a", "b", "c", "d"),
				[]string{"b"},
			},
			false,
		},
		{
			"isExcludeDir",
			args{
				"ab",
				[]string{"a"},
			},
			false,
		},
		{
			"isExcludeDir",
			args{
				"a",
				[]string{"a"},
			},
			true,
		},
		{
			"isExcludeDir",
			args{
				filepath.Join("a", ""),
				[]string{"a"},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.relPath, func(t *testing.T) {
			if got := isExcludeDir(tt.args.relPath, tt.args.excludeSubdirs); got != tt.want {
				t.Errorf("isExcludeDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_doExpandCflags(t *testing.T) {
	cflags, _ := newCflags("cfg_test_data/libtasn1/include/")
	type args struct {
		str            string
		excludeSubdirs []string
		fn             func(s string) bool
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 string
	}{
		{
			"h",
			args{
				cflags,
				[]string{"internal"},
				func(s string) bool {
					ext := filepath.Ext(s)
					return ext == ".h"
				},
			},
			[]string{
				"libtasn1.h",
			},
			cflags,
		},
		{
			"hh",
			args{
				cflags,
				[]string{"internal"},
				func(s string) bool {
					ext := filepath.Ext(s)
					return ext == ".hh"
				},
			},
			[]string{
				"b.hh",
			},
			cflags,
		},
		{
			"hh&h",
			args{
				cflags,
				[]string{"internal"},
				func(s string) bool {
					ext := filepath.Ext(s)
					return ext == ".hh" || ext == ".h"
				},
			},
			[]string{
				"b.hh",
				"libtasn1.h",
			},
			cflags,
		},
		{
			"hh&h-no-excludes",
			args{
				cflags,
				[]string{},
				func(s string) bool {
					ext := filepath.Ext(s)
					return ext == ".hh" || ext == ".h"
				},
			},
			[]string{
				"b.hh",
				"internal/a.h",
				"libtasn1.h",
			},
			cflags,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := doExpandCflags(tt.args.str, tt.args.excludeSubdirs, tt.args.fn)
			sort.Strings(got)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("doExpandCflags() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("doExpandCflags() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestExpandName(t *testing.T) {
	type args struct {
		name         string
		dir          string
		libsOrCflags string
	}
	tests := []struct {
		name             string
		args             args
		wantExpandPrefix string
		wantOrg          string
	}{
		{
			"cflags",
			args{
				"libcjson",
				"",
				"cflags",
			},
			"-I/",
			"$(pkg-config --cflags libcjson)",
		},
		{
			"libs",
			args{
				"libcjson",
				"",
				"libs",
			},
			"-",
			"$(pkg-config --libs libcjson)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExpand, gotOrg := ExpandName(tt.args.name, tt.args.dir, tt.args.libsOrCflags)
			if !strings.HasPrefix(gotExpand, tt.wantExpandPrefix) {
				t.Errorf("ExpandName() gotExpand = %v, want %v", gotExpand, tt.wantExpandPrefix)
			}
			if gotOrg != tt.wantOrg {
				t.Errorf("ExpandName() gotOrg = %v, want %v", gotOrg, tt.wantOrg)
			}
		})
	}
}

func TestExpandLibsName(t *testing.T) {
	type args struct {
		name string
		dir  string
	}
	tests := []struct {
		name             string
		args             args
		wantExpandPrefix string
		wantOrg          string
	}{
		{
			"",
			args{
				"libcjson",
				"",
			},
			"-",
			"$(pkg-config --libs libcjson)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.name, func(t *testing.T) {
			gotExpand, gotOrg := ExpandLibsName(tt.args.name, tt.args.dir)
			if !strings.HasPrefix(gotExpand, tt.wantExpandPrefix) {
				t.Errorf("ExpandLibsName() gotExpand = %v, want %v", gotExpand, tt.wantExpandPrefix)
			}
			if gotOrg != tt.wantOrg {
				t.Errorf("ExpandLibsName() gotOrg = %v, want %v", gotOrg, tt.wantOrg)
			}
		})
	}
}

func TestExpandCflags(t *testing.T) {
	type args struct {
		originCFlags string
		exts         []string
		excludeDirs  []string
	}
	tests := []struct {
		name             string
		args             args
		wantIncludes     []string
		wantExpandPrefix string
		wantOrg          string
	}{
		{
			"libcjson",
			args{
				"$(pkg-config --cflags libcjson)",
				[]string{".h"},
				[]string{""},
			},
			[]string{"cJSON_Utils.h", "cJSON.h"},
			"-I/",
			"$(pkg-config --cflags libcjson)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIncludes, gotExpand, gotOrg := ExpandCflags(tt.args.originCFlags, tt.args.exts, tt.args.excludeDirs)
			sort.Strings(gotIncludes)
			sort.Strings(tt.wantIncludes)
			if !equalIncludes(gotIncludes, tt.wantIncludes) {
				t.Errorf("ExpandCflags() gotIncludes = %v, want %v", gotIncludes, tt.wantIncludes)

			}
			if !strings.HasPrefix(gotExpand, tt.wantExpandPrefix) {
				t.Errorf("ExpandCflags() gotExpand = %v, want %v", gotExpand, tt.wantExpandPrefix)
			}
			if gotOrg != tt.wantOrg {
				t.Errorf("ExpandCflags() gotOrg = %v, want %v", gotOrg, tt.wantOrg)
			}
		})
	}
}

func TestExpandCFlagsName(t *testing.T) {
	type args struct {
		name        string
		exts        []string
		excludeDirs []string
	}
	tests := []struct {
		name             string
		args             args
		wantIncludes     []string
		wantExpandPrefix string
		wantOrg          string
	}{
		{
			"libcjson",
			args{
				"libcjson",
				[]string{".h"},
				[]string{},
			},
			[]string{"cjson/cJSON_Utils.h", "cjson/cJSON.h"},
			"-I/",
			"$(pkg-config --cflags libcjson)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIncludes, gotExpand, gotOrg := ExpandCFlagsName(tt.args.name, tt.args.exts, tt.args.excludeDirs)
			sort.Strings(gotIncludes)
			sort.Strings(tt.wantIncludes)
			if !equalIncludes(gotIncludes, tt.wantIncludes) {
				t.Errorf("ExpandCFlagsName() gotIncludes = %v, want %v", gotIncludes, tt.wantIncludes)
			}
			if !strings.HasPrefix(gotExpand, tt.wantExpandPrefix) {
				t.Errorf("ExpandCFlagsName() gotExpand = %v, want %v", gotExpand, tt.wantExpandPrefix)
			}
			if gotOrg != tt.wantOrg {
				t.Errorf("ExpandCFlagsName() gotOrg = %v, want %v", gotOrg, tt.wantOrg)
			}
		})
	}
}

func Test_expandCFlagsAndLibs(t *testing.T) {
	config := NewLLCppConfig("libcjson", WithSort)
	type args struct {
		name        string
		cfg         *LLCppConfig
		dir         string
		exts        []string
		excludeDirs []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"libcjson",
			args{
				"libcjson",
				config,
				"",
				[]string{},
				[]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expandCFlagsAndLibs(tt.args.name, tt.args.cfg, tt.args.dir, tt.args.exts, tt.args.excludeDirs)
			if !strings.HasPrefix(config.CFlags, "-I") ||
				!strings.HasPrefix(config.Libs, "-") {
				t.Errorf("%s expand cflags and libs fail", tt.args.name)
			}
		})
	}
}

func Test_findDepSlice(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 string
	}{
		{
			"cjson",
			args{
				[]string{"cJSON_Utils.o:", "cjson/cJSON_Utils.h", "cjson/cJSON.h"},
			},
			[]string{"cjson/cJSON.h"},
			"cJSON_Utils.o:cjson/cJSON_Utils.h",
		},
		{
			"cjson",
			args{
				[]string{"cJSON_Utils.o:cjson/cJSON_Utils.h", "cjson/cJSON.h"},
			},
			[]string{"cjson/cJSON.h"},
			"cJSON_Utils.o:cjson/cJSON_Utils.h",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := findDepSlice(tt.args.lines)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findDepSlice() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("findDepSlice() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

type parseEntry struct {
	isDir bool
	name  string
	typ   fs.FileMode
}

func (p *parseEntry) Name() string {
	return p.name
}

func (p *parseEntry) IsDir() bool {
	return p.isDir
}

func (p *parseEntry) Type() fs.FileMode {
	return p.typ
}

func (p *parseEntry) Info() (os.FileInfo, error) {
	return nil, nil
}

func Test_parseFileEntry(t *testing.T) {
	cflags, inc := newCflags("cfg_test_data/cjson/include/")
	path := inc + "cJSON.h"
	internalPath := filepath.Join(inc, "internal")
	objFile := NewObjFileString("cJSON.o:cJSON.h")
	objFile.Deps = []string{}
	type args struct {
		cflags         string
		trimStr        string
		path           string
		d              fs.DirEntry
		exts           []string
		excludeSubdirs []string
	}
	tests := []struct {
		name string
		args args
		want *ObjFile
	}{
		{
			"cjson_ok",
			args{
				cflags,
				inc,
				path,
				&parseEntry{name: "cJSON.h"},
				[]string{".h"},
				[]string{},
			},
			objFile,
		},
		{
			"cjson_exts_not_found",
			args{
				cflags,
				inc,
				path,
				&parseEntry{name: "cJSON.h"},
				[]string{".hh"},
				[]string{},
			},
			nil,
		},
		{
			"cjson_not_rel_path",
			args{
				cflags,
				inc,
				"cJSON.h",
				&parseEntry{name: "cJSON.h"},
				[]string{".h"},
				[]string{},
			},
			objFile,
		},
		{
			"cjson_exclude_dir",
			args{
				cflags,
				inc,
				filepath.Join(internalPath, "a.h"),
				&parseEntry{name: "a.h"},
				[]string{".h"},
				[]string{"internal"},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseFileEntry(tt.args.cflags, tt.args.trimStr, tt.args.path, tt.args.d, tt.args.exts, tt.args.excludeSubdirs)
			if tt.want != nil && got != nil && !got.IsEqual(tt.want) {
				t.Errorf("parseFileEntry() = %v, want %v", got, tt.want)
			}
			if tt.want == nil {
				if got != nil {
					t.Errorf("parseFileEntry() = %v, want nil", got)
				}
			}
		})
	}
}

func Test_parseCFlagsEntry(t *testing.T) {
	cflags, inc := newCflags("cfg_test_data/cjson/include/")
	type args struct {
		cflags         string
		cflag          string
		exts           []string
		excludeSubdirs []string
	}
	tests := []struct {
		name string
		args args
		want *CflagEntry
	}{
		{
			"cjson_prefix_with_I",
			args{
				cflags,
				cflags,
				[]string{".h"},
				[]string{"internal"},
			},
			&CflagEntry{
				Include: inc,
				ObjFiles: []*ObjFile{
					{OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cJSON.h"}},
					{OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
				},
			},
		},
		{
			"cjson_not_prefix_with_I",
			args{
				cflags,
				inc,
				[]string{".h"},
				[]string{},
			},
			nil,
		},
		{
			"cjson_walk_dir_err",
			args{
				cflags,
				"-I/a/b/c",
				[]string{".h"},
				[]string{},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCFlagsEntry(tt.args.cflags, tt.args.cflag, tt.args.exts, tt.args.excludeSubdirs)
			if got != nil && tt.want != nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCFlagsEntry() = %v, want %v", got, tt.want)
			}
			if tt.want == nil {
				if got != nil {
					t.Errorf("parseCFlagsEntry() = %v, want nil", got)
				}
			}
		})
	}
}

func Test_sortIncludes(t *testing.T) {
	cflags, _ := newCflags("cfg_test_data/cjson/include/")
	depsCflags, _ := newCflags("cfg_test_data/deps/")
	cfg := &LLCppConfig{
		Name:   "libcjson",
		CFlags: cflags,
		Libs:   "$(pkg-config --libs libcjson)",
	}
	depCfg := &LLCppConfig{
		Name:   "deps",
		CFlags: depsCflags,
		Libs:   "",
	}
	type args struct {
		expandCflags   string
		cfg            *LLCppConfig
		exts           []string
		excludeSubdirs []string
	}
	tests := []struct {
		name        string
		args        args
		wantInclude []string
	}{
		{
			"cjson",
			args{
				cflags,
				cfg,
				[]string{".h"},
				[]string{"internal"},
			},
			[]string{
				"cJSON_Utils.h",
				"cJSON.h",
			},
		},
		{
			"deps/case1",
			args{
				filepath.Join(depsCflags, "case1"),
				depCfg,
				[]string{".h"},
				[]string{},
			},
			[]string{
				"a.h",
				"b.h",
			},
		},
		{
			"deps/case2",
			args{
				filepath.Join(depsCflags, "case2"),
				depCfg,
				[]string{".h"},
				[]string{},
			},
			[]string{
				"a.h",
				"b.h",
				"c.h",
			},
		},
		{
			"deps/case3_recircle",
			args{
				filepath.Join(depsCflags, "case3"),
				depCfg,
				[]string{".h"},
				[]string{},
			},
			[]string{
				"a.h",
				"b.h",
			},
		},
		{
			"deps/case4_recircle",
			args{
				filepath.Join(depsCflags, "case4"),
				depCfg,
				[]string{".h"},
				[]string{},
			},
			[]string{
				"a.h",
				"b.h",
				"f.h",
				"c.h",
				"e.h",
				"g.h",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortIncludes(tt.args.expandCflags, tt.args.cfg, tt.args.exts, tt.args.excludeSubdirs)
			if !reflect.DeepEqual(tt.args.cfg.Include, tt.wantInclude) {
				t.Errorf("sortIncludes() = %v, want %v", tt.args.cfg.Include, tt.wantInclude)
			}
		})
	}
}

func TestNewLLCppConfig(t *testing.T) {
	type args struct {
		name string
		flag FlagMode
	}
	tests := []struct {
		name string
		args args
		want *LLCppConfig
	}{
		{
			"libcjson",
			args{
				"libcjson",
				WithSort,
			},
			&LLCppConfig{
				Name:         "libcjson",
				CFlags:       "$(pkg-config --cflags libcjson)",
				Libs:         "$(pkg-config --libs libcjson)",
				Include:      nil,
				Deps:         nil,
				TrimPrefixes: []string{},
				Cplusplus:    false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLLCppConfig(tt.args.name, tt.args.flag); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLLCppConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizePackageName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"lua5.4",
			args{
				"lua5.4",
			},
			"lua5_4",
		},
		{
			"tree-sitter",
			args{
				"tree-sitter",
			},
			"tree_sitter",
		},
		{
			"python-3.12-embed",
			args{
				"python-3.12-embed",
			},
			"python_3_12_embed",
		},
		{
			"libmpdec++",
			args{
				"libmpdec++",
			},
			"libmpdec",
		},
		{
			"startWithDigit",
			args{
				"3json",
			},
			"_3json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizePackageName(tt.args.name); got != tt.want {
				t.Errorf("NormalizePackageName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenCfg(t *testing.T) {
	if runtime.GOOS == "linux" {
		return
	}
	_, bdwgcCfgFilePath := newCflags("cfg_test_data/bdw-gc/llcppg.cfg")
	_, libffiCfgFilePath := newCflags("cfg_test_data/libffi/llcppg.cfg")

	type args struct {
		name           string
		flag           FlagMode
		exts           []string
		excludeSubdirs []string
	}
	tests := []struct {
		name    string
		args    args
		want    *bytes.Buffer
		wantErr bool
	}{
		{
			"libcjson",
			args{
				"libcjson",
				WithSort,
				[]string{".h"},
				[]string{},
			},
			bytes.NewBufferString("{\n\t\"name\": \"libcjson\",\n\t\"cflags\": \"$(pkg-config --cflags libcjson)\",\n\t\"libs\": \"$(pkg-config --libs libcjson)\",\n\t\"include\": [\n\t\t\"cjson/cJSON_Utils.h\",\n\t\t\"cjson/cJSON.h\",\n\t\t\"cJSON_Utils.h\",\n\t\t\"cJSON.h\"\n\t],\n\t\"deps\": null,\n\t\"trimPrefixes\": [],\n\t\"cplusplus\": false\n}\n"),
			false,
		},
		{
			"bdw-gc",
			args{
				"bdw-gc",
				WithSort,
				[]string{".h"},
				[]string{},
			},
			readFile(bdwgcCfgFilePath),
			false,
		},
		{
			"libffi",
			args{
				"libffi",
				WithSort,
				[]string{".h"},
				[]string{},
			},
			readFile(libffiCfgFilePath),
			false,
		},
		{
			"empty_name",
			args{
				"",
				WithSort,
				[]string{".h"},
				[]string{},
			},
			nil,
			true,
		},
		{
			"expand",
			args{
				"libcjson",
				WithSort | WithExpand,
				[]string{".h"},
				[]string{},
			},
			nil,
			false,
		},
		{
			"expand_not_sort",
			args{
				"libcjson",
				WithExpand,
				[]string{".h"},
				[]string{},
			},
			nil,
			false,
		},
		{
			"normal_not_sort",
			args{
				"libcjson",
				0,
				[]string{".h"},
				[]string{},
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenCfg(tt.args.name, tt.args.flag, tt.args.exts, tt.args.excludeSubdirs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenCfg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.flag&WithExpand != 0 {
				if got.Len() <= 0 {
					t.Errorf("GenCfg() = %v, want expaned", got)
				}
			} else {
				if tt.args.flag&WithSort != 0 && !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GenCfg() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func newCflags(relDir string) (cflags string, path string) {
	dir, file := filepath.Split(relDir)
	wd, _ := os.Getwd()
	inc := joinPath(wd, dir)
	if len(file) > 0 {
		inc += file
	}
	trimStr := "-I" + inc
	return trimStr, inc
}

func joinPath(dir, rel string) string {
	path := filepath.Join(dir, rel)
	if !strings.HasSuffix(path, string(filepath.Separator)) {
		path += string(filepath.Separator)
	}
	return path
}

func equalIncludes(gotIncludes, wantIncludes []string) bool {
	if len(gotIncludes) != len(wantIncludes) {
		return false
	}
	for i := range gotIncludes {
		got := gotIncludes[i]
		want := wantIncludes[i]
		_, gotFile := filepath.Split(got)
		_, wantFile := filepath.Split(want)
		if gotFile != wantFile {
			return false
		}
	}
	return true
}

func readFile(filepath string) *bytes.Buffer {
	buf, err := os.ReadFile(filepath)
	if err != nil {
		return bytes.NewBufferString("")
	}
	return bytes.NewBuffer(buf)
}

func Test_getDir(t *testing.T) {
	type args struct {
		relPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"getDir/",
			args{
				relPath: "a/",
			},
			"a",
		},
		{
			"getDir",
			args{
				relPath: "a",
			},
			"a",
		},
		{
			"getDir/a/b/c/d",
			args{
				relPath: "a/b/c/d",
			},
			"a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDir(tt.args.relPath); got != tt.want {
				t.Errorf("getDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
