package errs

type AnonymousFuncNotSupportError struct {
}

func (p *AnonymousFuncNotSupportError) Error() string {
	return "anonymous function not supported"
}

func NewAnonymousFuncNotSupportError() *AnonymousFuncNotSupportError {
	return &AnonymousFuncNotSupportError{}
}
