// Created by Hisen at 2022/3/3.
package userservice

import (
	"context"
	"encoding/base64"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/hanxinhisen/moss/internal/pkg/middleware"
	"github.com/hanxinhisen/moss/internal/pkg/middleware/auth"
	"github.com/hanxinhisen/moss/internal/userservice/store"
	v1 "github.com/hanxinhisen/moss/internal/userservice/store/model/v1"
	"github.com/hanxinhisen/moss/pkg/log"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/spf13/viper"

	"net/http"
	"strings"
	"time"
)

func newJWTAuth() middleware.AuthStrategy {
	ginjwt, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            viper.GetString("jwt.Realm"),
		SigningAlgorithm: "HS256",
		Key:              []byte(viper.GetString("jwt.key")),
		Timeout:          viper.GetDuration("jwt.timeout"),
		MaxRefresh:       viper.GetDuration("jwt.max-refresh"),
		Authenticator:    authenticator(),
		LoginResponse:    loginResponse(),
		LogoutResponse: func(c *gin.Context, code int) {
			c.JSON(http.StatusOK, nil)
		},
		RefreshResponse: refreshResponse(),
		PayloadFunc:     payloadFunc(),
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)

			return claims[jwt.IdentityKey]
		},
		IdentityKey:  middleware.UsernameKey,
		Authorizator: authorizator(),
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		SendCookie:    true,
		TimeFunc:      time.Now,
		// TODO: HTTPStatusMessageFunc:
	})

	return auth.NewJWTStrategy(*ginjwt)
}

func authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if v, ok := data.(string); ok {
			log.L(c).Infof("user `%s` is authenticated.", v)

			return true
		}

		return false
	}
}

const (
	APIServerAudience = "iam.api.marmotedu.com"

	APIServerIssuer = "iam-apiserver"
)

func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		claims := jwt.MapClaims{
			"iss": APIServerIssuer,
			"aud": APIServerAudience,
		}
		if u, ok := data.(*v1.User); ok {
			claims[jwt.IdentityKey] = u.Name
			claims["sub"] = u.Name
		}

		return claims
	}
}

func refreshResponse() func(c *gin.Context, code int, message string, time time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(http.StatusOK, gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

func loginResponse() func(c *gin.Context, code int, message string, time time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		c.JSON(http.StatusOK, gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	}
}

func newAutoAuth() middleware.AuthStrategy {
	return auth.NewAutoStrategy(newBasicAuth().(auth.BasicStrategy), newJWTAuth().(auth.JWTStrategy))
}

func newBasicAuth() middleware.AuthStrategy {
	return auth.NewBasicStrategy(func(username string, password string) bool {
		user, err := store.Client().Users().Get(context.TODO(), username, metav1.GetOptions{})
		if err != nil {
			return false
		}
		if err := user.Compare(password); err != nil {
			return false
		}
		user.LoginedAt = time.Now()
		_ = store.Client().Users().Update(context.TODO(), user, metav1.UpdateOptions{})
		return true
	})

}

type loginInfo struct {
	Username string `form:"username" json:"username" binding:"required,username"`
	Password string `form:"password" json:"password" binding:"required,password"`
}

func authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		log.Infof("-------------")

		var login loginInfo
		var err error

		if c.Request.Header.Get("Authorization") != "" {
			login, err = parseWithHeader(c)
		} else {
			log.Infof("-------------")

			login, err = parseWithBody(c)
		}
		log.Infof("-------------")

		if err != nil {
			return "", jwt.ErrFailedAuthentication
		}
		user, err := store.Client().Users().Get(c, login.Username, metav1.GetOptions{})
		if err != nil {
			log.Errorf("get user information failed: %s", err.Error())

			return "", jwt.ErrFailedAuthentication
		}

		// Compare the login password with the user password.
		if err := user.Compare(login.Password); err != nil {
			return "", jwt.ErrFailedAuthentication
		}

		user.LoginedAt = time.Now()
		_ = store.Client().Users().Update(c, user, metav1.UpdateOptions{})

		return user, nil
	}
}

func parseWithHeader(c *gin.Context) (loginInfo, error) {
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		log.Errorf("get basic string from Authorization header failed")
		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	payload, err := base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		log.Errorf("decode basic string: %s", err.Error())

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		log.Errorf("parse payload failed")

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	return loginInfo{
		Username: pair[0],
		Password: pair[1],
	}, nil

}

func parseWithBody(c *gin.Context) (loginInfo, error) {
	var login loginInfo
	log.Infof("-------------")

	if err := c.ShouldBindJSON(&login); err != nil {
		log.Infof("-------------")

		log.Errorf("parse login parameters: %s", err.Error())

		return loginInfo{}, jwt.ErrFailedAuthentication
	}

	return login, nil
}
