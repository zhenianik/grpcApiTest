package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zhenianik/grpcApiTest/pkg/cache"
)

type UserService struct {
	Storage UserStorage
	Sender  EventSender
	Cache   cache.Cache
}

func New(storage UserStorage, sender EventSender, cache cache.Cache) *UserService {
	return &UserService{
		Storage: storage,
		Sender:  sender,
		Cache:   cache,
	}
}

func (ur *UserService) Get(ctx context.Context) ([]*User, error) {
	var resp []byte
	var err error
	var userList []*User

	resp, err = ur.Cache.Get(ctx, "grpc_test_users")
	if len(resp) != 0 {
		err = json.Unmarshal(resp, &userList)
		if err != nil {
			return nil, fmt.Errorf("unmarshal userlist error: %w", err)
		}
		return userList, nil
	}

	userList, err = ur.Storage.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting user error: %w", err)
	}

	js, err := json.Marshal(userList)
	if err != nil {
		return nil, fmt.Errorf("marshal userlist error: %w", err)
	}

	err = ur.Cache.Set(ctx, "grpc_test_users", js)
	if err != nil {
		return nil, fmt.Errorf("setting cache error: %w", err)
	}

	return userList, nil

}

func (ur *UserService) Add(ctx context.Context, user *User) (UserID, error) {
	if user.Name == "" || user.Email == "" {
		return 0, fmt.Errorf("adding user into db error: %w", EMPTY_USER_ERROR)
	}

	id, err := ur.Storage.AddUser(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("adding user into db error: %w", err)
	}

	err = ur.Cache.Delete(ctx, "grpc_test_users")
	if err != nil {
		return id, fmt.Errorf("deleting from cache error: %w", err)
	}

	err = ur.Sender.Send(id, fmt.Sprintf("%s %s", user.Name, user.Email), time.Now())
	if err != nil {
		return id, fmt.Errorf("eventsender adding error: %w", err)
	}

	return id, nil
}

func (ur *UserService) Remove(ctx context.Context, id UserID) error {
	err := ur.Storage.RemoveUser(ctx, id)
	if err != nil {
		return fmt.Errorf("removing user from db error: %w", err)
	}

	err = ur.Cache.Delete(ctx, "grpc_test_users")
	if err != nil {
		return fmt.Errorf("deleting from cache error: %w", err)
	}

	err = ur.Sender.Send(id, "", time.Now())
	if err != nil {
		return fmt.Errorf("eventsender removing error: %w", err)
	}

	return nil
}
