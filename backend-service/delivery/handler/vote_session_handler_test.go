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
)

func TestVoteItemsHandler_OpenVoteSession(t *testing.T) {
	// setup
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockVoteSessionUseCase := new(appmock.MockVoteSessionUseCase)
		mockVoteSessionUseCase.On("OpenVoteSession", uint(1)).Return(nil)

		h := &VoteSessionsHandler{
			VoteSessionUseCase: mockVoteSessionUseCase,
		}

		h.OpenVoteSession(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid session ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		h := &VoteSessionsHandler{
			VoteSessionUseCase: new(appmock.MockVoteSessionUseCase),
		}

		h.OpenVoteSession(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error opening vote session", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockVoteSessionUseCase := new(appmock.MockVoteSessionUseCase)
		mockVoteSessionUseCase.On("OpenVoteSession", uint(1)).Return(errors.New("error"))

		h := &VoteSessionsHandler{
			VoteSessionUseCase: mockVoteSessionUseCase,
		}

		h.OpenVoteSession(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestVoteItemsHandler_GetOpenVoteSession(t *testing.T) {
	// setup
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockVoteSession := &domain.VoteSession{}
		mockVoteSessionUseCase := new(appmock.MockVoteSessionUseCase)
		mockVoteSessionUseCase.On("GetOpenVoteSession").Return(mockVoteSession, nil)

		h := &VoteSessionsHandler{
			VoteSessionUseCase: mockVoteSessionUseCase,
		}

		h.GetOpenVoteSession(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Error getting open vote session", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mockVoteSessionUseCase := new(appmock.MockVoteSessionUseCase)
		mockVoteSessionUseCase.On("GetOpenVoteSession").Return(nil, errors.New("error"))

		h := &VoteSessionsHandler{
			VoteSessionUseCase: mockVoteSessionUseCase,
		}

		h.GetOpenVoteSession(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestVoteItemsHandler_CloseVoteSession(t *testing.T) {
	// setup
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockVoteSessionUseCase := new(appmock.MockVoteSessionUseCase)
		mockVoteSessionUseCase.On("CloseVoteSession", uint(1)).Return(nil)

		h := &VoteSessionsHandler{
			VoteSessionUseCase: mockVoteSessionUseCase,
		}

		h.CloseVoteSession(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid session ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		h := &VoteSessionsHandler{
			VoteSessionUseCase: new(appmock.MockVoteSessionUseCase),
		}

		h.CloseVoteSession(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error closing vote session", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockVoteSessionUseCase := new(appmock.MockVoteSessionUseCase)
		mockVoteSessionUseCase.On("CloseVoteSession", uint(1)).Return(errors.New("error"))

		h := &VoteSessionsHandler{
			VoteSessionUseCase: mockVoteSessionUseCase,
		}

		h.CloseVoteSession(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
