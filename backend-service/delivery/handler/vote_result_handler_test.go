package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/appmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVoteResultsHandler_GetVoteResultsBySession(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "session_id", Value: "1"}}

		mockVoteResultUseCase := new(appmock.MockVoteResultUsecase)
		mockVoteResultUseCase.On("GetVoteResultsBySession", mock.Anything).Return([]domain.VoteResult{}, nil)

		h := &VoteResultsHandler{
			VoteResultUseCase: mockVoteResultUseCase,
		}
		h.GetVoteResultsBySession(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid session ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "session_id", Value: "invalid"}}

		h := &VoteResultsHandler{}
		h.GetVoteResultsBySession(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("VoteResultUseCase.GetVoteResultsBySession returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "session_id", Value: "1"}}

		mockVoteResultUseCase := new(appmock.MockVoteResultUsecase)
		mockVoteResultUseCase.On("GetVoteResultsBySession", mock.Anything).Return([]domain.VoteResult{}, errors.New("error"))

		h := &VoteResultsHandler{
			VoteResultUseCase: mockVoteResultUseCase,
		}
		h.GetVoteResultsBySession(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
