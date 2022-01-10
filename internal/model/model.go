package model

import (
	"github.com/zhenianik/grpcApiTest/internal/controller/api"
)

type User struct {
	Id    int64
	Name  string
	Email string
}

func (u *User) Decode(user *api.AddUserReq) {
	u.Name = user.Name
	u.Email = user.Email
}

func (u *User) Encode() *api.User {
	return &api.User{
		Id:    u.Id,
		Name:  u.Name,
		Email: u.Email,
	}
}
