package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/krittawatcode/vote-items/backend/domain/apperrors"
	"github.com/krittawatcode/vote-items/backend/domain/interface/user_interface"
	"github.com/krittawatcode/vote-items/backend/domain/model"
)

// Handler struct holds required services for handler to function
type UserHandler struct {
	UserUseCase user_interface.UserUseCase
}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type UserConfig struct {
	R           *gin.Engine
	UserUseCase user_interface.UserUseCase
}

// Me handler calls services for getting
// a user's details
// Me handler calls services for getting
// a user's details
func (h *UserHandler) Me(c *gin.Context) {
	// A *model.User will eventually be added to context in middleware
	user, exists := c.Get("user")

	// This shouldn't happen, as our middleware ought to throw an error.
	// This is an extra safety measure
	// We'll extract this logic later as it will be common to all handler
	// methods which require a valid user
	if !exists {
		log.Printf("Unable to extract user from request context for unknown reason: %v\n", c)
		err := apperrors.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})

		return
	}

	uid := user.(*model.User).UID

	// gin.Context satisfies go's context.Context interface
	u, err := h.UserUseCase.Get(c, uid)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", uid, err)
		e := apperrors.NewNotFound("user", uid.String())

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": u,
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
