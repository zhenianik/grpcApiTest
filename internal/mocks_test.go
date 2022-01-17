// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package internal_test is a generated GoMock package.
package internal_test

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	internal "github.com/zhenianik/grpcApiTest/internal"
)

// MockUserStorage is a mock of UserStorage interface.
type MockUserStorage struct {
	ctrl     *gomock.Controller
	recorder *MockUserStorageMockRecorder
}

// MockUserStorageMockRecorder is the mock recorder for MockUserStorage.
type MockUserStorageMockRecorder struct {
	mock *MockUserStorage
}

// NewMockUserStorage creates a new mock instance.
func NewMockUserStorage(ctrl *gomock.Controller) *MockUserStorage {
	mock := &MockUserStorage{ctrl: ctrl}
	mock.recorder = &MockUserStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserStorage) EXPECT() *MockUserStorageMockRecorder {
	return m.recorder
}

// AddUser mocks base method.
func (m *MockUserStorage) AddUser(ctx context.Context, user *internal.User) (internal.UserID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", ctx, user)
	ret0, _ := ret[0].(internal.UserID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUser indicates an expected call of AddUser.
func (mr *MockUserStorageMockRecorder) AddUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockUserStorage)(nil).AddUser), ctx, user)
}

// GetUsers mocks base method.
func (m *MockUserStorage) GetUsers(ctx context.Context) ([]*internal.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", ctx)
	ret0, _ := ret[0].([]*internal.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockUserStorageMockRecorder) GetUsers(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockUserStorage)(nil).GetUsers), ctx)
}

// RemoveUser mocks base method.
func (m *MockUserStorage) RemoveUser(ctx context.Context, id internal.UserID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUser", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUser indicates an expected call of RemoveUser.
func (mr *MockUserStorageMockRecorder) RemoveUser(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUser", reflect.TypeOf((*MockUserStorage)(nil).RemoveUser), ctx, id)
}

// MockEventSender is a mock of EventSender interface.
type MockEventSender struct {
	ctrl     *gomock.Controller
	recorder *MockEventSenderMockRecorder
}

// MockEventSenderMockRecorder is the mock recorder for MockEventSender.
type MockEventSenderMockRecorder struct {
	mock *MockEventSender
}

// NewMockEventSender creates a new mock instance.
func NewMockEventSender(ctrl *gomock.Controller) *MockEventSender {
	mock := &MockEventSender{ctrl: ctrl}
	mock.recorder = &MockEventSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventSender) EXPECT() *MockEventSenderMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockEventSender) Send(id internal.UserID, name string, time time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", id, name, time)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockEventSenderMockRecorder) Send(id, name, time interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockEventSender)(nil).Send), id, name, time)
}
