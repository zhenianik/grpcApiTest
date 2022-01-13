package model

import (
	"github.com/zhenianik/grpcApiTest/internal/controller/grpc/api"
)

type UserID int64

type User struct {
	Id    UserID
	Name  string
	Email string
}

func (u *User) Decode(user *api.AddUserReq) {
	u.Name = user.Name
	u.Email = user.Email
}

func (u *User) Encode() *api.User {
	return &api.User{
		Id:    int64(u.Id),
		Name:  u.Name,
		Email: u.Email,
	}
}
