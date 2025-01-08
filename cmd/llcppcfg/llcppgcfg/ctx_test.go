package llcppgcfg

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewObjFile(t *testing.T) {
	type args struct {
		oFile string
		hFile string
	}
	tests := []struct {
		name string
		args args
		want *ObjFile
	}{
		{
			"cJSON",
			args{
				"cJSON.o",
				"cJSON.h",
			},
			&ObjFile{
				OFile: "cJSON.o",
				HFile: "cJSON.h",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewObjFile(tt.args.oFile, tt.args.hFile); !got.IsEqual(tt.want) {
				t.Errorf("NewObjFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewObjFileString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want *ObjFile
	}{
		{
			"cjson_ok",
			args{
				"cJSON.o:cJSON.h",
			},
			&ObjFile{
				OFile: "cJSON.o",
				HFile: "cJSON.h",
			},
		}, {
			"cjson_fail",
			args{
				"cJSON.o cJSON.h",
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewObjFileString(tt.args.str)
			if got != nil && tt.want != nil && !got.IsEqual(tt.want) {
				t.Errorf("NewObjFileString() = %v, want %v", got, tt.want)
			}
			if tt.want == nil {
				if got != nil {
					t.Errorf("NewObjFileString() = %v, want nil", got)
				}
			}
		})
	}
}

func TestObjFile_String(t *testing.T) {
	tests := []struct {
		name string
		o    *ObjFile
		want string
	}{
		{
			"cjson",
			&ObjFile{
				OFile: "cJSON.o",
				HFile: "cJSON.h",
			},
			"{OFile:cJSON.o, HFile:cJSON.h, Deps:[]}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.String(); got != tt.want {
				t.Errorf("ObjFile.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCflagEntry_String(t *testing.T) {
	tests := []struct {
		name string
		c    *CflagEntry
		want string
	}{
		{
			"cjson",
			&CflagEntry{
				Include: "/user/local/opt/cjson/include",
				ObjFiles: []*ObjFile{
					{
						OFile: "cJSON.o",
						HFile: "cJSON.h",
					},
					{
						OFile: "cJSON_Utils.o",
						HFile: "cJSON_Utils.h",
					},
				},
			},
			"{Include:/user/local/opt/cjson/include, ObjFiles:[{OFile:cJSON.o, HFile:cJSON.h, Deps:[]} {OFile:cJSON_Utils.o, HFile:cJSON_Utils.h, Deps:[]}]}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("CflagEntry.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDepCtx(t *testing.T) {
	_, inc := newCflags("cfg_test_data/cjson/include")
	type args struct {
		cflagEntry *CflagEntry
	}
	tests := []struct {
		name string
		args args
		want *DepCtx
	}{
		{
			"cjson",
			args{
				&CflagEntry{
					inc,
					[]*ObjFile{
						{OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
						{OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
					},
				},
			},
			&DepCtx{
				&CflagEntry{
					inc,
					[]*ObjFile{
						{OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
						{OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
					},
				},
				map[int]*ObjFile{
					0: {OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
					1: {OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
				},
				map[string]int{
					"cJSON_Utils.h": 0,
					"cJSON.h":       1,
				},
				map[*ObjFile][]int{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDepCtx(tt.args.cflagEntry); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDepCtx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDepCtx_GetObjFileByRelPath(t *testing.T) {
	_, inc := newCflags("cfg_test_data/cjson/include")
	type args struct {
		relPath string
	}
	tests := []struct {
		name  string
		p     *DepCtx
		args  args
		want  *ObjFile
		want1 int
	}{
		{
			"cjson_exist",
			&DepCtx{
				&CflagEntry{
					inc,
					[]*ObjFile{
						{OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
						{OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
					},
				},
				map[int]*ObjFile{
					0: {OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
					1: {OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
				},
				map[string]int{
					"cJSON_Utils.h": 0,
					"cJSON.h":       1,
				},
				map[*ObjFile][]int{},
			},
			args{"cJSON_Utils.h"},
			&ObjFile{OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
			0,
		},
		{
			"cjson_not_exist",
			&DepCtx{
				&CflagEntry{
					inc,
					[]*ObjFile{
						{OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
						{OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
					},
				},
				map[int]*ObjFile{
					0: {OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
					1: {OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
				},
				map[string]int{
					"cJSON_Utils.h": 0,
					"cJSON.h":       1,
				},
				map[*ObjFile][]int{},
			},
			args{"cJSON_Utilss.h"},
			nil,
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.p.GetObjFileByRelPath(tt.args.relPath)
			if got != nil && tt.want != nil {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("DepCtx.GetObjFileByRelPath() got = %v, want %v", got, tt.want)
				}
			}
			if tt.want == nil {
				if got != nil {
					t.Errorf("DepCtx.GetObjFileByRelPath() got = %v, want nil", got)
				}
			}
			if got1 != tt.want1 {
				t.Errorf("DepCtx.GetObjFileByRelPath() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDepCtx_GetObjFileByID(t *testing.T) {
	_, inc := newCflags("cfg_test_data/cjson/include")
	type args struct {
		id int
	}
	tests := []struct {
		name string
		p    *DepCtx
		args args
		want *ObjFile
	}{
		{
			"exist",
			&DepCtx{
				&CflagEntry{
					inc,
					[]*ObjFile{
						{OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
						{OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
					},
				},
				map[int]*ObjFile{
					0: {OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
					1: {OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
				},
				map[string]int{
					"cJSON_Utils.h": 0,
					"cJSON.h":       1,
				},
				map[*ObjFile][]int{},
			},
			args{0},
			&ObjFile{OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
		},
		{
			"not_exist",
			&DepCtx{
				&CflagEntry{
					inc,
					[]*ObjFile{
						{OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
						{OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
					},
				},
				map[int]*ObjFile{
					0: {OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
					1: {OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
				},
				map[string]int{
					"cJSON_Utils.h": 0,
					"cJSON.h":       1,
				},
				map[*ObjFile][]int{},
			},
			args{3},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.GetObjFileByID(tt.args.id)
			if got != nil && tt.want != nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DepCtx.GetObjFileByID() = %v, want %v", got, tt.want)
			}
			if tt.want == nil {
				if got != nil {
					t.Errorf("DepCtx.GetObjFileByID() = %v, want nil", got)
				}
			}
		})
	}
}

func TestDepCtx_GetIDByRelPath(t *testing.T) {
	_, inc := newCflags("cfg_test_data/cjson/include")
	type args struct {
		relPath string
	}
	tests := []struct {
		name string
		p    *DepCtx
		args args
		want int
	}{
		{
			"id>=0",
			&DepCtx{
				&CflagEntry{
					inc,
					[]*ObjFile{
						{OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
						{OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
					},
				},
				map[int]*ObjFile{
					0: {OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
					1: {OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
				},
				map[string]int{
					"cJSON_Utils.h": 0,
					"cJSON.h":       1,
				},
				map[*ObjFile][]int{},
			},
			args{"cJSON_Utils.h"},
			0,
		},
		{
			"id==-1",
			&DepCtx{
				&CflagEntry{
					inc,
					[]*ObjFile{
						{OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
						{OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
					},
				},
				map[int]*ObjFile{
					0: {OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
					1: {OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
				},
				map[string]int{
					"cJSON_Utils.h": 0,
					"cJSON.h":       1,
				},
				map[*ObjFile][]int{},
			},
			args{"cJSON_Utilss.h"},
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GetIDByRelPath(tt.args.relPath); got != tt.want {
				t.Errorf("DepCtx.GetIDByRelPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDepCtx_GetInclude(t *testing.T) {
	_, inc := newCflags("cfg_test_data/cjson/include")
	tests := []struct {
		name string
		p    *DepCtx
		want string
	}{
		{
			"cjson",
			&DepCtx{
				&CflagEntry{
					inc,
					[]*ObjFile{
						{OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
						{OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
					},
				},
				map[int]*ObjFile{
					0: {OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cjson/cJSON.h"}},
					1: {OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
				},
				map[string]int{
					"cJSON_Utils.h": 0,
					"cJSON.h":       1,
				},
				map[*ObjFile][]int{},
			},
			inc,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GetInclude(); got != tt.want {
				t.Errorf("DepCtx.GetInclude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDepCtx_ExpandDeps(t *testing.T) {
	_, inc := newCflags("cfg_test_data/cjson/include")
	type args struct {
		objFile *ObjFile
	}
	tests := []struct {
		name        string
		p           *DepCtx
		args        args
		wantDepsMap map[*ObjFile][]int
	}{
		{
			"cjson",
			&DepCtx{
				&CflagEntry{
					inc,
					[]*ObjFile{
						{OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cJSON.h"}},
						{OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
					},
				},
				map[int]*ObjFile{
					0: {OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cJSON.h"}},
					1: {OFile: "cJSON.o", HFile: "cJSON.h", Deps: []string{}},
				},
				map[string]int{
					"cJSON_Utils.h": 0,
					"cJSON.h":       1,
				},
				map[*ObjFile][]int{},
			},
			args{
				&ObjFile{OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cJSON.h"}},
			},
			map[*ObjFile][]int{
				{OFile: "cJSON_Utils.o", HFile: "cJSON_Utils.h", Deps: []string{"cJSON.h"}}: {1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.ExpandDeps(tt.args.objFile)
			got := fmt.Sprintf("%v", tt.p.depsMap)
			want := fmt.Sprintf("%v", tt.wantDepsMap)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("ExpandDeps() = %v, want %v", tt.p.depsMap, tt.wantDepsMap)
			}
		})
	}
}

func Test_removeDups(t *testing.T) {
	type args[TT comparable] struct {
		s []TT
	}
	type CaseType[TT comparable] struct {
		name string
		args args[TT]
		want []TT
	}
	tests := []CaseType[int]{
		{
			"ints",
			args[int]{
				[]int{1, 1, 7, 7, 7, 2, 3, 5, 6, 6},
			},
			[]int{1, 7, 2, 3, 5, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeDups(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeDups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjFile_IsEqual(t *testing.T) {
	type args struct {
		o *ObjFile
	}
	tests := []struct {
		name string
		p    *ObjFile
		args args
		want bool
	}{
		{
			"hfile_not_equal",
			&ObjFile{
				OFile: "a.o",
				HFile: "a.h",
				Deps:  []string{},
			},
			args{
				&ObjFile{
					OFile: "a.o",
					HFile: "b.h",
					Deps:  []string{},
				},
			},
			false,
		},
		{
			"ofile_not_equal",
			&ObjFile{
				OFile: "a.o",
				HFile: "a.h",
				Deps:  []string{},
			},
			args{
				&ObjFile{
					OFile: "b.o",
					HFile: "a.h",
					Deps:  []string{},
				},
			},
			false,
		},
		{
			"dep_len_not_equal",
			&ObjFile{
				OFile: "a.o",
				HFile: "a.h",
				Deps:  []string{},
			},
			args{
				&ObjFile{
					OFile: "a.o",
					HFile: "a.h",
					Deps:  []string{"c.h"},
				},
			},
			false,
		},
		{
			"dep_not_equal",
			&ObjFile{
				OFile: "a.o",
				HFile: "a.h",
				Deps:  []string{"b.h"},
			},
			args{
				&ObjFile{
					OFile: "a.o",
					HFile: "a.h",
					Deps:  []string{"c.h"},
				},
			},
			false,
		},
		{
			"equal",
			&ObjFile{
				OFile: "a.o",
				HFile: "a.h",
				Deps:  []string{"b.h"},
			},
			args{
				&ObjFile{
					OFile: "a.o",
					HFile: "a.h",
					Deps:  []string{"b.h"},
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.IsEqual(tt.args.o); got != tt.want {
				t.Errorf("ObjFile.IsEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
