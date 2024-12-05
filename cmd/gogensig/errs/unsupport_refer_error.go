package errs

import "fmt"

type UnsupportedReferError struct {
	typ any
}

func (p *UnsupportedReferError) Error() string {
	return fmt.Sprintf("unsupported refer: %T", p.typ)
}

func NewUnsupportedReferError(typ any) *UnsupportedReferError {
	return &UnsupportedReferError{typ: typ}
}
