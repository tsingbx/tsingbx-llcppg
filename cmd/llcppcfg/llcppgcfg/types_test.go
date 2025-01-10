package llcppgcfg

import (
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
