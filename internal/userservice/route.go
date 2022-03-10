// Created by Hisen at 2022/3/3.
package userservice

import (
	"github.com/gin-gonic/gin"
	"github.com/hanxinhisen/moss/internal/pkg/code"
	"github.com/hanxinhisen/moss/internal/pkg/middleware"
	"github.com/hanxinhisen/moss/internal/pkg/middleware/auth"
	"github.com/hanxinhisen/moss/internal/userservice/controller/v1/user"
	"github.com/hanxinhisen/moss/internal/userservice/store/mysql"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"

	_ "github.com/hanxinhisen/moss/pkg/validator"
)

func initRoute(g *gin.Engine) {
	installMiddleware(g)
	installController(g)

}

func installController(g *gin.Engine) *gin.Engine {
	jwtStrategy, _ := newJWTAuth().(auth.JWTStrategy)
	g.POST("/login", jwtStrategy.LoginHandler)
	g.POST("/logout", jwtStrategy.LogoutHandler)
	g.POST("/refresh", jwtStrategy.RefreshHandler)

	auto := newAutoAuth()
	g.NoRoute(auto.AuthFunc(), func(context *gin.Context) {
		core.WriteResponse(context, errors.WithCode(code.ErrPageNotFound, "Page Not Found"), nil)
	})

	storeIns, _ := mysql.GetMysqlFactoryOr(nil)
	v1 := g.Group("/v1")

	{
		// todo route miss
		userv1 := v1.Group("/users")
		{
			controller := user.NewUserController(storeIns)
			userv1.POST("", controller.Create)
			userv1.Use(auto.AuthFunc(), middleware.Validation())
			userv1.GET(":name", controller.Get)

		}
		v1.Use(auto.AuthFunc())
	}
	return g
}
func installMiddleware(g *gin.Engine) {
}
