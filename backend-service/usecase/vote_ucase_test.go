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
}
