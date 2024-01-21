package route

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/krittawatcode/vote-items/backend/delivery/handler"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	userHandler := &handler.UserHandler{}

	// user handler
	v1 := r.Group(os.Getenv("API_URL"))
	{
		v1.GET("me", userHandler.Me)
		v1.POST("signUp", userHandler.SignUp)
	}

	return r
}
