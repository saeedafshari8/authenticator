package middleware

import (
	"fmt"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type AuthInfo struct {
	IdentityKey   string
	Realm         string
	Secret        string
	TokenTimeout  time.Duration
	MaxRefresh    time.Duration
	Authenticator func(login *Login) (*Account, error)
	Authorizator  func(account *Account) (bool, error)
}

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type HttpStatusResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var JwtAuthentication = func(authInfo *AuthInfo) *jwt.GinJWTMiddleware {

	if (*authInfo).TokenTimeout == 0 {
		(*authInfo).TokenTimeout = time.Hour
	}
	if (*authInfo).MaxRefresh == 0 {
		(*authInfo).MaxRefresh = time.Hour
	}

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       (*authInfo).Realm,
		Key:         []byte((*authInfo).Secret),
		Timeout:     (*authInfo).TokenTimeout,
		MaxRefresh:  (*authInfo).MaxRefresh,
		IdentityKey: (*authInfo).IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if account, ok := data.(*Account); ok {
				return jwt.MapClaims{
					(*authInfo).IdentityKey: account.Email,
					"firstName":             account.FirstName,
					"lastName":              account.LastName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &Account{
				Email:     claims[(*authInfo).IdentityKey].(string),
				FirstName: claims["firstName"].(string),
				LastName:  claims["lastName"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var login *Login
			if err := c.ShouldBind(login); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			account, err := (*authInfo).Authenticator(login)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return account, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if account, ok := data.(*Account); ok {
				authorized, err := (*authInfo).Authorizator(account)
				return err == nil && authorized
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, HttpStatusResponse{
				Code:    fmt.Sprintf("%d", code),
				Message: message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return authMiddleware
}
