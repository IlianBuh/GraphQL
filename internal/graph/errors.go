package graph

import (
	"errors"

	"github.com/IlianBuh/GraphQL/internal/clients/sso"
	ssocodes "github.com/IlianBuh/GraphQL/internal/clients/sso/sso-codes"
)

const (
	InvalidArgument = "invalid arguments"
	Internal        = "internal"
)

func handleSsoError(err *sso.Error) error {
	switch err.Code {
	case ssocodes.InvalidArgument:
		return sendErr(InvalidArgument, err)
	case ssocodes.Unknown, ssocodes.Internal:
		return sendErr(Internal, err)
	default:
		return sendErr(Internal, err)
	}
}

func sendErr(mess string, err error) error {
	return errors.New(mess + ": " + err.Error())
}
