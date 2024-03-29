package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/krittawatcode/vote-items/backend-service/delivery/handler/middleware"
	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/apperror"
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
	// Create a handler (which will later have injected usecase)
	h := &UserHandler{
		UserUseCase:  uu,
		TokenUseCase: tu,
	}

	// Create an user group
	// ug = user group shorthand
	ug := router.Group(baseUrl)

	if gin.Mode() != gin.TestMode {
		ug.Use(middleware.Timeout(timeout, apperror.NewServiceUnavailable()))
		ug.GET("/me", middleware.AuthUser(h.TokenUseCase), h.Me)
	} else {
		ug.GET("/me", h.Me)
	}

	ug.POST("/signUp", h.SignUp)
	ug.POST("/singIn", h.SignIn)
	ug.POST("/tokens", h.Tokens)
}

// @Summary Get user details
// @Description Get details of the current user
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} domain.User "Successfully retrieved user details"
// @Failure 400 {object} domain.ErrorResponse "Bad Request"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /users/me [get]
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

// @Summary Sign up a new user
// @Description Sign up a new user with email and password
// @Tags users
// @Accept  json
// @Produce  json
// @Param   email     body    string     true    "Email"
// @Param   password  body    string     true    "Password"
// @Success 201 {object} domain.TokenPair "Successfully signed up and returned tokens"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users/signUp [post]
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

// @Summary Sign in an existing user
// @Description Sign in an existing user with email and password
// @Tags users
// @Accept  json
// @Produce  json
// @Param   email     body    string     true    "Email"
// @Param   password  body    string     true    "Password"
// @Success 200 {object} domain.TokenPair "Successfully signed in and returned tokens"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users/signIn [post]
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

type tokensReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// Tokens handler
func (h *UserHandler) Tokens(c *gin.Context) {
	// bind JSON to req of type tokensRew
	var req tokensReq

	if ok := bindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()

	// verify refresh JWT
	refreshToken, err := h.TokenUseCase.ValidateRefreshToken(req.RefreshToken)

	if err != nil {
		c.JSON(apperror.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// get up-to-date user
	u, err := h.UserUseCase.Get(ctx, refreshToken.UID)

	if err != nil {
		c.JSON(apperror.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// create fresh pair of tokens
	tokens, err := h.TokenUseCase.NewPairFromUser(ctx, u, refreshToken.ID.String())

	if err != nil {
		log.Printf("Failed to create tokens for user: %+v. Error: %v\n", u, err.Error())

		c.JSON(apperror.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tokens": tokens,
	})
}
