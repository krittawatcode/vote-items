package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/backend-service/delivery/handler/middleware"
	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/apperror"
)

// Handler struct holds required services for handler to function
type VoteItemsHandler struct {
	Router          *gin.Engine
	VoteItemUseCase domain.VoteItemUseCase
	TokenUseCase    domain.TokenUseCase
	BaseUrl         string // base url for user routes
	TimeoutDuration time.Duration
}

// Does not return as it deals directly with a reference to the gin Engine
func NewVoteItemsHandler(router *gin.Engine, viu domain.VoteItemUseCase, tu domain.TokenUseCase, url string, timeout time.Duration) {
	h := &VoteItemsHandler{
		VoteItemUseCase: viu,
		TokenUseCase:    tu,
	}

	// Create an vote-items group
	g := router.Group(url)

	if gin.Mode() != gin.TestMode {
		g.Use(middleware.Timeout(timeout, apperror.NewServiceUnavailable()))
	}

	if gin.Mode() != gin.TestMode {
		// set up middle ware for time out
		g.Use(middleware.Timeout(timeout, apperror.NewServiceUnavailable()))
		// get all active vote items
		g.GET("/", middleware.AuthUser(h.TokenUseCase), h.FetchActiveVoteItems)
		// create a new vote item
		g.POST("/", middleware.AuthUser(h.TokenUseCase), h.CreateVoteItem)
		// update a vote item
		g.PUT("/:id", middleware.AuthUser(h.TokenUseCase), h.UpdateVoteItem)
		// delete a vote item
		g.DELETE("/:id", middleware.AuthUser(h.TokenUseCase), h.DeleteVoteItem)
		// clear all vote items
		g.DELETE("/", middleware.AuthUser(h.TokenUseCase), h.ClearVoteItem)
	}
}

// @Summary Get all active vote items
// @Description Retrieve all active vote items
// @Tags vote_items
// @Produce  json
// @Success 200 {array} domain.VoteItem "Successfully retrieved the active vote items"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /api/v1/vote_items [get]
// GET /api/v1/vote_items: Get all active vote items
func (h *VoteItemsHandler) FetchActiveVoteItems(c *gin.Context) {
	voteItems, err := h.VoteItemUseCase.FetchActive(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, voteItems)
}

// @Summary Create a new vote item
// @Description Create a new vote item with the provided fields
// @Tags vote_items
// @Accept  json
// @Produce  json
// @Param   voteItem     body    domain.VoteItem     true    "Vote Item"
// @Success 201 {object} domain.VoteItem "Successfully created the vote item"
// @Failure 400 {object} domain.ErrorResponse "Bad Request"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /api/v1/vote_items [post]
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

// @Summary Update a vote item
// @Description Update a vote item by ID
// @Tags vote_items
// @Accept  json
// @Produce  json
// @Param   id     path    string     true    "Vote Item ID"
// @Param   voteItem     body    domain.VoteItem     true    "Vote Item"
// @Success 200 {object} domain.SuccessResponse "Vote item updated successfully"
// @Failure 400 {object} domain.ErrorResponse "Bad Request"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /api/v1/vote_items/{id} [put]
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

// @Summary Delete a vote item by id
// @Description Delete a vote item by id
// @Tags vote_items
// @Accept  json
// @Produce  json
// @Param id path string true "Vote Item ID"
// @Success 200 {object} domain.SuccessResponse "Vote item deleted successfully"
// @Failure 400 {object} domain.ErrorResponse "Bad Request"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /api/v1/vote_items/{id} [delete]
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

// @Summary Clear all vote items
// @Description Clear all vote items
// @Tags vote_items
// @Accept  json
// @Produce  json
// @Success 200 {object} domain.SuccessResponse "Vote item cleared successfully"
// @Failure 400 {object} domain.ErrorResponse "Bad Request"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /api/v1/vote_items [delete]
// DELETE /api/v1/vote_items: Clear all vote items
func (h *VoteItemsHandler) ClearVoteItem(c *gin.Context) {
	err := h.VoteItemUseCase.ClearVoteItem(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
