package errs

type AnonymousFuncNotSupportError struct {
}

func (p *AnonymousFuncNotSupportError) Error() string {
	return "anonymous function not supported"
}

func NewAnonymousFuncNotSupportError() *AnonymousFuncNotSupportError {
	return &AnonymousFuncNotSupportError{}
}

type ModNotFoundError struct {
}

func (p *ModNotFoundError) Error() string {
	return "go.mod not found"
}

func NewModNotFoundError() *ModNotFoundError {
	return &ModNotFoundError{}
}

type TypeDefinedError struct {
	Name       string
	OriginName string
}

func (p *TypeDefinedError) Error() string {
	return "type " + p.Name + " already defined,original name is " + p.OriginName
}

func NewTypeDefinedError(name, originName string) *TypeDefinedError {
	return &TypeDefinedError{
		Name:       name,
		OriginName: originName,
	}
}
