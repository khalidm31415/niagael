package router

import (
	"api_gateway/entity"
	middelware "api_gateway/middleware"
	"api_gateway/service"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRouter(authService service.AuthService) *gin.Engine {
	authMiddleware := middelware.NewAuthtMiddleware(authService)

	r := gin.Default()

	type SignupInput struct {
		Username    string `binding:"required" json:"username"`
		DisplayName string `binding:"required" json:"display_name"`
		Password    string `binding:"required" json:"password"`
	}
	r.POST("/auth/signup", func(c *gin.Context) {
		var input SignupInput
		if err := c.ShouldBindJSON(&input); err != nil {
			zap.S().Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := authService.Signup(input.Username, input.Password, input.DisplayName); err != nil {
			zap.S().Error(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	})
	r.POST("/auth/login", authMiddleware.LoginHandler)
	r.GET("/auth/current-user", authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		user, _ := c.Get(middelware.IdentityKey)
		c.JSON(200, gin.H{
			"userID":   user.(entity.User).ID,
			"username": user.(entity.User).Username,
		})
	})
	r.POST("/auth/logout", authMiddleware.LogoutHandler)

	r.POST("/product/graphql", authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		user, _ := c.Get(middelware.IdentityKey)
		userID := user.(entity.User).ID
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = "product_api:2719"
			req.Host = "product_api:2719"
			req.Header["X-user-id"] = []string{userID}
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)

	})

	return r
}
