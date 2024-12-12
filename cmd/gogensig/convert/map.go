package convert

import "regexp"

type PkgMapping struct {
	Pattern string
	Package string
}

const (
	SysPkgC       = "github.com/goplus/llgo/c"
	SysPkgOs      = "github.com/goplus/llgo/c/os"
	SysPkgSetjmp  = "github.com/goplus/llgo/c/setjmp"
	SysPkgTime    = "github.com/goplus/llgo/c/time"
	SysPkgMath    = "github.com/goplus/llgo/c/math"
	SysPkgI18n    = "github.com/goplus/llgo/c/i18n"
	SysPkgComplex = "github.com/goplus/llgo/c/math/cmplx"

	// posix
	SysPkgPthread = "github.com/goplus/llgo/c/pthread"
	SysPkgUnixNet = "github.com/goplus/llgo/c/unix/net"
)

// IncPathToPkg determines the Go package for a given C include path.
//
// According to the C language specification, when including a standard library,
// such as stdio.h, certain declarations must be provided (e.g., FILE type).
// However, these types don't have to be declared in the header file itself.
// On MacOS, for example, the actual declaration exists in _stdio.h. Therefore,
// each standard library header file can be viewed as defining an interface,
// independent of its implementation.
//
// In our current requirements, the matching follows this order:
//  1. First match standard library interface headers (like stdio.h, stdint.h)
//     which define required types and functions
//  2. Then match implementation headers (like _stdio.h, sys/_types/_int8_t.h)
//     which contain the actual type definitions
//
// For example:
// - stdio.h as interface, specifies that FILE type must be provided
// - _stdio.h as implementation, provides the actual FILE definition on MacOS

var PkgMappings = []PkgMapping{
	// c std
	{Pattern: `(^|[^a-zA-Z0-9])stdint[^a-zA-Z0-9]`, Package: SysPkgC},
	{Pattern: `(^|[^a-zA-Z0-9])stddef[^a-zA-Z0-9]`, Package: SysPkgC},
	{Pattern: `(^|[^a-zA-Z0-9])stdio[^a-zA-Z0-9]`, Package: SysPkgC},
	{Pattern: `(^|[^a-zA-Z0-9])stdlib[^a-zA-Z0-9]`, Package: SysPkgC},
	{Pattern: `(^|[^a-zA-Z0-9])string[^a-zA-Z0-9]`, Package: SysPkgC},
	{Pattern: `(^|[^a-zA-Z0-9])stdbool[^a-zA-Z0-9]`, Package: SysPkgC},
	{Pattern: `(^|[^a-zA-Z0-9])stdarg[^a-zA-Z0-9]`, Package: SysPkgC},
	{Pattern: `(^|[^a-zA-Z0-9])limits[^a-zA-Z0-9]`, Package: SysPkgC},
	{Pattern: `(^|[^a-zA-Z0-9])ctype[^a-zA-Z0-9]`, Package: SysPkgC},
	{Pattern: `(^|[^a-zA-Z0-9])uchar[^a-zA-Z0-9]`, Package: SysPkgC},
	{Pattern: `(^|[^a-zA-Z0-9])wchar[^a-zA-Z0-9]`, Package: SysPkgC},
	{Pattern: `(^|[^a-zA-Z0-9])wctype[^a-zA-Z0-9]`, Package: SysPkgC},
	{Pattern: `(^|[^a-zA-Z0-9])inttypes[^a-zA-Z0-9]`, Package: SysPkgC},

	{Pattern: `(^|[^a-zA-Z0-9])signal[^a-zA-Z0-9]`, Package: SysPkgOs},
	{Pattern: `(^|[^a-zA-Z0-9])sig[a-zA-Z]*[^a-zA-Z0-9]`, Package: SysPkgOs},
	{Pattern: `(^|[^a-zA-Z0-9])assert[^a-zA-Z0-9]`, Package: SysPkgOs},
	{Pattern: `(^|[^a-zA-Z0-9])stdalign[^a-zA-Z0-9]`, Package: SysPkgOs},

	{Pattern: `(^|[^a-zA-Z0-9])setjmp[^a-zA-Z0-9]`, Package: SysPkgSetjmp},

	{Pattern: `(^|[^a-zA-Z0-9])math[^a-zA-Z0-9]`, Package: SysPkgMath},
	{Pattern: `(^|[^a-zA-Z0-9])fenv[^a-zA-Z0-9]`, Package: SysPkgMath},
	{Pattern: `(^|[^a-zA-Z0-9])complex[^a-zA-Z0-9]`, Package: SysPkgComplex},

	{Pattern: `(^|[^a-zA-Z0-9])time[^a-zA-Z0-9]`, Package: SysPkgTime},

	{Pattern: `(^|[^a-zA-Z0-9])pthread\w*`, Package: SysPkgPthread},

	{Pattern: `(^|[^a-zA-Z0-9])locale[^a-zA-Z0-9]`, Package: SysPkgI18n},

	// c posix
	{Pattern: `(^|[^a-zA-Z0-9])socket[^a-zA-Z0-9]`, Package: SysPkgUnixNet},
	{Pattern: `(^|[^a-zA-Z0-9])arpa[^a-zA-Z0-9]`, Package: SysPkgUnixNet},
	{Pattern: `(^|[^a-zA-Z0-9])netinet6?[^a-zA-Z0-9]`, Package: SysPkgUnixNet},
	{Pattern: `(^|[^a-zA-Z0-9])net[^a-zA-Z0-9]`, Package: SysPkgUnixNet},

	// impl file
	{Pattern: `_int\d+_t`, Package: SysPkgC},
	{Pattern: `_uint\d+_t`, Package: SysPkgC},
	{Pattern: `_size_t`, Package: SysPkgC},
	{Pattern: `_intptr_t`, Package: SysPkgC},
	{Pattern: `_uintptr_t`, Package: SysPkgC},
	{Pattern: `_ptrdiff_t`, Package: SysPkgC},

	{Pattern: `malloc`, Package: SysPkgC},
	{Pattern: `alloc`, Package: SysPkgC},

	{Pattern: `(^|[^a-zA-Z0-9])clock(id_t|_t)`, Package: SysPkgTime},
	{Pattern: `(^|[^a-zA-Z0-9])(i)?time\w*`, Package: SysPkgTime},
	{Pattern: `(^|[^a-zA-Z0-9])tm[^a-zA-Z0-9]`, Package: SysPkgTime},

	// before must the special type.h such as _pthread_types.h ....
	{Pattern: `\w+_t[^a-zA-Z0-9]`, Package: SysPkgC},
	{Pattern: `(^|[^a-zA-Z0-9])types[^a-zA-Z0-9]`, Package: SysPkgC},
	{Pattern: `(^|[^a-zA-Z0-9])sys[^a-zA-Z0-9]`, Package: SysPkgOs},
}

func IncPathToPkg(incPath string) (pkg string, isDefault bool) {
	for _, mapping := range PkgMappings {
		matched, err := regexp.MatchString(mapping.Pattern, incPath)
		if err != nil {
			panic(err)
		}
		if matched {
			return mapping.Package, false
		}
	}
	return SysPkgC, true
}
