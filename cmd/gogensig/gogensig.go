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

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/goplus/llcppg/_xtool/llcppsymg/args"
	"github.com/goplus/llcppg/cmd/gogensig/config"
	"github.com/goplus/llcppg/cmd/gogensig/convert"
	"github.com/goplus/llcppg/cmd/gogensig/convert/basic"
	"github.com/goplus/llcppg/cmd/gogensig/dbg"
	"github.com/goplus/llcppg/cmd/gogensig/unmarshal"
)

func main() {
	ags, remainArgs := args.ParseArgs(os.Args[1:], args.LLCPPG_SIGFETCH, nil)

	if ags.Help {
		printUsage()
		return
	}

	if ags.Verbose {
		dbg.SetDebugAll()
	}

	var data []byte
	var err error
	if ags.UseStdin {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(ags.CfgFile)
	}
	check(err)

	var cfg string
	for i := 0; i < len(remainArgs); i++ {
		arg := remainArgs[i]
		if strings.HasPrefix(arg, "-cfg=") {
			cfg = args.StringArg(arg, args.LLCPPG_CFG)
		}
	}
	if cfg == "" {
		cfg = args.LLCPPG_CFG
	}

	conf, err := config.GetCppgCfgFromPath(cfg)
	check(err)

	wd, err := os.Getwd()
	check(err)

	err = runGoCmds(wd, conf.Name)
	fmt.Println(err)

	p, _, err := basic.ConvertProcesser(&basic.Config{
		AstConvertConfig: convert.AstConvertConfig{
			PkgName:  conf.Name,
			CfgFile:  filepath.Join(wd, cfg),
			SymbFile: filepath.Join(wd, "llcppg.symb.json"),
			PubFile:  filepath.Join(wd, "llcppg.pub"),
		},
	})
	check(err)

	inputdata, err := unmarshal.UnmarshalFileSet(data)
	check(err)

	err = p.ProcessFileSet(inputdata)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func runGoCmds(wd, pkg string) error {
	dir := filepath.Join(wd, pkg)

	err := os.MkdirAll(dir, 0744)
	if err != nil {
		return err
	}

	err = os.Chdir(pkg)
	if err != nil {
		return err
	}

	err = config.RunCommand(dir, "go", "mod", "init", pkg)
	if err != nil {
		return err
	}

	return config.RunCommand(dir, "go", "get", "github.com/goplus/llgo@main")
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: gogensig [-v] [sigfetch-file]")
}
