package user

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/zhenianik/grpcApiTest/internal/user/cache"
	"github.com/zhenianik/grpcApiTest/internal/user/model"
	"github.com/zhenianik/grpcApiTest/internal/user/repository"
	"github.com/zhenianik/grpcApiTest/pkg/api"
	"github.com/zhenianik/grpcApiTest/pkg/logger"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
	"time"
)

type GRPCServer struct {
	api.UnimplementedUserServer
	db     *repository.UserRepository
	logger *logger.Logger
	cache  *cache.RedisCache
}

func NewGRPCServer(db *pgxpool.Pool, logger *logger.Logger, cache *cache.RedisCache) *GRPCServer {
	return &GRPCServer{
		db:     repository.New(db),
		logger: logger,
		cache:  cache,
	}
}

func (s *GRPCServer) Get(ctx context.Context, _ *emptypb.Empty) (*api.UserList, error) {
	var resp *api.UserList
	var err error

	resp, err = s.cache.Get(ctx, "grpc_test_users")
	if len(resp.Users) != 0 {
		return resp, nil
	}

	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range users {
		resp.Users = append(resp.Users, v.Encode())
	}

	err = s.cache.Set(ctx, "grpc_test_users", resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *GRPCServer) Add(ctx context.Context, request *api.AddRequest) (*api.Response, error) {
	var user model.User
	user.Decode(request.Body.User)

	id, err := s.db.AddUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	err = s.logger.LogRegistration(id, fmt.Sprintf("%s %s", request.Body.User.Name, request.Body.User.Email), time.Now())
	if err != nil {
		return nil, err
	}

	return &api.Response{
		Message: fmt.Sprintf("User with name %s and email %s was added", request.Body.User.Name, request.Body.User.Email),
	}, nil
}

func (s *GRPCServer) Remove(ctx context.Context, request *api.RemoveRequest) (*api.Response, error) {
	err := s.db.RemoveUser(ctx, request.Body.User.Id)
	if err != nil {
		return nil, err
	}

	return &api.Response{
		Message: `User #` + strconv.Itoa(int(request.Body.User.Id)) + ` was removed.`,
	}, nil
}
