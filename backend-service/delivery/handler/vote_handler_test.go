package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/appmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVotesHandler_CastVote(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userId := uuid.New()
	voteItemId := uuid.New()

	t.Run("Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		vote := domain.Vote{UserID: userId, VoteItemID: voteItemId}
		jsonVote, _ := json.Marshal(vote)
		reqBody := bytes.NewReader(jsonVote)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/votes", reqBody)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user", &domain.User{UID: userId})

		mockVoteUseCase := new(appmock.MockVoteUseCase)
		mockVoteUseCase.On("Create", mock.Anything, mock.Anything).Return(nil)

		h := &VotesHandler{
			VoteUseCase: mockVoteUseCase,
		}
		h.CastVote(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("VoteUseCase.Create returns error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		vote := domain.Vote{UserID: userId, VoteItemID: voteItemId}
		jsonVote, _ := json.Marshal(vote)
		reqBody := bytes.NewReader(jsonVote)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/votes", reqBody)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user", &domain.User{UID: userId})

		mockVoteUseCase := new(appmock.MockVoteUseCase)
		mockVoteUseCase.On("Create", mock.Anything, mock.Anything).Return(errors.New("error"))

		h := &VotesHandler{
			VoteUseCase: mockVoteUseCase,
		}
		h.CastVote(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := bytes.NewReader([]byte(`{"invalid":"json"`)) // invalid json
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/votes", reqBody)
		c.Request.Header.Set("Content-Type", "application/json")

		h := &VotesHandler{}
		h.CastVote(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("User not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		vote := domain.Vote{UserID: userId, VoteItemID: voteItemId}
		jsonVote, _ := json.Marshal(vote)
		reqBody := bytes.NewReader(jsonVote)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/votes", reqBody)
		c.Request.Header.Set("Content-Type", "application/json")

		h := &VotesHandler{}
		h.CastVote(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
