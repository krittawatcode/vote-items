package handler

import (
	"bytes"
	"context"
	"encoding/csv"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/backend-service/delivery/handler/middleware"
	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/apperror"
)

// Handler struct holds required services for handler to function
type VoteItemsHandler struct {
	Router             *gin.Engine
	VoteItemUseCase    domain.VoteItemUseCase
	VoteSessionUseCase domain.VoteSessionUseCase
	VoteUseCase        domain.VoteUseCase
	TokenUseCase       domain.TokenUseCase
	BaseUrl            string // base url for user routes
	TimeoutDuration    time.Duration
}

// Does not return as it deals directly with a reference to the gin Engine
func NewVoteItemsHandler(router *gin.Engine, viu domain.VoteItemUseCase, vu domain.VoteUseCase, vsu domain.VoteSessionUseCase, tu domain.TokenUseCase, baseUrl string, timeout time.Duration) {
	h := &VoteItemsHandler{
		VoteItemUseCase:    viu,
		VoteSessionUseCase: vsu,
		VoteUseCase:        vu,
		TokenUseCase:       tu,
	}

	// Create an vote-items group
	g := router.Group(baseUrl)

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
		// get all active vote items
		g.GET("/vote_items", middleware.AuthUser(h.TokenUseCase), h.FetchActiveVoteItems)
		// create a new vote item
		g.POST("/vote_items", middleware.AuthUser(h.TokenUseCase), h.CreateVoteItem)
		// update a vote item
		g.PUT("/vote_items/:id", middleware.AuthUser(h.TokenUseCase), h.UpdateVoteItem)
		// delete a vote item
		g.DELETE("/vote_items/:id", middleware.AuthUser(h.TokenUseCase), h.DeleteVoteItem)
		// clear all vote items
		g.DELETE("/vote_items", middleware.AuthUser(h.TokenUseCase), h.ClearVoteItem)
		// cast a vote
		g.POST("/votes", middleware.AuthUser(h.TokenUseCase), h.CastVote)
		// get vote results by session id
		// GET /api/v1/vote_results/{session_id}: Get vote results by session id
		// GET /api/v1/vote_results/{session_id}?format=csv: Get vote results by session id in CSV format
		g.GET("/vote_results/:session_id", middleware.AuthUser(h.TokenUseCase), h.GetVoteResultsBySession)
	}
}

// GET /api/v1/vote_items: Get all active vote items
func (h *VoteItemsHandler) OpenVoteSession(c *gin.Context) {
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

// GET /api/v1/vote_sessions/open: Get open vote session
func (h *VoteItemsHandler) GetOpenVoteSession(c *gin.Context) {
	voteSession, err := h.VoteSessionUseCase.GetOpenVoteSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, voteSession)
}

// PUT /api/v1/vote_sessions/{id}/close: Close a vote session
func (h *VoteItemsHandler) CloseVoteSession(c *gin.Context) {
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

// GET /api/v1/vote_items: Get all active vote items
func (h *VoteItemsHandler) FetchActiveVoteItems(c *gin.Context) {
	voteItems, err := h.VoteItemUseCase.FetchActive(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, voteItems)
}

// POST /api/v1/vote_items: Create a new vote item
func (h *VoteItemsHandler) CreateVoteItem(c *gin.Context) {
	var voteItem domain.VoteItem
	if err := c.ShouldBindJSON(&voteItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": apperror.NewBadRequest(err.Error())})
		return
	}

	err := h.VoteItemUseCase.Create(c.Request.Context(), &voteItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, voteItem)
}

// PUT /api/v1/vote_items/{id}: Update item
func (h *VoteItemsHandler) UpdateVoteItem(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Received id: %v", id)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": apperror.NewBadRequest("ID is required")})
		return
	}

	var voteItem *domain.VoteItem
	if err := c.ShouldBindJSON(&voteItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": apperror.NewBadRequest(err.Error())})
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": apperror.NewBadRequest("Invalid ID format")})
		return
	}
	log.Printf("Received uid: %v", uid)

	voteItem.ID = uid

	log.Printf("Received voteItem: %+v", voteItem)

	ctx := context.Background()
	err = h.VoteItemUseCase.Update(ctx, voteItem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// DELETE /api/v1/vote_items/{id}: Delete a vote item by id
func (h *VoteItemsHandler) DeleteVoteItem(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": apperror.NewBadRequest("ID is required")})
		return
	}

	vid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": apperror.NewBadRequest("Invalid ID format")})
		return
	}

	ctx := context.Background()
	err = h.VoteItemUseCase.Delete(ctx, vid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// DELETE /api/v1/vote_items: Clear all vote items
func (h *VoteItemsHandler) ClearVoteItem(c *gin.Context) {
	err := h.VoteItemUseCase.ClearVoteItem(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": apperror.NewInternal()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// POST /api/v1/votes: Cast a vote
func (h *VoteItemsHandler) CastVote(c *gin.Context) {
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

// GET /api/v1/vote_results/{session_id}: Get vote results by session id
// GET /api/v1/vote_results/{session_id}?format=csv: Get vote results by session id in CSV format
func (h *VoteItemsHandler) GetVoteResultsBySession(c *gin.Context) {
	sessionIDStr := c.Param("session_id")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	voteResults, err := h.VoteUseCase.GetVoteResultsBySession(uint(sessionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	format := c.DefaultQuery("format", "json")
	if format == "csv" {
		// Create a CSV writer writing to the HTTP response
		c.Writer.Header().Set("Content-Type", "text/csv")
		c.Writer.Header().Set("Content-Disposition", "attachment;filename=vote_results.csv")
		buf := &bytes.Buffer{}
		writer := csv.NewWriter(buf)

		// Write the header
		err = writer.Write([]string{"ID", "Description", "Name", "VoteCount", "SessionID", "IsActive"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Write the data
		for _, voteItem := range voteResults {
			record := []string{
				voteItem.VoteItemID.String(),
				voteItem.VoteItemName,
				strconv.Itoa(int(voteItem.VoteCount)),
			}
			err = writer.Write(record)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		writer.Flush()
		if err = writer.Error(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV data"})
			return
		}

		// Set the necessary headers to instruct the browser to download the file
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename=vote_results.csv")
		c.Header("Content-Type", "application/octet-stream")

		// Write the buffer contents to the response
		c.String(http.StatusOK, buf.String())

	} else {
		c.JSON(http.StatusOK, voteResults)
	}
}
