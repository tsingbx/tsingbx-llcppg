package cmdout

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

func GetOut(cmd *exec.Cmd, dir string) (string, error) {
	if cmd == nil {
		return "", newNilError()
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

func NewExecCommand(cmdStr string, args ...string) *exec.Cmd {
	cmdStr = strings.TrimSpace(cmdStr)
	return exec.Command(cmdStr, args...)
}

func ExpandString(str string, dir string) (expand string, org string) {
	org = str
	str = strings.ReplaceAll(str, "(", "{")
	str = strings.ReplaceAll(str, ")", "}")
	expandStr := os.Expand(str, func(s string) string {
		args := strings.Fields(s)
		if len(args) == 0 {
			return ""
		}
		outString, err := GetOut(NewExecCommand(args[0], args[1:]...), dir)
		if err != nil {
			return ""
		}
		return outString
	})
	expand = strings.TrimSpace(expandStr)
	return expand, org
}

type nilError struct {
}

func (p *nilError) Error() string {
	return "nil error"
}

func newNilError() *nilError {
	return &nilError{}
}
