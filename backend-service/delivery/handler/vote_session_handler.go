package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/krittawatcode/vote-items/backend-service/delivery/handler/middleware"
	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/apperror"
)

// Handler struct holds required services for handler to function
type VoteSessionsHandler struct {
	Router             *gin.Engine
	VoteSessionUseCase domain.VoteSessionUseCase
	TokenUseCase       domain.TokenUseCase
	Url                string // base url for vote session routes
	TimeoutDuration    time.Duration
}

// Does not return as it deals directly with a reference to the gin Engine
func NewVoteSessionsHandler(router *gin.Engine, vsu domain.VoteSessionUseCase, tu domain.TokenUseCase, url string, timeout time.Duration) {
	h := &VoteSessionsHandler{
		VoteSessionUseCase: vsu,
		TokenUseCase:       tu,
	}

	// Create an vote-sessions group
	g := router.Group(url)

	if gin.Mode() != gin.TestMode {
		g.Use(middleware.Timeout(timeout, apperror.NewServiceUnavailable()))
	}

	if gin.Mode() != gin.TestMode {
		// set up middle ware for time out
		g.Use(middleware.Timeout(timeout, apperror.NewServiceUnavailable()))
		// get current open vote session
		g.GET("/vote_sessions/open", middleware.AuthUser(h.TokenUseCase), h.GetOpenVoteSession)
		// create a new vote session
		g.PUT("/vote_sessions/:id/open", middleware.AuthUser(h.TokenUseCase), h.OpenVoteSession)
		// close a vote session
		g.PUT("/vote_sessions/:id/close", middleware.AuthUser(h.TokenUseCase), h.CloseVoteSession)
	}
}

// PUT /vote_sessions/:id/open // Open a vote session
// OpenVoteSession opens a vote session
// @Summary Open a vote session
// @Description Open a vote session by ID
// @Tags vote_sessions
// @Accept  json
// @Produce  json
// @Param id path int true "Session ID"
// @Success 200 {object} domain.SuccessResponse "Vote session opened successfully"
// @Failure 400 {object} domain.ErrorResponse "Bad Request"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /api/v1/vote_sessions/:id/open [put]
func (h *VoteSessionsHandler) OpenVoteSession(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	err = h.VoteSessionUseCase.OpenVoteSession(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Vote session opened successfully"})
}

// @Summary Get open vote session
// @Description Retrieve the currently open vote session
// @Tags vote_sessions
// @Produce  json
// @Success 200 {object} domain.VoteSession "Successfully retrieved the open vote session"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /api/v1/vote_sessions/open [get]
// GET /api/v1/vote_sessions/open: Get open vote session
func (h *VoteSessionsHandler) GetOpenVoteSession(c *gin.Context) {
	voteSession, err := h.VoteSessionUseCase.GetOpenVoteSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, voteSession)
}

// @Summary Close a vote session
// @Description Close a vote session by ID
// @Tags vote_sessions
// @Accept  json
// @Produce  json
// @Param   id     path    int     true    "Vote Session ID"
// @Success 200 {object} domain.SuccessResponse "Vote session closed successfully"
// @Failure 400 {object} domain.ErrorResponse "Bad Request"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /api/v1/vote_sessions/{id}/close [put]
// PUT /api/v1/vote_sessions/{id}/close: Close a vote session
func (h *VoteSessionsHandler) CloseVoteSession(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	err = h.VoteSessionUseCase.CloseVoteSession(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Vote session closed successfully"})
}
