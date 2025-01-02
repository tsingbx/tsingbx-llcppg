package main

import (
	"gpgerror"

	"github.com/goplus/llgo/c"
)

var GPGERRSOURCEMASK = gpgerror.SourceTGPGERRSOURCEDIM - 1
var GPGERRCODEMASK = gpgerror.CodeTGPGERRCODEDIM - 1
var GPGERRSOURCESHIFT = 24

func ErrMake(source gpgerror.SourceT, code gpgerror.CodeT) gpgerror.ErrorT {
	if code == gpgerror.CodeTGPGERRNOERROR {
		return gpgerror.ErrorT(gpgerror.CodeTGPGERRNOERROR)
	}

	return gpgerror.ErrorT(((c.Int(source) & c.Int(GPGERRSOURCEMASK)) << c.Int(GPGERRSOURCESHIFT)) | (c.Int(code) & c.Int(GPGERRCODEMASK)))
}

func main() {
	gpgerror.Init()
	err := ErrMake(gpgerror.SourceTGPGERRSOURCEUSER1, gpgerror.CodeTGPGERRBADKEY)
	errStr := err.Strerror()
	source := err.Strsource()
	c.Printf(c.Str("Error: %s\n"), errStr)
	c.Printf(c.Str("Source: %s\n"), source)
}
