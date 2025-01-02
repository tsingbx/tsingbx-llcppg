/*
 * Copyright (c) 2024 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package types

import "fmt"

type ObjFile struct {
	OFile string   `json:"ofile"`
	HFile string   `json:"hfile"`
	Deps  []string `json:"deps"`
}

func (o *ObjFile) String() string {
	return fmt.Sprintf("{OFile:%s, HFile:%s, Deps:%v}", o.OFile, o.HFile, o.Deps)
}

type CflagEntry struct {
	Include  string    `json:"include"`
	ObjFiles []ObjFile `json:"objfiles"`
}

func (c *CflagEntry) String() string {
	return fmt.Sprintf("{Include:%s, ObjFiles:%v}", c.Include, c.ObjFiles)
}

// Config represents a configuration for the llcppg tool.
type Config struct {
	Name         string       `json:"name"`
	CFlags       string       `json:"cflags"`
	CflagEntrys  []CflagEntry `json:"cflagEntrys"`
	Libs         string       `json:"libs"`
	Include      []string     `json:"include"`
	Deps         []string     `json:"deps"`
	TrimPrefixes []string     `json:"trimPrefixes"`
	Cplusplus    bool         `json:"cplusplus"`
}

type SymbolInfo struct {
	Mangle string `json:"mangle"` // C++ Symbol
	CPP    string `json:"c++"`    // C++ function name
	Go     string `json:"go"`     // Go function name
}
