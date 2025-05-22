package clangtool

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PreprocessConfig struct {
	File    string
	IsCpp   bool
	Args    []string
	OutFile string
}

func Preprocess(cfg *PreprocessConfig) error {
	args := []string{"-E"}
	args = append(args, defaultArgs(cfg.IsCpp)...)
	args = append(args, cfg.Args...)
	args = append(args, cfg.File)
	args = append(args, "-o", cfg.OutFile)
	cmd := exec.Command("clang", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ComposeIncludes create Include list
// #include <file1.h>
// #include <file2.h>
func ComposeIncludes(files []string, outfile string) error {
	var str string
	for _, file := range files {
		str += ("#include <" + file + ">\n")
	}
	return os.WriteFile(outfile, []byte(str), 0644)
}

func GetIncludePaths(isCpp bool) []string {
	args := []string{"-E", "-v"}
	args = append(args, defaultArgs(isCpp)...)
	args = append(args, "/dev/null")
	cmd := exec.Command("clang", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	return ParseClangIncOutput(string(output))
}

func ParseClangIncOutput(output string) []string {
	var paths []string
	start := strings.Index(output, "#include <...> search starts here:")
	end := strings.Index(output, "End of search list.")
	if start == -1 || end == -1 {
		return paths
	}
	content := output[start:end]
	lines := strings.Split(content, "\n")
	for _, line := range lines[1:] {
		for _, item := range strings.Fields(line) {
			if path := strings.TrimSpace(item); filepath.IsAbs(path) {
				paths = append(paths, path)
			}
		}
	}
	return paths
}

func defaultArgs(isCpp bool) []string {
	args := []string{"-x", "c"}
	if isCpp {
		args = []string{"-x", "c++"}
	}
	return args
}
