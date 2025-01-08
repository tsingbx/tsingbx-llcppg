package llcppgcfg

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"strings"

	"github.com/goplus/llcppg/cmdout"
)

const LINUX = "linux"

func IsLibInDir(dir string, lib string) bool {
	if runtime.GOOS != LINUX {
		return false
	}
	cmd := cmdout.NewExecCommand("find", dir, "-name", lib+".so*")
	out, err := cmdout.GetOut(cmd, "")
	if err != nil {
		return false
	}
	if out == "" {
		return false
	}
	return strings.Contains(out, lib)
}

func SearchLib(lib string) (string, error) {
	if runtime.GOOS != LINUX {
		return "", fmt.Errorf("only support linux")
	}
	ldCmd := cmdout.NewExecCommand("ld", "--verbose")
	ldRes, err := cmdout.GetOut(ldCmd, "")
	if err != nil {
		log.Fatal(err)
	}
	cmd := cmdout.NewExecCommand("grep", "SEARCH_DIR")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer stdin.Close()
		_, wrErr := io.WriteString(stdin, ldRes)
		panic(wrErr)
	}()
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	searchDirs := strings.Split(string(out), ";")
	libs := ""
	findDirs := make([]string, 0, len(searchDirs))
	for _, searchDir := range searchDirs {
		eqIndex := strings.Index(searchDir, "=")
		endQuoteIndex := strings.LastIndex(searchDir, "\"")
		if eqIndex != -1 && endQuoteIndex != -1 {
			dir := searchDir[eqIndex+1 : endQuoteIndex]
			if IsLibInDir(dir, lib) {
				findDirs = append(findDirs, dir)
			}
		}
	}
	maxDir := 0
	for _, dir := range findDirs {
		cnt := strings.Count(dir, "/")
		if cnt > maxDir {
			maxDir = cnt
			libs = dir
		}
	}
	return libs, nil
}
