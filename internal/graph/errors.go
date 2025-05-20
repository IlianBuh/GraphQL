package graph

import (
	"encoding/json"
	"errors"

	ssocodes "github.com/IlianBuh/GraphQL/internal/clients/sso/codes"
	serrors "github.com/IlianBuh/GraphQL/internal/clients/sso/errors"
)

var (
	CodeInvalidArgument = 1
	InvalidArgument     = GQLErrorInfo{
		Code:    1,
		Message: "invalid arguments",
	}
	Internal = GQLErrorInfo{
		Code:    2,
		Message: "internal",
	}
)

func handleError(err error) error {
	var ssoerr *serrors.Error
	if errors.As(err, &ssoerr) {
		switch ssoerr.Code {
		case ssocodes.InvalidArgument:
			gqlerror := GQLError{
				Base: err,
				Info: InvalidArgument,
			}
			return gqlerror
		case ssocodes.Unknown, ssocodes.Internal:
			gqlerror := GQLError{
				Base: err,
				Info: Internal,
			}
			return gqlerror
		default:
			gqlerror := GQLError{
				Base: err,
				Info: Internal,
			}
			return gqlerror
		}
	} else {
		gqlerror := GQLError{
			Base: err,
			Info: Internal,
		}
		return gqlerror
	}
}

func sendErr(info GQLErrorInfo, err error) error {
	return GQLError{
		Base: err,
		Info: info,
	}
}

type GQLError struct {
	Base error `json:"-"`
	Info GQLErrorInfo
}

type GQLErrorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (g GQLError) Error() string {
	bytes, _ := json.Marshal(g)
	return string(bytes)
}

func (g GQLError) Unwrap() error {
	return g.Base
}
