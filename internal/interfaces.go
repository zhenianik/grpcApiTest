package internal

import (
	"context"
	"time"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=internal_test

type (
	UserStorage interface {
		GetUsers(ctx context.Context) (users []*User, err error)
		AddUser(ctx context.Context, user *User) (id UserID, err error)
		RemoveUser(ctx context.Context, id UserID) error
	}

	EventSender interface {
		Send(id UserID, name string, time time.Time) error
	}
)
