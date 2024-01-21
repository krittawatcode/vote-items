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

// signUpReq is not exported, hence the lowercase name
// it is used for validation and json marshalling
type signUpReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

// Sign up handler
func (h *UserHandler) SignUp(c *gin.Context) {
	// define a variable to which we'll bind incoming
	// json body, {email, password}
	var req signUpReq

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	u := &model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.UserUseCase.SignUp(c, u)

	if err != nil {
		log.Printf("Failed to sign up user: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
}

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
