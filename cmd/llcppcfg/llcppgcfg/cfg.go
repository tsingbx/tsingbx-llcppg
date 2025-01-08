package llcppgcfg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"unicode"

	"github.com/goplus/llcppg/cmdout"
	"github.com/goplus/llcppg/types"
)

type RunMode int

type FlagMode int

const (
	NormalMode RunMode = iota
	ExpandMode
)

const (
	WithSort FlagMode = 1 << iota
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

func isExcludeDir(relPath string, excludeSubdirs []string) bool {
	if len(excludeSubdirs) == 0 {
		return false
	}
	index := strings.IndexRune(relPath, filepath.Separator)
	if index >= 0 {
		dir := relPath[:index]
		for _, subdir := range excludeSubdirs {
			if subdir == dir {
				return true
			}
		}
	}
	return false
}

func doExpandCflags(str string, excludeSubdirs []string, fn func(s string) bool) ([]string, string) {
	list := strings.Fields(str)
	contains := make(map[string]string, 0)
	for _, l := range list {
		trimStr := strings.TrimPrefix(l, "-I")
		trimStr += string(filepath.Separator)
		err := filepath.WalkDir(trimStr, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if !fn(d.Name()) {
				return nil
			}
			_, ok := contains[path]
			if !ok {
				relPath, errRel := filepath.Rel(trimStr, path)
				if errRel != nil {
					return errRel
				}
				if isExcludeDir(relPath, excludeSubdirs) {
					return nil
				}
				contains[path] = relPath
			}
			return nil
		})
		if err != nil {
			log.Println(err)
		}
	}

	includes := make([]string, 0)
	includeMap := make(map[string]struct{})
	for path, relPath := range contains {
		includeDir, found := strings.CutSuffix(path, relPath)
		if found {
			includeMap[includeDir] = struct{}{}
		}
		includes = append(includes, relPath)
	}
	var flagsBuilder strings.Builder
	for include := range includeMap {
		if flagsBuilder.Len() > 0 {
			flagsBuilder.WriteRune(' ')
		}
		flagsBuilder.WriteString("-I" + include)
	}
	flags := flagsBuilder.String()
	return includes, flags
}

func ExpandName(name string, dir string, libsOrCflags string) (expand string, org string) {
	originString := fmt.Sprintf("$(pkg-config --%s %s)", libsOrCflags, name)
	return cmdout.ExpandString(originString, dir)
}

func ExpandLibsName(name string, dir string) (expand string, org string) {
	return ExpandName(name, dir, "libs")
}

func ExpandCflags(originCFlags string, exts []string, excludeDirs []string) (includes []string, expand string, org string) {
	cflags, orgCflags := cmdout.ExpandString(originCFlags, "")
	expandIncludes, expandCflags := doExpandCflags(cflags, excludeDirs, func(s string) bool {
		ext := filepath.Ext(s)
		for _, e := range exts {
			if e == ext {
				return true
			}
		}
		return false
	})
	if len(expandCflags) > 0 {
		cflags = expandCflags
	}
	return expandIncludes, cflags, orgCflags
}

func ExpandCFlagsName(name string, exts []string, excludeDirs []string) (includes []string, expand string, org string) {
	originCFlags := fmt.Sprintf("$(pkg-config --cflags %s)", name)
	return ExpandCflags(originCFlags, exts, excludeDirs)
}

func expandCFlagsAndLibs(name string, cfg *LLCppConfig, dir string, exts []string, excludeDirs []string) {
	cfg.Include, cfg.CFlags, _ = ExpandCFlagsName(name, exts, excludeDirs)
	cfg.Libs, _ = ExpandLibsName(name, dir)
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

func parseFileEntry(trimStr, path string, d fs.DirEntry, exts []string, excludeSubdirs []string) *ObjFile {
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
	relPath, err := filepath.Rel(trimStr, path)
	if err != nil {
		return nil
	}
	if isExcludeDir(relPath, excludeSubdirs) {
		return nil
	}
	clangCmd := cmdout.NewExecCommand("clang", "-I"+trimStr, "-MM", relPath)
	outString, err := cmdout.GetOut(clangCmd, trimStr)
	if err != nil || outString == "" {
		objFile := NewObjFile(relPath, relPath)
		return objFile
	}
	outString = strings.ReplaceAll(outString, "\\\n", "\n")
	fields := strings.Fields(outString)
	lines, objFileStr := findDepSlice(fields)
	objFile := NewObjFileString(objFileStr, relPath)
	objFile.Deps = append(objFile.Deps, lines...)
	return objFile
}

func parseCFlagsEntry(l string, exts []string, excludeSubdirs []string) *CflagEntry {
	if !strings.HasPrefix(l, "-I") {
		return nil
	}
	trimStr := strings.TrimPrefix(l, "-I")
	trimStr += string(filepath.Separator)
	var cflagEntry CflagEntry
	cflagEntry.Include = trimStr
	cflagEntry.ObjFiles = make([]*ObjFile, 0)
	err := filepath.WalkDir(trimStr, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		pObjFile := parseFileEntry(trimStr, path, d, exts, excludeSubdirs)
		if pObjFile != nil {
			cflagEntry.ObjFiles = append(cflagEntry.ObjFiles, pObjFile)
		}
		return nil
	})
	sort.Slice(cflagEntry.ObjFiles, func(i, j int) bool {
		return len(cflagEntry.ObjFiles[i].Deps) > len(cflagEntry.ObjFiles[j].Deps)
	})
	if err != nil {
		log.Println(err)
	}
	return &cflagEntry
}

func sortIncludes(expandCflags string, cfg *LLCppConfig, exts []string, excludeSubdirs []string) {
	list := strings.Fields(expandCflags)
	cflagEntryList := make([]*CflagEntry, 0)
	for _, l := range list {
		pCflagEntry := parseCFlagsEntry(l, exts, excludeSubdirs)
		if pCflagEntry != nil {
			cflagEntryList = append(cflagEntryList, pCflagEntry)
		}
	}
	cfg.Include = make([]string, 0)
	for _, cflagEntry := range cflagEntryList {
		depCtx := NewDepCtx(cflagEntry)
		for _, objFile := range cflagEntry.ObjFiles {
			depCtx.ExpandDeps(objFile)
		}
		for _, objFile := range cflagEntry.ObjFiles {
			cfg.Include = append(cfg.Include, objFile.HFile)
			expandDepIDs := depCtx.depsMap[objFile]
			for _, depID := range expandDepIDs {
				depObj := depCtx.GetObjFileByID(depID)
				cfg.Include = append(cfg.Include, depObj.HFile)
			}
		}
	}
	cfg.Include = removeDups(cfg.Include)
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

func GenCfg(name string, flag FlagMode, expand RunMode, exts []string, excludeSubdirs []string) (*bytes.Buffer, error) {
	if len(name) == 0 {
		return nil, newEmptyStringError("name")
	}
	cfg := NewLLCppConfig(name, flag)
	switch expand {
	case ExpandMode:
		if flag&WithSort != 0 {
			cfg.CFlags, _ = ExpandName(name, "", "cflags")
			cfg.Libs, _ = ExpandName(name, "", "libs")
			sortIncludes(cfg.CFlags, cfg, exts, excludeSubdirs)
		} else {
			expandCFlagsAndLibs(name, cfg, "", exts, excludeSubdirs)
		}
	case NormalMode:
		if flag&WithSort != 0 {
			cfg.CFlags, _ = ExpandName(name, "", "cflags")
			sortIncludes(cfg.CFlags, cfg, exts, excludeSubdirs)
		} else {
			cfg.Include, cfg.CFlags, _ = ExpandCFlagsName(name, exts, excludeSubdirs)
		}
	}

	cfg.Name = NormalizePackageName(cfg.Name)

	if runtime.GOOS == LINUX {
		libpath, _ := SearchLib(name)
		if len(libpath) > 0 {
			libs, err := cmdout.GetOut(cmdout.NewExecCommand("pkg-config", "--libs", name), "")
			if err == nil {
				cfg.Libs = fmt.Sprintf("-L%s %s", libpath, strings.TrimSpace(libs))
			}
		}
	}

	buf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buf)
	jsonEncoder.SetIndent("", "\t")
	err := jsonEncoder.Encode(cfg)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
