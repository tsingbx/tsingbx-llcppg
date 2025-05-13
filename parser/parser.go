/*
 * Copyright (c) 2025 The GoPlus Authors (goplus.org). All rights reserved.
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

package parser

import (
	"go/token"

	"github.com/goplus/llcppg/ast"
)

// Config represents the configuration for parsing C/C++ source files.
type Config struct {
}

// ParseFile parses a C/C++ source file and returns the corresponding AST.
func ParseFile(fset *token.FileSet, filename string, conf *Config) (f *ast.File, err error) {
	panic("todo")
}
