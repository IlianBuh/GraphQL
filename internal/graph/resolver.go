package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"context"

	"github.com/IlianBuh/GraphQL/internal/domain/models"
)

type SsoService interface {
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

	User(
		ctx context.Context,
		uuid int,
	) (*models.User, error)
	Users(
		ctx context.Context,
		uuid []int,
	) ([]*models.User, error)
	UsersExist(
		ctx context.Context,
		uuid []int,
	) (bool, error)
	UsersByLogin(
		ctx context.Context,
		login string,
	) ([]*models.User, error)
}

type Resolver struct {
	SSO SsoService
}
