package handler

import (
	"bytes"
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/krittawatcode/vote-items/backend-service/delivery/handler/middleware"
	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/apperror"
)

// Handler struct holds required services for handler to function
type VoteResultsHandler struct {
	Router            *gin.Engine
	VoteResultUseCase domain.VoteResultUseCase
	TokenUseCase      domain.TokenUseCase
	Url               string // base url for vote session routes
	TimeoutDuration   time.Duration
}

// Does not return as it deals directly with a reference to the gin Engine
func NewVoteResultsHandler(router *gin.Engine, vru domain.VoteResultUseCase, tu domain.TokenUseCase, url string, timeout time.Duration) {
	h := &VoteResultsHandler{
		VoteResultUseCase: vru,
		TokenUseCase:      tu,
	}

	// Create an vote-sessions group
	g := router.Group(url)

	if gin.Mode() != gin.TestMode {
		g.Use(middleware.Timeout(timeout, apperror.NewServiceUnavailable()))
	}

	if gin.Mode() != gin.TestMode {
		// set up middle ware for time out
		g.Use(middleware.Timeout(timeout, apperror.NewServiceUnavailable()))
		// get vote results by session id
		// GET /api/v1/vote_results/{session_id}: Get vote results by session id
		// GET /api/v1/vote_results/{session_id}?format=csv: Get vote results by session id in CSV format
		g.GET("/:session_id", middleware.AuthUser(h.TokenUseCase), h.GetVoteResultsBySession)
	}
}

// @Summary Get vote results by session id
// @Description Get vote results by session id. Can also return results in CSV format.
// @Tags vote_results
// @Accept  json
// @Produce  json
// @Produce  text/csv
// @Param session_id path int true "Session ID"
// @Param format query string false "Format of the response (json or csv)"
// @Success 200 {array} domain.VoteResult "Vote results successfully retrieved"
// @Failure 400 {object} domain.ErrorResponse "Bad Request"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /api/v1/vote_results/{session_id} [get]
// GET /api/v1/vote_results/{session_id}: Get vote results by session id
// GET /api/v1/vote_results/{session_id}?format=csv: Get vote results by session id in CSV format
func (h *VoteResultsHandler) GetVoteResultsBySession(c *gin.Context) {
	sessionIDStr := c.Param("session_id")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	voteResults, err := h.VoteResultUseCase.GetVoteResultsBySession(uint(sessionID))
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
