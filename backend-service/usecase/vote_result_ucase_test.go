package usecase

import (
	"testing"

	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/appmock"
	"github.com/stretchr/testify/assert"
)

func TestVoteResultUsecase(t *testing.T) {
	mockVoteResultRepo := new(appmock.MockVoteResultRepository)
	mockVoteResultUsecase := NewVoteResultUsecase(mockVoteResultRepo)

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

		mockVoteResultRepo.On("GetVoteResultsBySession", sessionID).Return(mockVoteResults, nil)

		voteResults, err := mockVoteResultUsecase.GetVoteResultsBySession(sessionID)

		assert.NoError(t, err)
		assert.Equal(t, mockVoteResults, voteResults)
		mockVoteResultRepo.AssertExpectations(t)
	})
}
