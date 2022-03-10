// Created by Hisen at 2022/3/3.
package auth

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/hanxinhisen/moss/internal/pkg/code"
	"github.com/hanxinhisen/moss/internal/pkg/middleware"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"
	"strings"
)

type BasicStrategy struct {
	compare func(username, password string) bool
}

func NewBasicStrategy(compare func(username string, password string) bool) BasicStrategy {
	return BasicStrategy{
		compare: compare,
	}
}

func (b BasicStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			core.WriteResponse(
				c,
				errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong."),
				nil,
			)
			c.Abort()

			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !b.compare(pair[0], pair[1]) {
			core.WriteResponse(
				c,
				errors.WithCode(code.ErrSignatureInvalid, "Authorization header format is wrong."),
				nil,
			)
			c.Abort()

			return
		}

		c.Set(middleware.UsernameKey, pair[0])

		c.Next()
	}
}

var _ middleware.AuthStrategy = &BasicStrategy{}
