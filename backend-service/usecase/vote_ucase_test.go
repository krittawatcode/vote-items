package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/appmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVoteUsecase(t *testing.T) {
	mockVoteRepo := new(appmock.MockVoteRepository)
	mockVoteUsecase := NewVoteUsecase(mockVoteRepo)

	t.Run("Create", func(t *testing.T) {
		mockVote := &domain.Vote{
			SessionID: uint(uuid.New().ID()),
			UserID:    uuid.New(),
		}

		mockVoteRepo.On("Create", mock.Anything, mockVote).Return(nil)

		err := mockVoteUsecase.Create(context.Background(), mockVote)

		assert.NoError(t, err)
		mockVoteRepo.AssertExpectations(t)
	})

	t.Run("GetVoteResultsBySession", func(t *testing.T) {
		sessionID := uint(1)

		mockVoteResults := []domain.VoteResult{
			{
				VoteItemID: uuid.New(),
				VoteCount:  10,
			},
			{
				VoteItemID: uuid.New(),
				VoteCount:  5,
			},
		}

		mockVoteRepo.On("GetVoteResultsBySession", sessionID).Return(mockVoteResults, nil)

		voteResults, err := mockVoteUsecase.GetVoteResultsBySession(sessionID)

		assert.NoError(t, err)
		assert.Equal(t, mockVoteResults, voteResults)
		mockVoteRepo.AssertExpectations(t)
	})
}
