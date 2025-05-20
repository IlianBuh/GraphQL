package graph

import (
	"errors"

	ssocodes "github.com/IlianBuh/GraphQL/internal/clients/sso/codes"
	serrors "github.com/IlianBuh/GraphQL/internal/clients/sso/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
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
			gqlerror := gqlerror.Error{
				Err: ssoerr,
				Extensions: map[string]interface{}{
					"code":    InvalidArgument.Code,
					"message": InvalidArgument.Message,
				},
			}
			return &gqlerror
		case ssocodes.Unknown, ssocodes.Internal:
			gqlerror := gqlerror.Error{
				Err: ssoerr,
				Extensions: map[string]interface{}{
					"code":    Internal.Code,
					"message": Internal.Message,
				},
			}
			return &gqlerror
		default:
			gqlerror := gqlerror.Error{
				Err: ssoerr,
				Extensions: map[string]interface{}{
					"code":    Internal.Code,
					"message": Internal.Message,
				},
			}
			return &gqlerror
		}
	} else {
		gqlerror := gqlerror.Error{
			Err: err,
			Extensions: map[string]interface{}{
				"code":    Internal.Code,
				"message": Internal.Message,
			},
		}
		return &gqlerror
	}
}

func sendErr(info GQLErrorInfo, err error) error {
	return &gqlerror.Error{
		Err: err,
		Extensions: map[string]interface{}{
			"code":    info.Code,
			"message": info.Message,
		},
	}
}

type GQLErrorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
