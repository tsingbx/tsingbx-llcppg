package main

import (
	"fmt"

	"github.com/goplus/llcppg/_xtool/llcppsymg/syspath"
)

func main() {
	TestLdOutput()
	TestClangIncOutput()
}

func TestLdOutput() {
	fmt.Println("=== TestLdOutput ===")
	res := syspath.ParseLdOutput(
		`GNU ld (GNU Binutils for Ubuntu) 2.42
  Supported emulations:
   aarch64linux
   aarch64elf
   aarch64elf32
   aarch64elf32b
   aarch64elfb
   armelf
   armelfb
   aarch64linuxb
   aarch64linux32
   aarch64linux32b
   armelfb_linux_eabi
   armelf_linux_eabi
using internal linker script:
==================================================
/* Script for -z combreloc */
/* Copyright (C) 2014-2024 Free Software Foundation, Inc.
   Copying and distribution of this script, with or without modification,
   are permitted in any medium without royalty provided the copyright
   notice and this notice are preserved.  */
OUTPUT_FORMAT("elf64-littleaarch64", "elf64-bigaarch64",
              "elf64-littleaarch64")
OUTPUT_ARCH(aarch64)
ENTRY(_start)
SEARCH_DIR("=/usr/local/lib/aarch64-linux-gnu"); SEARCH_DIR("=/lib/aarch64-linux-gnu"); SEARCH_DIR("=/usr/lib/aarch64-linux-gnu"); SEARCH_DIR("=/usr/local/lib"); SEARCH_DIR("=/lib"); SEARCH_DIR("=/usr/lib"); SEARCH_DIR("=/usr/aarch64-linux-gnu/lib");
SECTIONS
{
  /* Read-only sections, merged into text segment: */
  PROVIDE (__executable_start = SEGMENT_START("text-segment", 0x400000)); . = SEGMENT_START("text-segment", 0x400000) + SIZEOF_HEADERS;
  .interp         : { *(.interp) }
  .note.gnu.build-id  : { *(.note.gnu.build-id) }
  .hash           : { *(.hash) }
  .gnu.hash       : { *(.gnu.hash) }
  .dynsym         : { *(.dynsym) }
  .dynstr         : { *(.dynstr) }
`)
	fmt.Println(res)
}

func TestClangIncOutput() {
	fmt.Println("=== TestClangIncOutput ===")
	res := syspath.ParseClangIncOutput(
		`Ubuntu clang version 18.1.3 (1ubuntu1)
Target: aarch64-unknown-linux-gnu
Thread model: posix
InstalledDir: /usr/bin
Found candidate GCC installation: /usr/bin/../lib/gcc/aarch64-linux-gnu/13
Selected GCC installation: /usr/bin/../lib/gcc/aarch64-linux-gnu/13
Candidate multilib: .;@m64
Selected multilib: .;@m64
 (in-process)
 "/usr/lib/llvm-18/bin/clang" -cc1 -triple aarch64-unknown-linux-gnu -E -disable-free -clear-ast-before-backend -disable-llvm-verifier -discard-value-names -main-file-name null -mrelocation-model pic -pic-level 2 -pic-is-pie -mframe-pointer=non-leaf -fmath-errno -ffp-contract=on -fno-rounding-math -mconstructor-aliases -funwind-tables=2 -target-cpu generic -target-feature +v8a -target-feature +fp-armv8 -target-feature +neon -target-abi aapcs -debugger-tuning=gdb -fdebug-compilation-dir=/root/llcppg -v -fcoverage-compilation-dir=/root/llcppg -resource-dir /usr/lib/llvm-18/lib/clang/18 -internal-isystem /usr/lib/llvm-18/lib/clang/18/include -internal-isystem /usr/local/include -internal-isystem /usr/bin/../lib/gcc/aarch64-linux-gnu/13/../../../../aarch64-linux-gnu/include -internal-externc-isystem /usr/include/aarch64-linux-gnu -internal-externc-isystem /include -internal-externc-isystem /usr/include -ferror-limit 19 -fno-signed-char -fgnuc-version=4.2.1 -fskip-odr-check-in-gmf -fcolor-diagnostics -target-feature +outline-atomics -target-feature -fmv -faddrsig -D__GCC_HAVE_DWARF2_CFI_ASM=1 -o - -x c /dev/null
clang -cc1 version 18.1.3 based upon LLVM 18.1.3 default target aarch64-unknown-linux-gnu
ignoring nonexistent directory "/usr/bin/../lib/gcc/aarch64-linux-gnu/13/../../../../aarch64-linux-gnu/include"
ignoring nonexistent directory "/include"
#include "..." search starts here:
#include <...> search starts here:
 /usr/lib/llvm-18/lib/clang/18/include
 /usr/local/include
 /usr/include/aarch64-linux-gnu
 /usr/include
End of search list.
# 1 "/dev/null"
# 1 "<built-in>" 1
# 1 "<built-in>" 3
# 399 "<built-in>" 3
# 1 "<command line>" 1
# 1 "<built-in>" 2
# 1 "/dev/null" 2
`)
	fmt.Println(res)
}
