package follow

import (
	"context"
	"log/slog"

	sgrpc "github.com/IlianBuh/GraphQL/internal/clients/sso/grpc"
	"github.com/IlianBuh/GraphQL/internal/domain/models"
	"github.com/IlianBuh/GraphQL/internal/lib/mapper"
	"github.com/IlianBuh/GraphQL/internal/lib/sl"
	followv1 "github.com/IlianBuh/SSO_Protobuf/gen/go/follow"
	"google.golang.org/grpc"
)

type FollowClient struct {
	log    *slog.Logger
	follow followv1.FollowClient
}

func NewClient(log *slog.Logger, cc *grpc.ClientConn) *FollowClient {
	follow := followv1.NewFollowClient(cc)

	return &FollowClient{
		log:    log,
		follow: follow,
	}
}

func (c *FollowClient) Follow(
	ctx context.Context,
	srcId int,
	targetId int,
) error {
	const op = "sso-grpc-client.Follow"
	log := c.log.With("op", op)
	log.Info("starting to follow user")

	_, err := c.follow.Follow(
		ctx,
		&followv1.FollowRequest{
			Src:    int32(srcId),
			Target: int32(targetId),
		},
	)
	if err != nil {
		log.Error("failed to follow", sl.Err(err))
		return sgrpc.HandleError(op, err, log)
	}

	log.Info("user is followeds")
	return nil
}

func (c *FollowClient) Unfollow(
	ctx context.Context,
	srcId int,
	targetId int,
) error {
	const op = "sso-grpc-client.Unollow"
	log := c.log.With("op", op)
	log.Info("starting to unfollow user")

	_, err := c.follow.Unfollow(
		ctx,
		&followv1.UnfollowRequest{
			Src:    int32(srcId),
			Target: int32(targetId),
		},
	)
	if err != nil {
		log.Error("failed to unfollow", sl.Err(err))
		return sgrpc.HandleError(op, err, log)
	}

	log.Info("user is unfollowed")
	return nil
}

func (c *FollowClient) FollowersList(
	ctx context.Context,
	userID int32,
) ([]*models.User, error) {
	const op = "sso-grpc-client.FollowersList"
	log := c.log.With("op", op)
	log.Info("starting to fetch followers")

	resp, err := c.follow.Followers(
		ctx,
		&followv1.FollowersRequest{
			Uuid: userID,
		},
	)
	if err != nil {
		log.Error("failed to get followers' list", sl.Err(err))
		return nil, sgrpc.HandleError(op, err, log)
	}

	log.Info("followers list is gotten")
	return mapper.MGrpcUserToDomain(resp.GetUser()), nil
}

func (c *FollowClient) FolloweesList(
	ctx context.Context,
	userID int32,
) ([]*models.User, error) {
	const op = "sso-grpc-client.FolloweesList"
	log := c.log.With("op", op)
	log.Info("starting to fetch followees")

	resp, err := c.follow.Followees(
		ctx,
		&followv1.FolloweesRequest{
			Uuid: userID,
		},
	)
	if err != nil {
		log.Error("failed to get followees' list", sl.Err(err))
		return nil, sgrpc.HandleError(op, err, log)
	}

	log.Info("followees list is gotten")
	return mapper.MGrpcUserToDomain(resp.GetUser()), nil
}
