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

	"github.com/goplus/llcppg/_xtool/llcppsigfetch/dbg"
	"github.com/goplus/llcppg/_xtool/llcppsigfetch/parse"
	"github.com/goplus/llcppg/_xtool/llcppsymg/args"
	"github.com/goplus/llcppg/_xtool/llcppsymg/clangutils"
	"github.com/goplus/llcppg/_xtool/llcppsymg/config"
	"github.com/goplus/llcppg/llcppg"
	"github.com/goplus/llgo/c"
	"github.com/goplus/llgo/c/cjson"
)

func main() {
	ags, remainArgs := args.ParseArgs(os.Args[1:], args.LLCPPG_CFG, map[string]bool{
		"--extract": true,
	})

	if ags.Help {
		printUsage()
		return
	}
	if ags.VerboseSigfetchParse {
		dbg.SetDebugParse()
	}
	if ags.Verbose {
		dbg.SetDebugAll()
	}
	extract := false
	out := false

	var extractFile string
	isTemp := false
	isCpp := true
	parseAllComments := true
	otherArgs := []string{}

	for i := 0; i < len(remainArgs); i++ {
		arg := remainArgs[i]
		switch {
		case arg == "--extract":
			extract = true
			if i+1 < len(remainArgs) && !strings.HasPrefix(remainArgs[i+1], "-") {
				extractFile = remainArgs[i+1]
				i++
			} else {
				fmt.Fprintln(os.Stderr, "Error: --extract requires a valid file argument")
				printUsage()
				os.Exit(1)
			}
		case strings.HasPrefix(arg, "-out="):
			out = args.BoolArg(arg, false)
		case strings.HasPrefix(arg, "-temp="):
			isTemp = args.BoolArg(arg, false)
		case strings.HasPrefix(arg, "-cpp="):
			isCpp = args.BoolArg(arg, true)
		case strings.HasPrefix(arg, "-parse-all-comment"):
			parseAllComments = args.BoolArg(arg, true)
		case strings.HasPrefix(arg, "-ClangResourceDir="):
			// temp to avoid call clang  in llcppsigfetch,will cause hang
			parse.ClangResourceDir = args.StringArg(arg, "")
		case strings.HasPrefix(arg, "-ClangSearchPath="):
			// temp to avoid call clang  in llcppsigfetch,will cause hang
			parse.ClangSearchPath = strings.Split(args.StringArg(arg, ""), ",")
		default:
			otherArgs = append(otherArgs, arg)
		}
	}

	if extract {
		if ags.Verbose {
			fmt.Fprintln(os.Stderr, "runExtract: extractFile:", extractFile)
			fmt.Fprintln(os.Stderr, "isTemp:", isTemp)
			fmt.Fprintln(os.Stderr, "isCpp:", isCpp)
			fmt.Fprintln(os.Stderr, "parseAllComment:", parseAllComments)
			fmt.Fprintln(os.Stderr, "out:", out)
			fmt.Fprintln(os.Stderr, "otherArgs:", otherArgs)
		}
		runExtract(extractFile, isTemp, isCpp, out, otherArgs, ags.Verbose, parseAllComments)
	} else {
		if ags.Verbose {
			fmt.Fprintln(os.Stderr, "runFromConfig: config file:", ags.CfgFile)
			fmt.Fprintln(os.Stderr, "use stdin:", ags.UseStdin)
			fmt.Fprintln(os.Stderr, "output to file:", out)
		}
		runFromConfig(ags.CfgFile, ags.UseStdin, out, ags.Verbose, ags.ParseAllComment)
	}

}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  llcppsigfetch [-v] [-out=<bool>] [config_file]")
	fmt.Println("  OR")
	fmt.Println("  llcppsigfetch --extract <file> [-out=<bool>] [-temp=<bool>] [-cpp=<bool>] [-v] [args...]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  [config_file]: Path to the configuration file (use '-' for stdin)")
	fmt.Println("                   If not provided, uses default 'llcppg.cfg'")
	fmt.Println("  -out=<bool>:     Optional. Set to 'true' to output results to a file,")
	fmt.Println("                   'false' (default) to output to stdout")
	fmt.Println("                   This option can be used with both modes")
	fmt.Println("")
	fmt.Println("  --extract:       Extract information from a single file")
	fmt.Println("    <file>:        Path to the file to process, or file content if -temp=true")
	fmt.Println("    -temp=<bool>:  Optional. Set to 'true' if <file> contains file content,")
	fmt.Println("                   'false' (default) if it's a file path")
	fmt.Println("    -cpp=<bool>:   Optional. Set to 'true' if the language is C++ (default: true)")
	fmt.Println("                   If not present, <file> is a file path")
	fmt.Println("    [args]:        Optional additional arguments")
	fmt.Println("                   Default for C++: -x c++")
	fmt.Println("                   Default for C: -x c")
	fmt.Println("")
	fmt.Println("  --help, -h:      Show this help message")
	fmt.Println("")
	fmt.Println("Note: The two usage modes are mutually exclusive. Use either [<config_file>] OR --extract, not both.")
}

func runFromConfig(cfgFile string, useStdin bool, outputToFile bool, verbose bool, parseAllComment bool) {
	var data []byte
	var err error
	if useStdin {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(cfgFile)
	}
	if verbose {
		if useStdin {
			fmt.Fprintln(os.Stderr, "runFromConfig: read from stdin")
		} else {
			fmt.Fprintln(os.Stderr, "runFromConfig: read from file", cfgFile)
		}
	}
	check(err)

	conf, err := config.GetConf(data)
	check(err)
	defer conf.Delete()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse config file:", cfgFile)
		os.Exit(1)
	}

	converter, err := parse.Do(&parse.ParseConfig{
		Conf: conf.Config,
	}, parseAllComment)

	check(err)
	info := converter.Output()
	str := info.Print()
	defer cjson.FreeCStr(str)
	defer info.Delete()
	outputResult(str, outputToFile)
}

func runExtract(content string, isTemp bool, isCpp bool, outToFile bool, otherArgs []string, verbose bool, parseAllComment bool) {
	var file string
	cflags := otherArgs
	if isTemp {
		temp, err := os.Create(clangutils.TEMP_FILE)
		if err != nil {
			panic(err)
		}
		defer temp.Close()
		defer os.Remove(file)
		temp.Write([]byte(content))
		file = temp.Name()
		cflags = append(cflags, "-I"+filepath.Dir(file))
	} else {
		file = content
	}

	converter, err := parse.Do(&parse.ParseConfig{
		Conf: &llcppg.Config{
			Include: []string{file},
			CFlags:  strings.Join(cflags, ""),
		},
	}, parseAllComment)
	check(err)

	_, err = converter.Convert()
	check(err)

	result := converter.Output()
	cstr := result.Print()
	outputResult(cstr, outToFile)
	cjson.FreeCStr(cstr)
	result.Delete()
	converter.Dispose()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func outputResult(result *c.Char, outputToFile bool) {
	if outputToFile {
		outputFile := "llcppg.sigfetch.json"
		err := os.WriteFile(outputFile, []byte(c.GoString(result)), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Results saved to %s\n", outputFile)
	} else {
		c.Printf(c.Str("%s"), result)
	}
}
