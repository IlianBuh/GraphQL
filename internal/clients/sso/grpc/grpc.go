package grpc

import (
	"context"
	"log/slog"

	"github.com/IlianBuh/GraphQL/internal/clients/sso"
	ssocodes "github.com/IlianBuh/GraphQL/internal/clients/sso/sso-codes"
	"github.com/IlianBuh/GraphQL/internal/domain/models"
	"github.com/IlianBuh/GraphQL/internal/lib/net"
	"github.com/IlianBuh/GraphQL/internal/lib/sl"
	e "github.com/IlianBuh/GraphQL/pkg/errors"
	authv1 "github.com/IlianBuh/SSO_Protobuf/gen/go/auth"
	userinfov1 "github.com/IlianBuh/SSO_Protobuf/gen/go/userinfo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	log        *slog.Logger
	auth       authv1.AuthClient
	userinfo   userinfov1.UserInfoClient
	connection *grpc.ClientConn
}

func NewClient(
	log *slog.Logger,
	port int,
	host string,
) (*Client, error) {
	const op = "sso-grpc-client.New"
	cc, err := grpc.NewClient(net.Join(host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	auth := authv1.NewAuthClient(cc)
	userinfo := userinfov1.NewUserInfoClient(cc)

	if err != nil {
		return nil, e.Fail(op, err)
	}

	return &Client{
		log:        log,
		auth:       auth,
		userinfo:   userinfo,
		connection: cc,
	}, nil
}

func (c *Client) SignUp(
	ctx context.Context,
	login string,
	email string,
	password string,
) (token string, _ error) {
	const op = "sso-grpc-client.SignUp"
	log := c.log.With("op", op)
	log.Info("starting to sign up user",
		slog.String("login", login),
	)

	resp, err := c.auth.SignUp(
		ctx,
		&authv1.SignUpRequest{
			Login:    login,
			Email:    email,
			Password: password,
		},
	)
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			log.Warn("invalid arguments", sl.Err(err))
			return "",
				sso.NewError(e.Fail(op, err), ssocodes.InvalidArgument)
		case codes.Internal:
			log.Error("internal error", sl.Err(err))
			return "",
				sso.NewError(e.Fail(op, err), ssocodes.Internal)
		default:
			log.Error("unknown error was received from sso", sl.Err(err))
			return "",
				sso.NewError(e.Fail(op, err), ssocodes.Unknown)
		}
	}

	log.Info("user is signed up")
	return resp.GetToken(), nil
}

func (c *Client) LogIn(
	ctx context.Context,
	login string,
	password string,
) (token string, _ error) {

	const op = "sso-grpc-client.Login"
	log := c.log.With("op", op)
	log.Info("starting to log in user",
		slog.String("login", login),
	)

	resp, err := c.auth.Login(
		ctx,
		&authv1.LoginRequest{
			Login:    login,
			Password: password,
		},
	)
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument:
			log.Warn("invalid arguments", sl.Err(err))
			return "",
				sso.NewError(e.Fail(op, err), ssocodes.InvalidArgument)
		case codes.Internal:
			log.Error("internal error", sl.Err(err))
			return "",
				sso.NewError(e.Fail(op, err), ssocodes.Internal)
		default:
			log.Error("unknown error was received from sso", sl.Err(err))
			return "",
				sso.NewError(e.Fail(op, err), ssocodes.Unknown)
		}
	}

	log.Info("user is logged in")
	return resp.GetToken(), nil
}

func (c *Client) FollowersList(
	ctx context.Context,
	userID int32,
) ([]models.User, error) {
	panic("implement me")
}

func (c *Client) Stop() error {
	const op = "op"
	log := c.log.With(slog.String("op", op))
	log.Info("stopping sso-grpc client")

	err := c.connection.Close()
	if err != nil {
		log.Error("failed to close connections", sl.Err(err))
		return e.Fail(op, err)
	}

	return nil
}
