package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/zhenianik/grpcApiTest/internal/controller/grpc/api"
	"github.com/zhenianik/grpcApiTest/internal/model"
	"github.com/zhenianik/grpcApiTest/internal/usecase"
	"github.com/zhenianik/grpcApiTest/pkg/cache/mock"
	mock2 "github.com/zhenianik/grpcApiTest/pkg/logger/mock"
)

var errInternalServErr = errors.New("internal server error")

type test struct {
	name string
	mock func()
	user *model.User
	res  interface{}
	err  error
}

func service(t *testing.T) (*usecase.Service, *MockUserRepo, *mock.MockCache, *mock2.MockLogger, *MockEventSender) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockUserRepo(mockCtl)
	cache := mock.NewMockCache(mockCtl)
	logger := mock2.NewMockLogger(mockCtl)
	es := NewMockEventSender(mockCtl)

	s := usecase.NewService(repo, es, cache, logger)
	return s, repo, cache, logger, es
}

func TestAdd(t *testing.T) {
	t.Parallel()

	s, repo, _, _, _ := service(t)
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
			res, err := s.Add(ctx, addReq)

			require.Equal(t, res, tc.res)
			require.ErrorIs(t, err, tc.err)
		})
	}
}
