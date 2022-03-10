// Created by Hisen at 2022/3/3.
package auth

import (
	ginjwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/hanxinhisen/moss/internal/pkg/middleware"
)

const AuthzAudience = "cmdb.moss.hanxin.cn"

type JWTStrategy struct {
	ginjwt.GinJWTMiddleware
}

func (j JWTStrategy) AuthFunc() gin.HandlerFunc {
	return j.MiddlewareFunc()
}

var _ middleware.AuthStrategy = &JWTStrategy{}

func NewJWTStrategy(jwtMiddleware ginjwt.GinJWTMiddleware) JWTStrategy {
	return JWTStrategy{jwtMiddleware}
}
