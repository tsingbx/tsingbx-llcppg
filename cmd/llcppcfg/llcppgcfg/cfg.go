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

func CmdOutString(cmd *exec.Cmd) (string, error) {
	if cmd == nil {
		return "", NewNilError()
	}
	outBuf := bytes.NewBufferString("")
	cmd.Stdin = os.Stdin
	cmd.Stdout = outBuf
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

func ExpandString(str string) string {
	str = strings.ReplaceAll(str, "(", "{")
	str = strings.ReplaceAll(str, ")", "}")
	expandStr := os.Expand(str, func(s string) string {
		args := strings.Fields(s)
		if len(args) == 0 {
			return ""
		}
		outString, err := CmdOutString(ExecCommand(args[0], args[1:]...))
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

func ExpandLibsName(name string) string {
	originLibs := fmt.Sprintf("${pkg-config --libs %s}", name)
	return ExpandString(originLibs)
}

func ExpandCflags(originCFlags string) ([]string, string) {
	cflags := ExpandString(originCFlags)
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

func expandCFlagsAndLibs(name string, cfg *LLCppConfig) {
	cfg.Include, cfg.CFlags = ExpandCFlagsName(name)
	cfg.Libs = ExpandLibsName(name)
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

func GenCfg(name string, cpp bool, expand bool) (*bytes.Buffer, error) {
	if len(name) == 0 {
		return nil, NewEmptyStringError("name")
	}
	cfg := NewLLCppConfig(name, cpp)
	if expand {
		expandCFlagsAndLibs(name, cfg)
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
