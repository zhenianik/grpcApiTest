package server

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/zhenianik/grpcApiTest/internal/db/model"
	"github.com/zhenianik/grpcApiTest/internal/db/repository"
	"github.com/zhenianik/grpcApiTest/internal/server/cache"
	"github.com/zhenianik/grpcApiTest/pkg/api"
	"github.com/zhenianik/grpcApiTest/pkg/dbLogger"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
	"time"
)

type GRPCServer struct {
	api.UnimplementedUserServer
	db       *repository.UserRepository
	dbLogger *dbLogger.Logger
	cache    *cache.RedisCache
	Logger   *logrus.Logger
}

func NewGRPCServer(db *pgxpool.Pool, dbLogger *dbLogger.Logger, cache *cache.RedisCache, logger *logrus.Logger) *GRPCServer {
	return &GRPCServer{
		db:       repository.New(db),
		dbLogger: dbLogger,
		cache:    cache,
		Logger:   logger,
	}
}

func (s *GRPCServer) Get(ctx context.Context, _ *emptypb.Empty) (*api.UserList, error) {
	var resp *api.UserList
	var err error

	resp, err = s.cache.Get(ctx, "grpc_test_users")
	if len(resp.Users) != 0 {
		s.Logger.Debug("getting user from cache")
		return resp, nil
	}

	users, err := s.db.GetUsers(ctx)
	if err != nil {
		s.Logger.Error(fmt.Errorf("getting user list from db error: %w", err))
		return nil, err
	}

	for _, v := range users {
		resp.Users = append(resp.Users, v.Encode())
	}

	err = s.cache.Set(ctx, "grpc_test_users", resp)
	if err != nil {
		s.Logger.Error(fmt.Errorf("setting cache into db error: %w", err))
		return nil, err
	}

	return resp, nil
}

func (s *GRPCServer) Add(ctx context.Context, request *api.AddRequest) (*api.Response, error) {
	var user model.User
	user.Decode(request.Body.User)

	id, err := s.db.AddUser(ctx, &user)
	if err != nil {
		s.Logger.Error(fmt.Errorf("adding user into db error: %w", err))
		return nil, err
	}

	err = s.dbLogger.LogRegistration(id, fmt.Sprintf("%s %s", request.Body.User.Name, request.Body.User.Email), time.Now())
	if err != nil {
		s.Logger.Error(fmt.Errorf("adding dbLogger error: %w", err))
		return nil, err
	}

	return &api.Response{
		Message: fmt.Sprintf("User with name %s and email %s was added", request.Body.User.Name, request.Body.User.Email),
	}, nil
}

func (s *GRPCServer) Remove(ctx context.Context, request *api.RemoveRequest) (*api.Response, error) {
	err := s.db.RemoveUser(ctx, request.Body.User.Id)
	if err != nil {
		s.Logger.Error(fmt.Errorf("removing user from db error: %w", err))
		return nil, err
	}

	return &api.Response{
		Message: `User #` + strconv.Itoa(int(request.Body.User.Id)) + ` was removed.`,
	}, nil
}
