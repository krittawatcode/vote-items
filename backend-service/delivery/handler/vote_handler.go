package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/krittawatcode/vote-items/backend-service/delivery/handler/middleware"
	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/apperror"
)

// Handler struct holds required services for handler to function
type VotesHandler struct {
	Router          *gin.Engine
	VoteUseCase     domain.VoteUseCase
	TokenUseCase    domain.TokenUseCase
	Url             string // url for user routes
	TimeoutDuration time.Duration
}

// Does not return as it deals directly with a reference to the gin Engine
func NewVotesHandler(router *gin.Engine, vu domain.VoteUseCase, tu domain.TokenUseCase, url string, timeout time.Duration) {
	h := &VotesHandler{
		VoteUseCase:  vu,
		TokenUseCase: tu,
	}

	// Create an votes group
	g := router.Group(url)
	if gin.Mode() != gin.TestMode {
		// set up middle ware for time out
		g.Use(middleware.Timeout(timeout, apperror.NewServiceUnavailable()))
		g.POST("/", middleware.AuthUser(h.TokenUseCase), h.CastVote)
	}
}

// @Summary Cast a vote
// @Description Cast a vote
// @Tags vote
// @Accept  json
// @Produce  json
// @Param vote body domain.Vote true "Vote payload"
// @Success 201 {object} domain.Vote "Vote successfully cast"
// @Failure 400 {object} domain.ErrorResponse "Bad Request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /api/v1/votes [post]
// POST /api/v1/votes: Cast a vote
func (h *VotesHandler) CastVote(c *gin.Context) {
	var vote domain.Vote
	if err := c.ShouldBindJSON(&vote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": apperror.NewBadRequest(err.Error())})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	vote.UserID = user.(*domain.User).UID

	err := h.VoteUseCase.Create(c.Request.Context(), &vote)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, vote)
}
