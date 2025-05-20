package userinfo

import (
	"context"
	"log/slog"

	sgrpc "github.com/IlianBuh/GraphQL/internal/clients/sso/grpc"
	"github.com/IlianBuh/GraphQL/internal/domain/models"
	"github.com/IlianBuh/GraphQL/internal/lib/mapper"
	userinfov1 "github.com/IlianBuh/SSO_Protobuf/gen/go/userinfo"
	"google.golang.org/grpc"
)

type UserInfoClient struct {
	log      *slog.Logger
	userinfo userinfov1.UserInfoClient
}

func NewClient(log *slog.Logger, cc *grpc.ClientConn) *UserInfoClient {
	userinfo := userinfov1.NewUserInfoClient(cc)
	return &UserInfoClient{
		log:      log,
		userinfo: userinfo,
	}
}

func (c *UserInfoClient) User(ctx context.Context, uuid int) (*models.User, error) {
	const op = "userinfo-client.User"
	log := c.log.With(slog.String("op", op))
	log.Info("starting to fetch user")

	resp, err := c.userinfo.User(ctx,
		&userinfov1.UserRequest{
			Uuid: int32(uuid),
		},
	)
	if err != nil {
		return nil, sgrpc.HandleError(op, err, log)
	}

	log.Info("fetching single users is completed")
	return mapper.GrpcUserToDomain(resp.GetUser()), nil
}

func (c *UserInfoClient) Users(ctx context.Context, uuids []int) ([]*models.User, error) {
	const op = "userinfo-client.Users"
	log := c.log.With(slog.String("op", op))
	log.Info("starting to fetch many users")

	resp, err := c.userinfo.Users(ctx,
		&userinfov1.UsersRequest{
			Uuids: mapper.NumsTToNumsE[int, int32](uuids),
		},
	)
	if err != nil {
		return nil, sgrpc.HandleError(op, err, log)
	}

	log.Info("fetching many users is completed")
	return mapper.MGrpcUserToDomain(resp.GetUsers()), nil
}

func (c *UserInfoClient) UsersExist(ctx context.Context, uuid []int) (bool, error) {
	const op = "userinfo-client.UsersExist"
	log := c.log.With(slog.String("op", op))
	log.Info("starting to check user's existing")

	isExists, err := c.userinfo.UsersExist(ctx,
		&userinfov1.UsersExistRequest{
			Uuid: mapper.NumsTToNumsE[int, int32](uuid),
		},
	)
	if err != nil {
		return false, sgrpc.HandleError(op, err, log)
	}

	log.Info("checking completed")
	return isExists.GetExist(), nil
}

func (c *UserInfoClient) UsersByLogin(ctx context.Context, login string) ([]*models.User, error) {
	const op = "userinfo-client.UsersByLogin"
	log := c.log.With(slog.String("op", op))
	log.Info("starting to fetch users by login", slog.String("login", login))

	users, err := c.userinfo.UsersByLogin(ctx, &userinfov1.UsersByLoginRequest{
		Login: login,
	})
	if err != nil {
		return nil, sgrpc.HandleError(op, err, log)
	}

	return mapper.MGrpcUserToDomain(users.GetUsers()), nil
}
