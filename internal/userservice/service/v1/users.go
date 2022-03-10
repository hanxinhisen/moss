// Created by Hisen at 2022/3/3.
package v1

import (
	"context"
	"github.com/hanxinhisen/moss/internal/pkg/code"
	"github.com/hanxinhisen/moss/internal/userservice/store"
	v1 "github.com/hanxinhisen/moss/internal/userservice/store/model/v1"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/errors"
	"regexp"
)

type UserSrv interface {
	Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error
	Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) error
	Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error
	Get(ctx context.Context, username string, opts metav1.GetOptions) (*v1.User, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error)
	ListWithBadPerformance(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error)
	ChangePassword(ctx context.Context, user *v1.User) error
}

type userService struct {
	store store.Factory
}

func newUsers(srv *service) *userService {
	return &userService{store: srv.store}
}
func (u *userService) Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) error {
	if err := u.store.Users().Create(ctx, user, opts); err != nil {
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'idx_name'", err.Error()); match {
			return errors.WithCode(code.ErrUserAlreadyExist, err.Error())
		}
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

func (u userService) Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) error {
	//TODO implement me
	return nil
}

func (u userService) Delete(ctx context.Context, username string, opts metav1.DeleteOptions) error {
	//TODO implement me
	return nil
}

func (u userService) DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error {
	//TODO implement me
	panic("implement me")
}

func (u userService) Get(ctx context.Context, username string, opts metav1.GetOptions) (*v1.User, error) {
	user, err := u.store.Users().Get(ctx, username, opts)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u userService) List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error) {
	//TODO implement me
	return nil, nil
}

func (u userService) ListWithBadPerformance(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error) {
	//TODO implement me
	return nil, nil
}

func (u userService) ChangePassword(ctx context.Context, user *v1.User) error {
	//TODO implement me
	return nil
}

var _ UserSrv = (*userService)(nil)
