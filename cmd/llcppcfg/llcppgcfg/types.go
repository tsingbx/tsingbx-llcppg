package llcppgcfg

import (
	"fmt"
	"strings"
)

type ObjFile struct {
	OFile string
	HFile string
	Deps  []string
}

func NewObjFile(oFile, hFile string) *ObjFile {
	return &ObjFile{
		OFile: oFile,
		HFile: hFile,
		Deps:  make([]string, 0),
	}
}

func NewObjFileString(str string) *ObjFile {
	fields := strings.Split(str, ":")
	if len(fields) != 2 {
		return nil
	}
	return NewObjFile(fields[0], fields[1])
}

func (p *ObjFile) IsEqual(o *ObjFile) bool {
	if p.HFile != o.HFile {
		return false
	}
	if p.OFile != o.OFile {
		return false
	}
	if len(p.Deps) != len(o.Deps) {
		return false
	}
	for i := range p.Deps {
		if p.Deps[i] != o.Deps[i] {
			return false
		}
	}
	return true
}

func (p *ObjFile) String() string {
	return fmt.Sprintf("{OFile:%s, HFile:%s, Deps:%v}", p.OFile, p.HFile, p.Deps)
}

type CflagEntry struct {
	Include  string
	ObjFiles []*ObjFile
}

func (c *CflagEntry) String() string {
	return fmt.Sprintf("{Include:%s, ObjFiles:%v}", c.Include, c.ObjFiles)
}
