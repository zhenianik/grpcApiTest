// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/zhenianik/grpcApiTest/internal/controller/api"
	"github.com/zhenianik/grpcApiTest/internal/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	User interface {
		Add(ctx context.Context, request *api.AddRequest) (*api.Response, error)
		Get(ctx context.Context, _ *emptypb.Empty) (*api.UserList, error)
		Remove(ctx context.Context, request *api.RemoveRequest) (*api.Response, error)
	}

	UserRepo interface {
		GetUsers(ctx context.Context) (users []*model.User, err error)
		AddUser(ctx context.Context, user *model.User) (id int64, err error)
		RemoveUser(ctx context.Context, id int64) error
	}

	Cache interface {
		Set(ctx context.Context, key string, user *api.UserList) error
		Get(ctx context.Context, key string) (*api.UserList, error)
	}
)
