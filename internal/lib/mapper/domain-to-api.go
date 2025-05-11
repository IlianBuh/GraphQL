package mapper

import (
	"github.com/IlianBuh/GraphQL/internal/domain/models"
	"github.com/IlianBuh/GraphQL/internal/graph/model"
)

func UserToApi(users []*models.User) []*model.User {
	res := make([]*model.User, len(users))

	for i, user := range users {
		res[i] = &model.User{
			ID:    int32(user.Id),
			Name:  user.Login,
			Email: user.Email,
		}
	}

	return res
}
