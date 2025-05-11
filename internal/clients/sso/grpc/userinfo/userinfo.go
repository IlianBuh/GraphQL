package userinfo

import (
	_ "github.com/IlianBuh/SSO_Protobuf/gen/go/user"
	userinfov1 "github.com/IlianBuh/SSO_Protobuf/gen/go/userinfo"
	"google.golang.org/grpc"
)

type UserInfoClient struct {
	userinfo userinfov1.UserInfoClient
}

func NewClient(cc *grpc.ClientConn) *UserInfoClient {
	userinfo := userinfov1.NewUserInfoClient(cc)
	return &UserInfoClient{
		userinfo: userinfo,
	}
}
