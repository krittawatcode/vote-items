package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler struct holds required services for handler to function
type UserHandler struct{}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type UserConfig struct {
	R *gin.Engine
}

// Me handler calls services for getting
// a user's details
func (h *UserHandler) Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it's me",
	})
}

// Sign up handler
// func (h *UserHandler) SignUp(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{
// 		"hello": "it's sign up",
// 	})
// }

// Sign in handler
// func (h *UserHandler) SignIn(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{
// 		"hello": "it's sign in",
// 	})
// }

// Sign out handler
// func (h *UserHandler) SignOut(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{
// 		"hello": "it's sign out",
// 	})
// }

// Tokens handler
// func (h *UserHandler) Tokens(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{
// 		"hello": "it's tokens",
// 	})
// }
