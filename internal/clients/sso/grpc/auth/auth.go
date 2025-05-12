package auth

import (
	"context"
	"log/slog"

	sgrpc "github.com/IlianBuh/GraphQL/internal/clients/sso/grpc"
	authv1 "github.com/IlianBuh/SSO_Protobuf/gen/go/auth"
	"google.golang.org/grpc"
)

type AuthClient struct {
	log  *slog.Logger
	auth authv1.AuthClient
}

func NewClient(log *slog.Logger, cc *grpc.ClientConn) *AuthClient {
	authclient := authv1.NewAuthClient(cc)

	return &AuthClient{
		log:  log,
		auth: authclient,
	}
}

func (c *AuthClient) SignUp(
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
		return "", sgrpc.HandleError(op, err, log)
	}

	log.Info("user is signed up")
	return resp.GetToken(), nil
}

func (c *AuthClient) LogIn(
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
		return "", sgrpc.HandleError(op, err, log)
	}

	log.Info("user is logged in")
	return resp.GetToken(), nil
}
