package llcppgcfg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/goplus/llcppg/cmdout"
	"github.com/goplus/llcppg/types"
)

type llcppCfgKey string

const (
	cfgLibsKey   llcppCfgKey = "libs"
	cfgCflagsKey llcppCfgKey = "cflags"
)

type FlagMode int

const (
	WithTab FlagMode = 1 << iota
	WithCpp
)

type emptyStringError struct {
	name string
}

func (p *emptyStringError) Error() string {
	return p.name + " can't be empty"
}

func newEmptyStringError(name string) *emptyStringError {
	return &emptyStringError{name: name}
}

type LLCppConfig types.Config

func getDir(relPath string) string {
	index := strings.IndexRune(relPath, filepath.Separator)
	if index < 0 {
		return relPath
	}
	return relPath[:index]
}

func isExcludeDir(relPath string, excludeSubdirs []string) bool {
	if len(excludeSubdirs) == 0 {
		return false
	}
	dir := getDir(relPath)
	for _, subdir := range excludeSubdirs {
		if subdir == dir {
			return true
		}
	}
	return false
}

func ExpandName(name string, dir string, cfgKey llcppCfgKey) string {
	originString := fmt.Sprintf("$(pkg-config --%s %s)", cfgKey, name)
	return cmdout.ExpandString(originString, dir)
}

func findDepSlice(lines []string) ([]string, string) {
	objFileString := ""
	iStart := 0
	numLines := len(lines)
	complete := false
	for i := 0; i < numLines && !complete; i++ {
		line := lines[i]
		if strings.ContainsRune(line, rune(':')) && !strings.HasSuffix(line, ":") {
			objFileString = line
			iStart = i + 1
			break
		}
		complete = true
		for j := i + 1; j < numLines; j++ {
			line2 := lines[j]
			if len(line2) > 0 {
				iStart = j + 1
				objFileString = line + line2
				break
			}
		}
	}
	if iStart < numLines {
		return lines[iStart:], objFileString
	}
	return []string{}, objFileString
}

func getClangArgs(cflags string, relpath string) []string {
	args := make([]string, 0)
	cflagsField := strings.Fields(cflags)
	args = append(args, cflagsField...)
	args = append(args, "-MM")
	args = append(args, relpath)
	return args
}

func parseFileEntry(cflags, trimCflag, path string, d fs.DirEntry, exts []string, excludeSubdirs []string) *ObjFile {
	if d.IsDir() || strings.HasPrefix(d.Name(), ".") {
		return nil
	}
	idx := len(exts)
	for i, ext := range exts {
		if strings.HasSuffix(d.Name(), ext) {
			idx = i
			break
		}
	}
	if idx == len(exts) {
		return nil
	}
	relPath, err := filepath.Rel(trimCflag, path)
	if err != nil {
		relPath = path
	}
	if isExcludeDir(relPath, excludeSubdirs) {
		return nil
	}
	args := getClangArgs(cflags, relPath)
	clangCmd := cmdout.NewExecCommand("clang", args...)
	outString, err := cmdout.GetOut(clangCmd, trimCflag)
	if err != nil || outString == "" {
		objFile := NewObjFile(relPath, relPath)
		return objFile
	}
	outString = strings.ReplaceAll(outString, "\\\n", "\n")
	fields := strings.Fields(outString)
	lines, objFileStr := findDepSlice(fields)
	objFile := NewObjFileString(objFileStr)
	objFile.Deps = append(objFile.Deps, lines...)
	return objFile
}

func parseCFlagsEntry(cflags, cflag string, exts []string, excludeSubdirs []string) *CflagEntry {
	if !strings.HasPrefix(cflag, "-I") {
		return nil
	}
	trimCflag := strings.TrimPrefix(cflag, "-I")
	if !strings.HasSuffix(trimCflag, string(filepath.Separator)) {
		trimCflag += string(filepath.Separator)
	}
	var cflagEntry CflagEntry
	cflagEntry.Include = trimCflag
	cflagEntry.ObjFiles = make([]*ObjFile, 0)
	err := filepath.WalkDir(trimCflag, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		pObjFile := parseFileEntry(cflags, trimCflag, path, d, exts, excludeSubdirs)
		if pObjFile != nil {
			cflagEntry.ObjFiles = append(cflagEntry.ObjFiles, pObjFile)
		}
		return nil
	})
	sort.Slice(cflagEntry.ObjFiles, func(i, j int) bool {
		return len(cflagEntry.ObjFiles[i].Deps) > len(cflagEntry.ObjFiles[j].Deps)
	})
	if err != nil {
		return nil
	}
	return &cflagEntry
}

func sortIncludes(expandCflags string, cfg *LLCppConfig, exts []string, excludeSubdirs []string) {
	list := strings.Fields(expandCflags)
	cflagEntryList := make([]*CflagEntry, 0)
	for _, cflag := range list {
		pCflagEntry := parseCFlagsEntry(expandCflags, cflag, exts, excludeSubdirs)
		if pCflagEntry != nil {
			cflagEntryList = append(cflagEntryList, pCflagEntry)
		}
	}
	cfg.Include = make([]string, 0)
	for _, cflagEntry := range cflagEntryList {
		for _, objFile := range cflagEntry.ObjFiles {
			cfg.Include = append(cfg.Include, objFile.HFile)
		}
	}
}

func NewLLCppConfig(name string, flag FlagMode) *LLCppConfig {
	cfg := &LLCppConfig{
		Name: name,
	}
	cfg.CFlags = fmt.Sprintf("$(pkg-config --cflags %s)", name)
	cfg.Libs = fmt.Sprintf("$(pkg-config --libs %s)", name)
	cfg.TrimPrefixes = []string{}
	cfg.Cplusplus = (flag&WithCpp != 0)
	return cfg
}

func NormalizePackageName(name string) string {
	fields := strings.FieldsFunc(name, func(r rune) bool {
		return !unicode.IsLetter(r) && r != '_' && !unicode.IsDigit(r)
	})
	if len(fields) > 0 {
		if len(fields[0]) > 0 && unicode.IsDigit(rune(fields[0][0])) {
			fields[0] = "_" + fields[0]
		}
	}
	return strings.Join(fields, "_")
}

func GenCfg(name string, flag FlagMode, exts []string, excludeSubdirs []string) (*bytes.Buffer, error) {
	if len(name) == 0 {
		return nil, newEmptyStringError("name")
	}
	cfg := NewLLCppConfig(name, flag)
	expandCFlags := ExpandName(name, "", cfgCflagsKey)
	sortIncludes(expandCFlags, cfg, exts, excludeSubdirs)

	cfg.Name = NormalizePackageName(cfg.Name)

	buf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buf)
	if flag&WithTab != 0 {
		jsonEncoder.SetIndent("", "\t")
	}
	err := jsonEncoder.Encode(cfg)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
