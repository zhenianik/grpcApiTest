package grpc

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zhenianik/grpcApiTest/internal"
	"github.com/zhenianik/grpcApiTest/internal/controller/grpc/api"
	"github.com/zhenianik/grpcApiTest/internal/grpc/mock"
	loggerMock "github.com/zhenianik/grpcApiTest/pkg/logger/mock"
)

func service(t *testing.T) (*mock.MockUserService, *loggerMock.MockLogger) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	userRepoMock := mock.NewMockUserService(mockCtl)
	logger := loggerMock.NewMockLogger(mockCtl)

	defer mockCtl.Finish()

	return userRepoMock, logger
}

func TestService_Get(t *testing.T) {
	us, l := service(t)

	users := []*api.User{
		{
			Id:    1,
			Name:  "Vasya",
			Email: "vasya@gmail.com",
		},
		{
			Id:    2,
			Name:  "Zhenia",
			Email: "zhenia@gmail.com",
		},
	}

	tests := []struct {
		name    string
		want    *api.UserList
		prepare func()
	}{
		{
			name: "empty user list",
			want: &api.UserList{},
			prepare: func() {
				us.EXPECT().Get(context.Background()).Return(nil, nil)
			},
		},
		{
			name: "user list with two users",
			want: &api.UserList{Users: users},
			prepare: func() {
				us.EXPECT().Get(context.Background()).Return([]*internal.User{
					{
						Id:    1,
						Name:  "Vasya",
						Email: "vasya@gmail.com",
					},
					{
						Id:    2,
						Name:  "Zhenia",
						Email: "zhenia@gmail.com",
					},
				}, nil)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.prepare()
			service := New(us, l)

			apiUserList, err := service.Get(context.Background(), nil)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, apiUserList)
		})
	}
}

func TestService_Add(t *testing.T) {
	us, l := service(t)

	emptyUser := &internal.User{
		Name:  "",
		Email: "",
	}

	notEmptyUser := &internal.User{
		Name:  "Alex",
		Email: "alex@gmail.com",
	}

	tests := []struct {
		name    string
		user    *api.AddRequest
		want    *api.AddResponse
		err     error
		prepare func()
	}{
		{
			name: "empty user",
			user: &api.AddRequest{},
			want: nil,
			err:  internal.EMPTY_USER_ERROR,
			prepare: func() {
				l.EXPECT().Error("adding user into db error: Empty User name/email").Return().AnyTimes()
				us.EXPECT().Add(context.Background(), emptyUser).Return(internal.UserID(0), internal.EMPTY_USER_ERROR)
			},
		},
		{
			name: "not empty user",
			user: &api.AddRequest{
				Name:  "Alex",
				Email: "alex@gmail.com",
			},
			want: &api.AddResponse{Id: 777},
			err:  nil,
			prepare: func() {
				us.EXPECT().Add(context.Background(), notEmptyUser).Return(internal.UserID(777), nil)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.prepare()
			service := New(us, l)

			resp, err := service.Add(context.Background(), tt.user)
			assert.Equal(t, err, tt.err)
			assert.Equal(t, tt.want, resp)
		})
	}
}

func TestService_Remove(t *testing.T) {
	us, l := service(t)

	tests := []struct {
		name    string
		userID  *api.RemoveRequest
		want    error
		prepare func()
	}{
		{
			name:   "empty id",
			userID: &api.RemoveRequest{Id: 0},
			want:   internal.EMPTY_USER_ERROR,
			prepare: func() {
				l.EXPECT().Error("removing user from db error: Empty User name/email").Return().AnyTimes()
				us.EXPECT().Remove(context.Background(), internal.UserID(0)).Return(internal.EMPTY_USER_ERROR)
			},
		},
		{
			name:   "not empty id",
			userID: &api.RemoveRequest{Id: 777},
			want:   nil,
			prepare: func() {
				us.EXPECT().Remove(context.Background(), internal.UserID(777)).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.prepare()
			service := New(us, l)

			_, err := service.Remove(context.Background(), tt.userID)
			assert.Equal(t, tt.want, err)
		})
	}
}
