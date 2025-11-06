package error

import "fmt"

type DomainErrorArg struct {
	Key   string
	Value any
}

type DomainError struct {
	args    []DomainErrorArg
	message string
	code    int

	parent error
}

func (e *DomainError) Error() string {
	errorf := "domain error - code: %d; message: %s"
	if len(e.args) > 0 {
		return fmt.Sprintf(errorf+"; args: %v", e.code, e.message, e.args)
	}

	return fmt.Sprintf(errorf, e.code, e.message)
}

func (e *DomainError) SetCode(code int) *DomainError {
	e.code = code

	return e
}

func (e *DomainError) SetArgs(args ...DomainErrorArg) *DomainError {
	e.args = args

	return e
}

func (e *DomainError) SetMessage(message string) *DomainError {
	e.message = message

	return e
}

func (e *DomainError) Code() int {
	return e.code
}

func (e *DomainError) Message() string {
	return e.message
}

func (e *DomainError) Wrap(message string) *DomainError {
	err := &DomainError{
		code:    e.code,
		message: fmt.Sprintf("%s: %s", e.message, message),
		args:    e.args,

		parent: e,
	}

	return err
}

func (e *DomainError) Unwrap() error {
	return e.parent
}

func (e *DomainError) WrapErr(err error) *DomainError {
	return e.Wrap(err.Error())
}

func (e *DomainError) Args() []DomainErrorArg {
	return e.args
}

func New(message string, args ...DomainErrorArg) *DomainError {
	return &DomainError{
		args:    args,
		message: message,
	}
}

func Arg(name string, val any) DomainErrorArg {
	return DomainErrorArg{
		Key:   name,
		Value: val,
	}
}
