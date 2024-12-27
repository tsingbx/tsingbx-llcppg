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

	"github.com/goplus/llcppg/_xtool/llcppsymg/args"
	"github.com/goplus/llcppg/_xtool/llcppsymg/config"
	"github.com/goplus/llcppg/_xtool/llcppsymg/config/cfgparse"
	"github.com/goplus/llcppg/_xtool/llcppsymg/dbg"
	"github.com/goplus/llcppg/_xtool/llcppsymg/parse"
	"github.com/goplus/llcppg/_xtool/llcppsymg/symbol"
)

func main() {
	symbFile := "llcppg.symb.json"

	ags, _ := args.ParseArgs(os.Args[1:], args.LLCPPG_CFG, nil)

	var data []byte
	var err error
	if ags.UseStdin {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(ags.CfgFile)
	}

	check(err)
	conf, err := config.GetConf(data)
	check(err)
	defer conf.Delete()

	if ags.VerboseParseIsMethod {
		dbg.SetDebugParseIsMethod()
	}

	if ags.Verbose {
		dbg.SetDebugSymbol()
		if ags.UseStdin {
			fmt.Println("Config From Stdin")
		} else {
			fmt.Println("Config From File", ags.CfgFile)
		}
		fmt.Println("Name:", conf.Name)
		fmt.Println("CFlags:", conf.CFlags)
		fmt.Println("Libs:", conf.Libs)
		fmt.Println("Include:", conf.Include)
		fmt.Println("TrimPrefixes:", conf.TrimPrefixes)
		fmt.Println("Cplusplus:", conf.Cplusplus)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse config file:", ags.CfgFile)
	}
	symbols, err := symbol.ParseDylibSymbols(conf.Libs)
	check(err)

	cflag := cfgparse.ParseCFlags(conf.CFlags)
	filepaths, notFounds, err := cflag.GenHeaderFilePaths(conf.Include)
	check(err)

	if ags.Verbose {
		fmt.Println("header file paths", filepaths)
		if len(notFounds) > 0 {
			fmt.Println("not found header files", notFounds)
		}
	}

	parseConfig := parse.NewSymbolProcessorConfig(cflag, conf.TrimPrefixes, conf.Cplusplus)
	headerInfos, err := parse.ParseHeaderFile(filepaths, parseConfig, false)
	check(err)

	symbolData, err := symbol.GenerateAndUpdateSymbolTable(symbols, headerInfos, symbFile)
	check(err)

	err = os.WriteFile(symbFile, symbolData, 0644)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
