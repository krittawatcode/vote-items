package route

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/krittawatcode/vote-items/user-service/delivery/handler"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	userHandler := &handler.UserHandler{}

	// user handler
	v1 := r.Group(os.Getenv("USER_SERVICE_API_URL"))
	{
		v1.GET("health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"message": "User service is running",
			})
		})
		v1.GET("me", userHandler.Me)
		v1.POST("signUp", userHandler.SignUp)
	}

	return r
}
