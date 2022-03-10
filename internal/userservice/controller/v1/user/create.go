// Created by Hisen at 2022/3/4.
package user

import (
	"github.com/gin-gonic/gin"
	"github.com/hanxinhisen/moss/internal/pkg/code"
	v1 "github.com/hanxinhisen/moss/internal/userservice/store/model/v1"
	"github.com/hanxinhisen/moss/pkg/log"
	"github.com/marmotedu/component-base/pkg/auth"
	"github.com/marmotedu/component-base/pkg/core"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/errors"
)

func (u *UserController) Create(c *gin.Context) {
	log.L(c).Info("user create function called.")
	var r v1.User

	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errors.WithCode(code.ErrBind, err.Error()), nil)
		return
	}

	if errs := r.Validate(); len(errs) != 0 {
		core.WriteResponse(c, errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), nil)
		return

	}

	r.Password, _ = auth.Encrypt(r.Password)
	r.Status = 1

	if err := u.srv.Users().Create(c, &r, metav1.CreateOptions{}); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, r)
}
