// Created by Hisen at 2022/3/4.
package v1

import "github.com/hanxinhisen/moss/internal/userservice/store"

type Service interface {
	Users() UserSrv
}

type service struct {
	store store.Factory
}

func (s *service) Users() UserSrv {
	return newUsers(s)
}

func NewService(store store.Factory) Service {
	return &service{store: store}
}
