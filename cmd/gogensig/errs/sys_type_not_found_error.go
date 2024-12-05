package errs

import "fmt"

//fmt.Errorf("sys type %s in %s not found in package %s, full path %s", name, info.IncPath, pkg, info.Path)

type SysTypeNotFoundError struct {
	sysType     string
	sysTypeFile string
	pkg         string
	fullPath    string
}

func (p *SysTypeNotFoundError) Error() string {
	return fmt.Sprintf("sys type %s in %s not found in package %s, full path %s",
		p.sysType, p.sysTypeFile, p.pkg, p.fullPath)
}

func NewSysTypeNotFoundError(sysType, sysTypeFile, pkg, fullPath string) *SysTypeNotFoundError {
	return &SysTypeNotFoundError{sysType: sysType, sysTypeFile: sysTypeFile, pkg: pkg, fullPath: fullPath}
}
