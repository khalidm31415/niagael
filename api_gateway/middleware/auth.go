package middelware

import (
	"api_gateway/entity"
	"api_gateway/service"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var IdentityKey = "id"

type LoginInput struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func NewAuthtMiddleware(authService service.AuthService) *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "niagael",
		Key:         []byte("JWT_SECRET_KEY"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: IdentityKey,
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals LoginInput
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			user, err := authService.Authenticate(username, password)
			if err != nil {
				zap.S().Error(err)

			}
			return user.ID, nil

		},
		PayloadFunc: func(userID interface{}) jwt.MapClaims {
			if userID != nil {
				return jwt.MapClaims{
					IdentityKey: userID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			userID := claims[IdentityKey].(string)
			user, err := authService.GetByID(userID)
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			if err != nil {
				zap.S().Error(err)
				return err
			}
			return *user
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(entity.User); ok {
				return true
			}
			return false
		},
		TimeFunc:       time.Now,
		SendCookie:     true,
		SecureCookie:   false,
		CookieHTTPOnly: true,
		CookieDomain:   "localhost",
		CookieName:     "token",
		TokenLookup:    "cookie:token",
		CookieSameSite: http.SameSiteDefaultMode,
	})

	if err != nil {
		zap.S().Fatal(err)
	}

	if err := authMiddleware.MiddlewareInit(); err != nil {
		zap.S().Fatal(err)
	}

	return authMiddleware
}
