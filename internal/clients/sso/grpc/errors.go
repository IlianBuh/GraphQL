package grpc

import (
	"log/slog"

	ssocodes "github.com/IlianBuh/GraphQL/internal/clients/sso/codes"
	"github.com/IlianBuh/GraphQL/internal/clients/sso/errors"
	"github.com/IlianBuh/GraphQL/internal/lib/sl"
	e "github.com/IlianBuh/GraphQL/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleError(op string, err error, log *slog.Logger) error {
	switch status.Code(err) {
	case codes.InvalidArgument:
		log.Warn("invalid arguments", sl.Err(err))
		return errors.NewError(e.Fail(op, err), ssocodes.InvalidArgument)
	case codes.Internal:
		log.Error("internal error", sl.Err(err))
		return errors.NewError(e.Fail(op, err), ssocodes.Internal)
	default:
		log.Error("unknown error was received from sso", sl.Err(err))
		return errors.NewError(e.Fail(op, err), ssocodes.Unknown)
	}
}
