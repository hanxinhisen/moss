// Created by Hisen at 2022/3/3.
package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/hanxinhisen/moss/internal/pkg/code"
	"github.com/hanxinhisen/moss/internal/pkg/middleware"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"
	"strings"
)

const authHeaderCount = 2

type AutoStrategy struct {
	basic middleware.AuthStrategy
	jwt   middleware.AuthStrategy
}

func NewAutoStrategy(basic, jwt middleware.AuthStrategy) AutoStrategy {
	return AutoStrategy{
		basic: basic,
		jwt:   jwt,
	}
}
func (a AutoStrategy) AuthFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		operator := middleware.AuthOperator{}
		authHeader := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(authHeader) != authHeaderCount {
			core.WriteResponse(
				c,
				errors.WithCode(code.ErrInvalidAuthHeader, "Authorization header format is wrong."),
				nil,
			)
			c.Abort()

			return
		}

		switch authHeader[0] {
		case "Basic":
			operator.SetStrategy(a.basic)
		case "Bearer":
			operator.SetStrategy(a.jwt)
			// a.JWT.MiddlewareFunc()(c)
		default:
			core.WriteResponse(c, errors.WithCode(code.ErrSignatureInvalid, "unrecognized Authorization header."), nil)
			c.Abort()

			return
		}

		operator.AuthFunc()(c)

		c.Next()
	}
}

var _ middleware.AuthStrategy = &AutoStrategy{}
