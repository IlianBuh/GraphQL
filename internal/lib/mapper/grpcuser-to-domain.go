package mapper

import (
	"github.com/IlianBuh/GraphQL/internal/domain/models"
	userv1 "github.com/IlianBuh/SSO_Protobuf/gen/go/user"
)

func MGrpcUserToDomain(users []*userv1.User) []*models.User {
	apiUsers := make([]*models.User, len(users))

	for i, user := range users {
		apiUsers[i] = &models.User{
			Id:    int(user.Uuid),
			Login: user.Login,
			Email: user.Email,
		}
	}

	return apiUsers
}

func GrpcUserToDomain(user *userv1.User) *models.User {
	return &models.User{
		Id:    int(user.Uuid),
		Login: user.Login,
		Email: user.Email,
	}
}
