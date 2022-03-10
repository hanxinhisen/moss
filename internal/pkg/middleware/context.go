// Created by Hisen at 2022/3/2.
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/hanxinhisen/moss/pkg/log"
)

const UsernameKey = "Username"

func Context() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set(log.KeyRequestID, context.GetString(XRequestIDKey))
		context.Set(log.KeyUsername, context.GetString(UsernameKey))
		context.Next()
	}
}
