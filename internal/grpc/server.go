package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/zhenianik/grpcApiTest/internal"
	"github.com/zhenianik/grpcApiTest/internal/controller/grpc/api"
	"github.com/zhenianik/grpcApiTest/pkg/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

//go:generate mockgen -source=server.go -destination=mock/server.go -package=mock

type (
	UserService interface {
		Get(ctx context.Context) ([]*internal.User, error)
		Add(ctx context.Context, user *internal.User) (internal.UserID, error)
		Remove(ctx context.Context, userID internal.UserID) error
	}
	EventSender interface {
		Send(id internal.UserID, name string, time time.Time) error
	}
)

type Service struct {
	api.UnimplementedUsersServer
	us UserService
	l  logger.Logger
}

func New(userService UserService, logger logger.Logger) *Service {
	return &Service{
		us: userService,
		l:  logger,
	}
}

func (s *Service) Get(ctx context.Context, _ *emptypb.Empty) (*api.UserList, error) {
	var err error
	userList := &api.UserList{}

	users, err := s.us.Get(ctx)
	if err != nil {
		s.l.Error(fmt.Errorf("getting user list error: %w", err).Error())
		return nil, err
	}

	for _, v := range users {
		userList.Users = append(userList.Users, v.Encode())
	}

	return userList, nil
}

func (s *Service) Add(ctx context.Context, request *api.AddRequest) (*api.AddResponse, error) {
	var user internal.User
	user.Decode(request)

	id, err := s.us.Add(ctx, &user)
	if err != nil {
		s.l.Error(fmt.Errorf("adding user into db error: %w", err).Error())
		return nil, err
	}

	return &api.AddResponse{Id: int64(id)}, nil
}

func (s *Service) Remove(ctx context.Context, request *api.RemoveRequest) (*emptypb.Empty, error) {
	err := s.us.Remove(ctx, internal.UserID(request.Id))
	if err != nil {
		s.l.Error(fmt.Errorf("removing user from db error: %w", err).Error())
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
