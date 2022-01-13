// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
	"time"

	"github.com/zhenianik/grpcApiTest/internal/controller/grpc/api"
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
		AddUser(ctx context.Context, user *model.User) (id model.UserID, err error)
		RemoveUser(ctx context.Context, id model.UserID) error
	}

	EventSender interface {
		Send(id model.UserID, name string, time time.Time) error
	}
)
