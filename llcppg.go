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
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/goplus/llcppg/_xtool/llcppsymg/args"
	"github.com/goplus/llcppg/types"
	"github.com/goplus/llgo/xtool/env"
)

var (
	verbose   bool
	vSymg     bool
	vSigfetch bool
	vGogen    bool
)

type CommandOptions struct {
	Name    string
	Args    []string
	Verbose bool
}

func command(opts CommandOptions) *exec.Cmd {
	args := opts.Args
	if opts.Verbose {
		args = append([]string{"-v"}, args...)
	}
	return exec.Command(opts.Name, args...)
}

func llcppsymg(conf []byte) error {
	cmd := command(CommandOptions{
		Name:    "llcppsymg",
		Args:    []string{"-"},
		Verbose: verbose || vSymg,
	})
	cmd.Stdin = bytes.NewReader(conf)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func llcppsigfetch(conf []byte, out *io.PipeWriter) {
	cmd := command(CommandOptions{
		Name:    "llcppsigfetch",
		Args:    []string{"-"},
		Verbose: verbose || vSigfetch,
	})
	cmd.Stdin = bytes.NewReader(conf)
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	check(err)
	out.Close()
}

func gogensig(in io.Reader, cfg string) error {
	cmd := command(CommandOptions{
		Name:    "gogensig",
		Args:    []string{"-", "-cfg=" + cfg},
		Verbose: verbose || vGogen,
	})
	cmd.Stdin = in
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	var symbGen, codeGen, help bool
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: llcppg [config-file] [-v|-vfetch|-vsymg|-vgogen] [-symbgen] [-codegen] [-h|--help]")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}
	flag.BoolVar(&verbose, "v", false, "Enable verbose output")
	flag.BoolVar(&vSigfetch, "vfetch", false, "Enable verbose of llcppsigfetch")
	flag.BoolVar(&vSymg, "vsymg", false, "Enable verbose of llcppsymg")
	flag.BoolVar(&vGogen, "vgogen", false, "Enable verbose of gogensig")
	flag.BoolVar(&symbGen, "symbgen", false, "Only use llcppsymg to generate llcppg.symb.json")
	flag.BoolVar(&codeGen, "codegen", false, "Only use (llcppsigfetch & gogensig) to generate go code binding")
	flag.BoolVar(&help, "h", false, "Display help information")
	flag.BoolVar(&help, "help", false, "Display help information")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	remainArgs := flag.Args()

	var cfgFile string
	if len(remainArgs) > 0 {
		cfgFile = remainArgs[0]
	} else {
		cfgFile = args.LLCPPG_CFG
	}

	f, err := os.Open(cfgFile)
	check(err)
	defer f.Close()

	var conf types.Config
	json.NewDecoder(f).Decode(&conf)
	conf.CFlags = env.ExpandEnv(conf.CFlags)
	conf.Libs = env.ExpandEnv(conf.Libs)

	b, err := json.MarshalIndent(&conf, "", "  ")
	check(err)

	if !codeGen {
		err = llcppsymg(b)
		check(err)
	}

	if !symbGen {
		r, w := io.Pipe()
		go llcppsigfetch(b, w)

		err = gogensig(r, cfgFile)
		check(err)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
