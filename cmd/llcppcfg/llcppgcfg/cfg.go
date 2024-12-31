package llcppgcfg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goplus/llcppg/types"
)

type CfgMode int

const (
	NormalMode CfgMode = iota
	ExpandMode
	SortMode
)

type LLCppConfig types.Config

type NilError struct {
}

func (p *NilError) Error() string {
	return "nil error"
}

func NewNilError() *NilError {
	return &NilError{}
}

type EmptyStringError struct {
	name string
}

func (p *EmptyStringError) Error() string {
	return p.name + " can't be empty"
}

func NewEmptyStringError(name string) *EmptyStringError {
	return &EmptyStringError{name: name}
}

func CmdOutString(cmd *exec.Cmd, dir string) (string, error) {
	if cmd == nil {
		return "", NewNilError()
	}
	outBuf := bytes.NewBufferString("")
	cmd.Stdin = os.Stdin
	cmd.Stdout = outBuf
	cmd.Env = os.Environ()
	if len(dir) > 0 {
		cmd.Dir = dir
	}
	err := cmd.Run()
	if err != nil {
		return outBuf.String(), err
	}
	return outBuf.String(), nil
}

func ExecCommand(cmdStr string, args ...string) *exec.Cmd {
	cmdStr = strings.TrimSpace(cmdStr)
	return exec.Command(cmdStr, args...)
}

func ExpandString(str string, dir string) string {
	str = strings.ReplaceAll(str, "(", "{")
	str = strings.ReplaceAll(str, ")", "}")
	expandStr := os.Expand(str, func(s string) string {
		args := strings.Fields(s)
		if len(args) == 0 {
			return ""
		}
		outString, err := CmdOutString(ExecCommand(args[0], args[1:]...), dir)
		if err != nil {
			return ""
		}
		return outString
	})
	return strings.TrimSpace(expandStr)
}

func doExpandCflags(str string, fn func(s string) bool) ([]string, string) {
	list := strings.Fields(str)
	contains := make(map[string]string, 0)
	for _, l := range list {
		trimStr := strings.TrimPrefix(l, "-I")
		trimStr += "/"
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

func ExpandLibsName(name string, dir string) string {
	originLibs := fmt.Sprintf("${pkg-config --libs %s}", name)
	return ExpandString(originLibs, dir)
}

func ExpandCflags(originCFlags string) ([]string, string) {
	cflags := ExpandString(originCFlags, "")
	expandIncludes, expandCflags := doExpandCflags(cflags, func(s string) bool {
		ext := filepath.Ext(s)
		return ext == ".h" || ext == ".hpp"
	})
	if len(expandCflags) > 0 {
		cflags = expandCflags
	}
	return expandIncludes, cflags
}

func ExpandCFlagsName(name string) ([]string, string) {
	originCFlags := fmt.Sprintf("${pkg-config --cflags %s}", name)
	return ExpandCflags(originCFlags)
}

func expandCFlagsAndLibs(name string, cfg *LLCppConfig, dir string) {
	cfg.Include, cfg.CFlags = ExpandCFlagsName(name)
	cfg.Libs = ExpandLibsName(name, dir)
}

func sortIncludes(name string, cfg *LLCppConfig, dir string) error {
	cfg.Include, cfg.CFlags = ExpandCFlagsName(name)
	cfg.Libs = ExpandLibsName(name, dir)
	originCFlags := fmt.Sprintf("${pkg-config --cflags %s}", name)
	// expand cflags
	cflags := ExpandString(originCFlags, "")
	// split cflags
	list := strings.Fields(cflags)
	// list include for every cflag
	cflagEntryList := make([]types.CflagEntry, 0)
	for _, l := range list {
		trimStr := strings.TrimPrefix(l, "-I")
		trimStr += "/"
		var cflagEntry types.CflagEntry
		cflagEntry.Cflag = l
		cflagEntry.ObjFiles = make([]types.ObjFile, 0)
		err := filepath.WalkDir(trimStr, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if !strings.HasSuffix(d.Name(), ".h") {
				return nil
			}
			if !strings.HasPrefix(d.Name(), ".") {
				relPath, err := filepath.Rel(trimStr, path)
				if err != nil {
					return nil
				}
				clangCmd := ExecCommand("clang", l, "-E", "-MM", relPath)
				outString, err := CmdOutString(clangCmd, trimStr)
				if err != nil {
					return nil
				}
				var objFile types.ObjFile
				objFile.Deps = make([]string, 0)
				lines := strings.Split(outString, "\n")
				for _, line := range lines {
					slashs := strings.Split(line, "\\")
					for _, slash := range slashs {
						if len(objFile.OFile) == 0 {
							kv := strings.Split(slash, ":")
							if len(kv) == 2 {
								objFile.OFile = kv[0]
								objFile.HFile = relPath
								dep := strings.TrimSpace(kv[1])
								dep = strings.TrimPrefix(dep, relPath)
								dep = strings.TrimSpace(dep)
								if len(dep) > 0 {
									objFile.Deps = append(objFile.Deps, dep)
								}
							}
						} else {
							if len(slash) > 0 {
								slash = strings.TrimSpace(slash)
								objFile.Deps = append(objFile.Deps, slash)
							}
						}
					}
				}
				cflagEntry.ObjFiles = append(cflagEntry.ObjFiles, objFile)
			}
			return nil
		})
		if err != nil {
			log.Println(err)
		}
		cflagEntryList = append(cflagEntryList, cflagEntry)
	}
	includeMap := make(map[string]struct{})
	cfg.Include = make([]string, 0)
	for _, cflagEntry := range cflagEntryList {
		for _, objFile := range cflagEntry.ObjFiles {
			if _, ok := includeMap[objFile.HFile]; !ok {
				includeMap[objFile.HFile] = struct{}{}
				cfg.Include = append(cfg.Include, objFile.HFile)
			}
		}
	}
	cfg.CflagEntrys = cflagEntryList
	return nil
}

func NewLLCppConfig(name string, isCpp bool) *LLCppConfig {
	cfg := &LLCppConfig{
		Name: name,
	}
	cfg.CFlags = fmt.Sprintf("$(pkg-config --cflags %s)", name)
	cfg.Libs = fmt.Sprintf("$(pkg-config --libs %s)", name)
	cfg.TrimPrefixes = []string{}
	cfg.Cplusplus = isCpp
	cfg.Include, _ = ExpandCFlagsName(name)
	return cfg
}

func GenCfg(name string, cpp bool, expand CfgMode) (*bytes.Buffer, error) {
	if len(name) == 0 {
		return nil, NewEmptyStringError("name")
	}
	cfg := NewLLCppConfig(name, cpp)
	switch expand {
	case ExpandMode:
		expandCFlagsAndLibs(name, cfg, "")
	case SortMode:
		sortIncludes(name, cfg, "")
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
