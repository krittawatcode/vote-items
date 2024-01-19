package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/krittawatcode/vote-items/backend/delivery/handler"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	userHandler := &handler.UserHandler{}

	// user handler
	v1 := r.Group("/v1")
	{
		v1.GET("me", userHandler.Me)
	}

	return r
}
