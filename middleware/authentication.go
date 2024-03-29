package middleware

import (
	"errors"
	"fmt"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type AuthInfo struct {
	IdentityKey          string
	Realm                string
	Secret               string
	LoginEndPoint        string
	RefreshTokenEndPoint string
	TokenTimeout         time.Duration
	MaxRefresh           time.Duration
	Authenticator        func(login *Login) (*Account, error)
	Authorizator         func(account *Account, request *http.Request) (bool, error)
}

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type LoginResponse struct {
	Code   int       `form:"code" json:"code" binding:"required"`
	Token  string    `form:"token" json:"token" binding:"required"`
	Expire time.Time `form:"expire" json:"expire" binding:"required"`
}

type HttpStatusResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var JwtAuthentication = func(authInfo *AuthInfo, router *gin.Engine) *jwt.GinJWTMiddleware {

	setDefaults(authInfo)

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           (*authInfo).Realm,
		Key:             []byte((*authInfo).Secret),
		Timeout:         (*authInfo).TokenTimeout,
		MaxRefresh:      (*authInfo).MaxRefresh,
		IdentityKey:     (*authInfo).IdentityKey,
		PayloadFunc:     payloadFunc(authInfo),
		IdentityHandler: identityHandler(authInfo),
		Authenticator:   authenticator(authInfo),
		Authorizator:    authorizator(authInfo),
		Unauthorized:    unauthorizedHandler(),
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// Refresh time can be longer than token timeout
	router.GET((*authInfo).RefreshTokenEndPoint, authMiddleware.RefreshHandler)
	router.POST((*authInfo).LoginEndPoint, authMiddleware.LoginHandler)

	return authMiddleware
}

func unauthorizedHandler() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, HttpStatusResponse{
			Code:    fmt.Sprintf("%d", code),
			Message: message,
		})
	}
}

func authorizator(authInfo *AuthInfo) func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if account, ok := data.(*Account); ok {
			if (*authInfo).Authorizator != nil {
				authorized, err := (*authInfo).Authorizator(account, c.Request)
				return err == nil && authorized
			}
			return true
		}
		return false
	}
}

func authenticator(authInfo *AuthInfo) func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var login Login
		if err := c.ShouldBind(&login); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		account, err := (*authInfo).Authenticator(&login)
		if err != nil {
			return nil, jwt.ErrFailedAuthentication
		}
		return account, nil
	}
}

func identityHandler(authInfo *AuthInfo) func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		return &Account{
			Email:     claims[(*authInfo).IdentityKey].(string),
			FirstName: claims["firstName"].(string),
			LastName:  claims["lastName"].(string),
			UserName:  claims["userName"].(string),
		}
	}
}

func payloadFunc(authInfo *AuthInfo) func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if account, ok := data.(*Account); ok {
			return jwt.MapClaims{
				(*authInfo).IdentityKey: account.Email,
				"userName":              account.UserName,
				"firstName":             account.FirstName,
				"lastName":              account.LastName,
			}
		}
		return jwt.MapClaims{}
	}
}

func setDefaults(authInfo *AuthInfo) {
	if (*authInfo).TokenTimeout == 0 {
		(*authInfo).TokenTimeout = time.Hour
	}
	if (*authInfo).MaxRefresh == 0 {
		(*authInfo).MaxRefresh = time.Hour
	}

	if (*authInfo).LoginEndPoint == "" {
		(*authInfo).LoginEndPoint = "/auth/login"
	}

	if (*authInfo).RefreshTokenEndPoint == "" {
		(*authInfo).RefreshTokenEndPoint = "/auth/refresh_token"
	}

	if (*authInfo).Authenticator == nil {
		panic(errors.New("authenticator is not provided"))
	}
}
