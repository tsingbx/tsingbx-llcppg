package errs

import "fmt"

type CantConvertError struct {
	from any
	to   string
}

func (p *CantConvertError) Error() string {
	return fmt.Sprintf("%v can't convert to %s", p.from, p.to)
}

func NewCantConvertError(from any, to string) *CantConvertError {
	return &CantConvertError{from: from, to: to}
}
