package sso

import ssocodes "github.com/IlianBuh/GraphQL/internal/clients/sso/sso-codes"

type Error struct {
	base    error
	message string
	Code    int
}

func NewError(err error, code int) *Error {
	return &Error{
		base:    err,
		message: ssocodes.Text(code),
		Code:    code,
	}
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Unwrap() error {
	return e.base
}
