package sso

import (
	"context"
	"log/slog"

	"github.com/IlianBuh/GraphQL/internal/clients/sso/grpc/auth"
	"github.com/IlianBuh/GraphQL/internal/clients/sso/grpc/follow"
	"github.com/IlianBuh/GraphQL/internal/clients/sso/grpc/userinfo"
	"github.com/IlianBuh/GraphQL/internal/domain/models"
	e "github.com/IlianBuh/GraphQL/internal/lib/errors"
	"github.com/IlianBuh/GraphQL/internal/lib/net"
	"github.com/IlianBuh/GraphQL/internal/lib/sl"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient interface {
	SignUp(
		ctx context.Context,
		login string,
		email string,
		password string,
	) (token string, _ error)
	LogIn(
		ctx context.Context,
		login string,
		password string,
	) (token string, _ error)
}

type FollowClient interface {
	FollowersList(
		ctx context.Context,
		userID int32,
	) ([]*models.User, error)
	FolloweesList(
		ctx context.Context,
		userID int32,
	) ([]*models.User, error)
	Follow(
		ctx context.Context,
		srcId int,
		targetId int,
	) error
	Unfollow(
		ctx context.Context,
		srcId int,
		targetId int,
	) error
}

type UserInfoClient interface {
	User(ctx context.Context, uuid int) (*models.User, error)
	Users(ctx context.Context, uuid []int) ([]*models.User, error)
	UsersExist(ctx context.Context, uuid []int) (bool, error)
	UsersByLogin(ctx context.Context, login string) ([]*models.User, error)
}

type SSOClient struct {
	log        *slog.Logger
	auth       AuthClient
	follow     FollowClient
	userinfo   UserInfoClient
	connection *grpc.ClientConn
}

func NewClient(
	log *slog.Logger,
	port int,
	host string,
) (*SSOClient, error) {
	const op = "sso-grpc-client.New"
	cc, err := grpc.NewClient(net.Join(host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, e.Fail(op, err)
	}

	auth := auth.NewClient(log, cc)
	follow := follow.NewClient(log, cc)
	userinfo := userinfo.NewClient(log, cc)

	return &SSOClient{
		log:        log,
		auth:       auth,
		follow:     follow,
		userinfo:   userinfo,
		connection: cc,
	}, nil
}

func (s *SSOClient) SignUp(
	ctx context.Context,
	login string,
	email string,
	password string,
) (token string, _ error) {
	return s.auth.SignUp(ctx, login, email, password)
}

func (s *SSOClient) LogIn(
	ctx context.Context,
	login string,
	password string,
) (token string, _ error) {
	return s.auth.LogIn(ctx, login, password)
}

func (s *SSOClient) FollowersList(
	ctx context.Context,
	userID int32,
) ([]*models.User, error) {
	return s.follow.FollowersList(ctx, userID)
}

func (s *SSOClient) FolloweesList(
	ctx context.Context,
	userID int32,
) ([]*models.User, error) {
	return s.follow.FolloweesList(ctx, userID)
}

func (s *SSOClient) Follow(
	ctx context.Context,
	srcId int,
	targetId int,
) error {
	return s.follow.Follow(ctx, srcId, targetId)
}

func (s *SSOClient) Unfollow(
	ctx context.Context,
	srcId int,
	targetId int,
) error {
	return s.follow.Unfollow(ctx, srcId, targetId)
}

func (s *SSOClient) Stop() error {
	const op = "sso-client.Stop"
	log := s.log.With(slog.String("op", op))
	log.Info("stopping sso-client", slog.String("op", op))

	err := s.connection.Close()
	if err != nil {
		log.Error("failed to close connection", sl.Err(err))
		return e.Fail(op, err)
	}

	log.Info("sso-clietn stopped")
	return nil
}

func (s *SSOClient) User(ctx context.Context, uuid int) (*models.User, error) {
	return s.userinfo.User(ctx, uuid)
}

func (s *SSOClient) Users(ctx context.Context, uuid []int) ([]*models.User, error) {
	return s.userinfo.Users(ctx, uuid)
}

func (s *SSOClient) UsersExist(ctx context.Context, uuid []int) (bool, error) {
	return s.userinfo.UsersExist(ctx, uuid)
}

func (s *SSOClient) UsersByLogin(ctx context.Context, login string) ([]*models.User, error) {
	return s.userinfo.UsersByLogin(ctx, login)
}
