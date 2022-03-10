// Created by Hisen at 2022/3/3.
package user

import (
	srvv1 "github.com/hanxinhisen/moss/internal/userservice/service/v1"
	"github.com/hanxinhisen/moss/internal/userservice/store"
)

type UserController struct {
	srv srvv1.Service
}

func NewUserController(store store.Factory) *UserController {
	return &UserController{srv: srvv1.NewService(store)}
}
