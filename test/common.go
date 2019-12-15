package test

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	m "github.com/saeedafshari8/authenticator/middleware"
	"net/http"
	"time"
)

func GetRouter() *gin.Engine {
	r := gin.Default()
	authInfo := &m.AuthInfo{
		IdentityKey:   "testIdentity",
		Realm:         "test Realm",
		Secret:        "secret",
		TokenTimeout:  time.Hour,
		MaxRefresh:    time.Hour,
		Authenticator: MockAuth,
	}
	r.POST("/v1/login", func(c *gin.Context) {
		var login m.Login
		if err := c.ShouldBindJSON(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, login)
	})
	r.Use(m.JwtAuthentication(authInfo).MiddlewareFunc())
	{
		r.POST("/v1/echo", func(c *gin.Context) {
			var login m.Login
			if err := c.ShouldBindJSON(&login); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, login)
		})
	}
	return r
}

func MockAuth(login *m.Login) (*m.Account, error) {
	if (login.Username == "admin" && login.Password == "admin") ||
		(login.Username == "test" && login.Password == "test") {
		return &m.Account{
			Email:     login.Username,
			LastName:  "Afshari",
			FirstName: "Saeed",
		}, nil
	}
	return nil, jwt.ErrFailedAuthentication
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}
