package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/appmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

		h := &VoteItemsHandler{
			VoteSessionUseCase: mockVoteSessionUseCase,
		}

		h.OpenVoteSession(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid session ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		h := &VoteItemsHandler{
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

		h := &VoteItemsHandler{
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

		h := &VoteItemsHandler{
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

		h := &VoteItemsHandler{
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

		h := &VoteItemsHandler{
			VoteSessionUseCase: mockVoteSessionUseCase,
		}

		h.CloseVoteSession(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid session ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		h := &VoteItemsHandler{
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

		h := &VoteItemsHandler{
			VoteSessionUseCase: mockVoteSessionUseCase,
		}

		h.CloseVoteSession(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
func TestVoteItemsHandler_FetchActiveVoteItems(t *testing.T) {
	// setup
	gin.SetMode(gin.TestMode)
	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Initialize the Request of gin.Context
		c.Request, _ = http.NewRequest("GET", "/api/v1/vote_items", nil)

		mockVoteItems := []domain.VoteItem{
			{ID: uuid.New(), Description: "Vote Item 1", Name: "Vote Item 1", VoteCount: 0, SessionID: 1, IsActive: true},
			{ID: uuid.New(), Description: "Vote Item 2", Name: "Vote Item 2", VoteCount: 0, SessionID: 1, IsActive: true},
			{ID: uuid.New(), Description: "Vote Item 3", Name: "Vote Item 3", VoteCount: 0, SessionID: 1, IsActive: true},
		}
		mockVoteItemUseCase := new(appmock.MockVoteItemUseCase)
		mockVoteItemUseCase.On("FetchActive", mock.Anything).Return(&mockVoteItems, nil)

		h := &VoteItemsHandler{
			VoteItemUseCase: mockVoteItemUseCase,
		}

		h.FetchActiveVoteItems(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestVoteItemsHandler_CreateVoteItem(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := strings.NewReader(`{"name":"Test Vote Item","description":"Test Description","is_active":true}`)
		c.Request = httptest.NewRequest(http.MethodPost, "/voteItems", reqBody)
		c.Request.Header.Set("Content-Type", "application/json")

		mockVoteItemUseCase := new(appmock.MockVoteItemUseCase)
		mockVoteItemUseCase.On("Create", mock.Anything, mock.Anything).Return(nil)

		h := &VoteItemsHandler{
			VoteItemUseCase: mockVoteItemUseCase,
		}
		h.CreateVoteItem(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("VoteItemUseCase.Create returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := strings.NewReader(`{"name":"Test Vote Item","description":"Test Description","is_active":true}`)
		c.Request = httptest.NewRequest(http.MethodPost, "/voteItems", reqBody)
		c.Request.Header.Set("Content-Type", "application/json")

		mockVoteItemUseCase := new(appmock.MockVoteItemUseCase)
		mockVoteItemUseCase.On("Create", mock.Anything, mock.Anything).Return(errors.New("error"))

		h := &VoteItemsHandler{
			VoteItemUseCase: mockVoteItemUseCase,
		}
		h.CreateVoteItem(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := strings.NewReader(`{"name":"","description":"Test Description","is_active":true}`) // missing name
		c.Request = httptest.NewRequest(http.MethodPost, "/voteItems", reqBody)
		c.Request.Header.Set("Content-Type", "application/json")

		h := &VoteItemsHandler{}
		h.CreateVoteItem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
