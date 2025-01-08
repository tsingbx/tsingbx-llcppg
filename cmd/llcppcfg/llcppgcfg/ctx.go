package llcppgcfg

import (
	"path/filepath"
	"sort"
)

type DepCtx struct {
	cflagEntry    *CflagEntry
	idMap         map[int]*ObjFile
	relPathMap    map[string]int
	depsMap       map[*ObjFile][]int
	processingMap map[*ObjFile]struct{}
}

func NewDepCtx(cflagEntry *CflagEntry) *DepCtx {
	relPathMap := make(map[string]int)
	idMap := make(map[int]*ObjFile)
	for idx, objFile := range cflagEntry.ObjFiles {
		relPathMap[objFile.HFile] = idx
		idMap[idx] = objFile
	}
	return &DepCtx{cflagEntry: cflagEntry, relPathMap: relPathMap, idMap: idMap, depsMap: make(map[*ObjFile][]int), processingMap: make(map[*ObjFile]struct{})}
}

func (p *DepCtx) GetObjFileByRelPath(relPath string) (*ObjFile, int) {
	id := p.GetIdByRelPath(relPath)
	if id >= 0 {
		return p.GetObjFileById(id), id
	}
	return nil, -1
}

func (p *DepCtx) GetObjFileById(id int) *ObjFile {
	objFile, ok := p.idMap[id]
	if ok {
		return objFile
	}
	return nil
}

func (p *DepCtx) GetIdByRelPath(relPath string) int {
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
	if _, ok := p.processingMap[objFile]; ok {
		return
	}
	p.processingMap[objFile] = struct{}{}
	if p.depsMap[objFile] == nil {
		p.depsMap[objFile] = make([]int, 0, len(p.idMap))
	}
	for _, dep := range objFile.Deps {
		relPath, _ := filepath.Rel(p.GetInclude(), dep)
		depObjFile, id := p.GetObjFileByRelPath(relPath)
		if depObjFile != nil && id >= 0 {
			p.depsMap[objFile] = append(p.depsMap[objFile], id)
			p.ExpandDeps(depObjFile)
			p.depsMap[objFile] = append(p.depsMap[objFile], p.depsMap[depObjFile]...)
		}
	}
}

func removeDupObjId(s []int) []int {
	if len(s) < 1 {
		return s
	}
	sort.Ints([]int(s))
	prev := 1
	for curr := 1; curr < len(s); curr++ {
		if s[curr-1] != s[curr] {
			s[prev] = s[curr]
			prev++
		}
	}
	return s[:prev]
}

func removeDupFilePath(s []string) []string {
	m := make(map[string]struct{})
	r := make([]string, 0)
	for _, ss := range s {
		_, ok := m[ss]
		if !ok {
			m[ss] = struct{}{}
			r = append(r, ss)
		}
	}
	return r
}
