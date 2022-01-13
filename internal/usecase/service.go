package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/zhenianik/grpcApiTest/internal/controller/grpc/api"
	"github.com/zhenianik/grpcApiTest/internal/model"
	"github.com/zhenianik/grpcApiTest/pkg/cache"
	"github.com/zhenianik/grpcApiTest/pkg/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	api.UnimplementedUserServer
	repo        UserRepo
	eventSender EventSender
	cache       cache.Cache
	logger      logger.Logger
}

func NewService(repo UserRepo, eventSender EventSender, cache cache.Cache, logger logger.Logger) *Service {
	return &Service{
		repo:        repo,
		eventSender: eventSender,
		cache:       cache,
		logger:      logger,
	}
}

func (s *Service) Get(ctx context.Context, _ *emptypb.Empty) (*api.UserList, error) {
	var resp []byte
	var err error
	userList := &api.UserList{}

	resp, err = s.cache.Get(ctx, "grpc_test_users")
	if len(resp) != 0 {
		s.logger.Debug("getting user from cache")
		err = json.Unmarshal(resp, &userList)
		if err != nil {
			return &api.UserList{}, err
		}
		return userList, nil
	}

	users, err := s.repo.GetUsers(ctx)
	if err != nil {
		s.logger.Error(fmt.Errorf("getting user list from db error: %w", err).Error())
		return nil, err
	}

	for _, v := range users {
		userList.Users = append(userList.Users, v.Encode())
	}

	js, err := json.Marshal(userList)
	if err != nil {
		s.logger.Error(fmt.Errorf("marshal user list error: %w", err).Error())
		return nil, err
	}

	err = s.cache.Set(ctx, "grpc_test_users", js)
	if err != nil {
		s.logger.Error(fmt.Errorf("setting cache into db error: %w", err).Error())
		return nil, err
	}

	return userList, nil
}

func (s *Service) Add(ctx context.Context, request *api.AddRequest) (*api.Response, error) {
	var user model.User
	user.Decode(request.Body.User)

	id, err := s.repo.AddUser(ctx, &user)
	if err != nil {
		s.logger.Error(fmt.Errorf("adding user into db error: %w", err).Error())
		return nil, err
	}

	err = s.eventSender.Send(id, fmt.Sprintf("%s %s", request.Body.User.Name, request.Body.User.Email), time.Now())
	if err != nil {
		s.logger.Error(fmt.Errorf("adding eventsender error: %w", err).Error())
		return nil, err
	}

	return &api.Response{
		Message: fmt.Sprintf("User with name %s and email %s was added", request.Body.User.Name, request.Body.User.Email),
	}, nil
}

func (s *Service) Remove(ctx context.Context, request *api.RemoveRequest) (*api.Response, error) {
	err := s.repo.RemoveUser(ctx, model.UserID(request.Body.User.Id))
	if err != nil {
		s.logger.Error(fmt.Errorf("removing user from db error: %w", err).Error())
		return nil, err
	}

	return &api.Response{
		Message: `User #` + strconv.Itoa(int(request.Body.User.Id)) + ` was removed.`,
	}, nil
}
