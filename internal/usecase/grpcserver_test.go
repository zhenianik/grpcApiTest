package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/zhenianik/grpcApiTest/internal/controller/api"
	"github.com/zhenianik/grpcApiTest/internal/model"
	"github.com/zhenianik/grpcApiTest/internal/usecase"
	dBlogger_mock "github.com/zhenianik/grpcApiTest/pkg/dbLogger/mocks"
	logger_mock "github.com/zhenianik/grpcApiTest/pkg/logger/mocks"
)

var errInternalServErr = errors.New("internal server error")

type test struct {
	name string
	mock func()
	user *model.User
	res  interface{}
	err  error
}

func grpcserver(t *testing.T) (*usecase.GRPCServer, *MockUserRepo, *MockCache, *logger_mock.MockInterface, *dBlogger_mock.MockInterface) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockUserRepo(mockCtl)
	cache := NewMockCache(mockCtl)
	logger := logger_mock.NewMockInterface(mockCtl)
	dBlogger := dBlogger_mock.NewMockInterface(mockCtl)

	grpcserver := usecase.NewGRPCServer(repo, dBlogger, cache, logger)
	return grpcserver, repo, cache, logger, dBlogger
}

func TestAdd(t *testing.T) {
	t.Parallel()

	grpcserver, repo, _, _, _ := grpcserver(t)
	user := &model.User{
		Id:    1,
		Name:  "Vasya",
		Email: "vasya@gmail.com",
	}

	addReq := &api.AddRequest{Body: &api.AddRequest_Body{User: &api.AddUserReq{Name: user.Name, Email: user.Email}}}
	ctx := context.Background()

	tests := []test{
		{
			name: "user id = 1",
			mock: func() {
				repo.EXPECT().AddUser(ctx, addReq).Return(user.Id, nil)
			},
			user: user,
			res:  user.Id,
			err:  nil,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tc.mock()
			res, err := grpcserver.Add(ctx, addReq)

			require.Equal(t, res, tc.res)
			require.ErrorIs(t, err, tc.err)
		})
	}
}
