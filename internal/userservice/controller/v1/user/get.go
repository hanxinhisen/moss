// Created by Hisen at 2022/3/10.
package user

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
)

func (u *UserController) Get(c *gin.Context) {
	user, err := u.srv.Users().Get(c, c.Param("name"), metav1.GetOptions{})

	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	core.WriteResponse(c, nil, user)
}
