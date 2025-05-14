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

package cl

import "github.com/goplus/llcppg/ast"

/*
import (
	"go/types"

	"github.com/goplus/llcppg/ast"
)

func toType(ctx *blockCtx, typ ast.Expr, flags int) types.Type {
	panic("todo")
	t, _ := toTypeEx(ctx, ctx.cb.Scope(), nil, typ, flags, false)
	return t
}
*/

// compileTypeDecl compiles a struct/union/class declaration.
func compileTypeDecl(ctx *blockCtx, decl *ast.TypeDecl, global, pub bool) {
}

// compileTypedef compiles a typedef declaration.
func compileTypedef(ctx *blockCtx, decl *ast.TypedefDecl, global, pub bool) {
	panic("todo")
}

// compileEnumTypeDecl compiles an enum type declaration.
func compileEnumTypeDecl(ctx *blockCtx, decl *ast.EnumTypeDecl, global bool) {
	panic("todo")
}
