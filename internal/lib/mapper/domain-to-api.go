package mapper

import (
	"github.com/IlianBuh/GraphQL/internal/domain/models"
	"github.com/IlianBuh/GraphQL/internal/graph/model"
)

func MUsersToApi(users []*models.User) []*model.User {
	res := make([]*model.User, len(users))

	for i, user := range users {
		res[i] = UserToApi(user)
	}

	return res
}

func UserToApi(user *models.User) *model.User {
	return &model.User{
		ID:    int32(user.Id),
		Login: user.Login,
		Email: user.Email,
	}
}
