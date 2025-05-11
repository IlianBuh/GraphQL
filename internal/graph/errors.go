package graph

import (
	"errors"

	ssocodes "github.com/IlianBuh/GraphQL/internal/clients/sso/codes"
	serrors "github.com/IlianBuh/GraphQL/internal/clients/sso/errors"
)

const (
	InvalidArgument = "invalid arguments"
	Internal        = "internal"
)

func handleSsoError(err *serrors.Error) error {
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
