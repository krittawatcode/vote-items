package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/krittawatcode/vote-items/user-service/delivery/handler/middleware"
	"github.com/krittawatcode/vote-items/user-service/domain"
	"github.com/krittawatcode/vote-items/user-service/domain/apperror"
)

// Handler struct holds required services for handler to function
type UserHandler struct {
	Router          *gin.Engine
	UserUseCase     domain.UserUseCase
	TokenUseCase    domain.TokenUseCase
	BaseUrl         string // base url for user routes
	TimeoutDuration time.Duration
}

// Does not return as it deals directly with a reference to the gin Engine
func NewUserHandler(router *gin.Engine, uu domain.UserUseCase, tu domain.TokenUseCase, baseUrl string, timeout time.Duration) {
	// Create a handler (which will later have injected services)
	h := &UserHandler{
		UserUseCase:  uu,
		TokenUseCase: tu,
	}

	// Create an account group
	g := router.Group(baseUrl)

	if gin.Mode() != gin.TestMode {
		g.Use(middleware.Timeout(timeout, apperror.NewServiceUnavailable()))
	}

	// Add a health check endpoint
	g.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "running"})
	})

	if gin.Mode() != gin.TestMode {
		g.Use(middleware.Timeout(timeout, apperror.NewServiceUnavailable()))
		g.GET("/me", middleware.AuthUser(h.TokenUseCase), h.Me)
	} else {
		g.GET("/me", h.Me)
	}

	g.POST("/signUp", h.SignUp)
	g.POST("/singIn", h.SignIn)
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
		err := apperror.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})

		return
	}

	u, ok := user.(*domain.User)
	if !ok {
		log.Printf("User is not of type *domain.User: %v\n", user)
		err := apperror.NewInternal()
		c.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	// use the Request Context
	ctx := c.Request.Context()
	validatedUser, err := h.UserUseCase.Get(ctx, u.UID)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", u.UID, err)
		e := apperror.NewNotFound("user", u.UID.String())

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": validatedUser,
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

	u := &domain.User{
		Email:    req.Email,
		Password: req.Password,
	}
	ctx := c.Request.Context()
	err := h.UserUseCase.SignUp(ctx, u)
	if err != nil {
		log.Printf("Failed to sign up user: %v\n", err.Error())

		c.JSON(apperror.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// create token pair as strings
	tokens, err := h.TokenUseCase.NewPairFromUser(ctx, u, "")
	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())

		// may eventually implement rollback logic here
		// meaning, if we fail to create tokens after creating a user,
		// we make sure to clear/delete the created user in the database

		c.JSON(apperror.Status(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"tokens": tokens,
	})
}

// signInReq is not exported
type signInReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

// SignIn used to authenticate extant user
func (h *UserHandler) SignIn(c *gin.Context) {
	var req signInReq

	if ok := bindData(c, &req); !ok {
		return
	}

	u := &domain.User{
		Email:    req.Email,
		Password: req.Password,
	}

	ctx := c.Request.Context()
	err := h.UserUseCase.SignIn(ctx, u)

	if err != nil {
		log.Printf("Failed to sign in user: %v\n", err.Error())
		c.JSON(apperror.Status(err), gin.H{
			"error": err,
		})
		return
	}

	tokens, err := h.TokenUseCase.NewPairFromUser(ctx, u, "")

	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())

		c.JSON(apperror.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}
