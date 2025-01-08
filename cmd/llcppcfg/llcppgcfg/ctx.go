package llcppgcfg

import (
	"fmt"
	"path/filepath"
	"strings"
)

type ObjFile struct {
	parent *ObjFile
	OFile  string
	HFile  string
	Deps   []string
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

type DepCtx struct {
	cflagEntry *CflagEntry
	idMap      map[int]*ObjFile
	relPathMap map[string]int
	depsMap    map[*ObjFile][]int
}

func NewDepCtx(cflagEntry *CflagEntry) *DepCtx {
	relPathMap := make(map[string]int)
	idMap := make(map[int]*ObjFile)
	for idx, objFile := range cflagEntry.ObjFiles {
		relPathMap[objFile.HFile] = idx
		idMap[idx] = objFile
	}
	return &DepCtx{cflagEntry: cflagEntry, relPathMap: relPathMap, idMap: idMap, depsMap: make(map[*ObjFile][]int)}
}

func (p *DepCtx) GetObjFileByRelPath(relPath string) (*ObjFile, int) {
	id := p.GetIDByRelPath(relPath)
	if id >= 0 {
		return p.GetObjFileByID(id), id
	}
	return nil, -1
}

func (p *DepCtx) GetObjFileByID(id int) *ObjFile {
	objFile, ok := p.idMap[id]
	if ok {
		return objFile
	}
	return nil
}

func (p *DepCtx) GetIDByRelPath(relPath string) int {
	id, ok := p.relPathMap[relPath]
	if ok {
		return id
	}
	return -1
}

func (p *DepCtx) GetInclude() string {
	return p.cflagEntry.Include
}

func (p *DepCtx) ExpandDeps(objFile *ObjFile) {
	if p.depsMap[objFile] != nil {
		return
	}
	if p.depsMap[objFile] == nil {
		p.depsMap[objFile] = make([]int, 0, len(p.idMap))
	}
	if len(objFile.Deps) == 0 {
		return
	}
	for _, dep := range objFile.Deps {
		relPath, err := filepath.Rel(p.GetInclude(), dep)
		if err != nil {
			relPath = dep
		}
		depObjFile, id := p.GetObjFileByRelPath(relPath)
		if depObjFile != nil && id >= 0 {
			depObjFile.parent = objFile
			isParentHFile := false
			for parent := depObjFile; parent != nil; parent = parent.parent {
				if relPath == parent.HFile {
					isParentHFile = true
					break
				}
			}
			if isParentHFile {
				if depObjFile.HFile != objFile.HFile {
					p.depsMap[objFile] = append(p.depsMap[objFile], id)
				}
				continue
			}
			p.depsMap[objFile] = append(p.depsMap[objFile], id)
			p.ExpandDeps(depObjFile)
			p.depsMap[objFile] = append(p.depsMap[objFile], p.depsMap[depObjFile]...)
		}
	}
}

func removeDups[T comparable](s []T) []T {
	m := make(map[T]struct{})
	r := make([]T, 0)
	for _, ss := range s {
		_, ok := m[ss]
		if !ok {
			m[ss] = struct{}{}
			r = append(r, ss)
		}
	}
	return r
}
