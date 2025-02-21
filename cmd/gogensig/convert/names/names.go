package names

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
)

// NameMapper handles name mapping and uniqueness for Go symbols
type NameMapper struct {
	count   map[string]int    // tracks count of each public name for uniqueness
	mapping map[string]string // maps original c names to Go names,like: foo(in c) -> Foo(in go)
}

func NewNameMapper() *NameMapper {
	return &NameMapper{
		count:   make(map[string]int),
		mapping: make(map[string]string),
	}
}

// returns a unique Go name for an original name
// For every go name, it will be unique.
func (m *NameMapper) GetUniqueGoName(name string, trimPrefixes []string) (string, bool) {
	pubName := m.GetGoName(name, trimPrefixes)

	if _, exists := m.mapping[name]; exists {
		return pubName, pubName != name
	}

	count := m.count[pubName]
	m.count[pubName]++
	if count > 0 {
		pubName = fmt.Sprintf("%s__%d", pubName, count)
	}

	return pubName, pubName != name
}

// returns the Go name for an original name,if the name is already mapped,return the mapped name
func (m *NameMapper) GetGoName(name string, trimPrefixes []string) string {
	if goName, exists := m.mapping[name]; exists {
		if goName == "" {
			return name
		}
		return goName
	}
	return GoName(name, trimPrefixes)
}

func (m *NameMapper) SetMapping(originName, newName string) {
	value := ""
	if originName != newName {
		value = newName
	}
	m.mapping[originName] = value
}

func GoName(name string, trimPrefixes []string) string {
	name = removePrefixedName(name, trimPrefixes)
	return PubName(name)
}

func removePrefixedName(name string, trimPrefixes []string) string {
	if len(trimPrefixes) == 0 {
		return name
	}
	for _, prefix := range trimPrefixes {
		if strings.HasPrefix(name, prefix) {
			return strings.TrimPrefix(name, prefix)
		}
	}
	return name
}

func ToCamelCase(s string, firstPartUpper bool) string {
	parts := strings.Split(s, "_")
	result := []string{}
	for i, part := range parts {
		if i == 0 && !firstPartUpper {
			result = append(result, part)
			continue
		}
		if len(part) > 0 {
			result = append(result, strings.ToUpper(part[:1])+part[1:])
		}
	}
	return strings.Join(result, "")
}

func PubName(name string) string {
	if len(name) == 0 {
		return name
	}
	fChar := name[0]
	if fChar == '_' {
		i := 0
		for i < len(name) && name[i] == '_' {
			i++
		}
		prefix := name[:i]
		return "X" + prefix + ToCamelCase(name[i:], false)
	}
	if unicode.IsDigit(rune(fChar)) {
		return "X" + ToCamelCase(name, false)
	}
	return ToCamelCase(name, true)
}

// /path/to/foo.h -> foo.go
// /path/to/_intptr.h -> X_intptr.go
func HeaderFileToGo(incPath string) string {
	_, fileName := filepath.Split(incPath)
	ext := filepath.Ext(fileName)
	if len(ext) > 0 {
		fileName = strings.TrimSuffix(fileName, ext)
	}
	if strings.HasPrefix(fileName, "_") {
		fileName = "X" + fileName
	}
	return fileName + ".go"
}
