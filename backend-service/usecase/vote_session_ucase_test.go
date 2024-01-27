package usecase

import (
	"context"
	"testing"

	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/appmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVoteSessionUsecase(t *testing.T) {
	mockVoteSessionRepo := new(appmock.MockVoteSessionRepository)
	mockVoteSessionUsecase := NewVoteSessionUsecase(mockVoteSessionRepo)

	t.Run("GetOpenVoteSession", func(t *testing.T) {
		mockVoteSession := &domain.VoteSession{
			ID: 1,
		}

		mockVoteSessionRepo.On("GetOpenVoteSession").Return(mockVoteSession, nil)

		voteSession, err := mockVoteSessionUsecase.GetOpenVoteSession()

		assert.NoError(t, err)
		assert.NotNil(t, voteSession)
		assert.Equal(t, voteSession, mockVoteSession)
		mockVoteSessionRepo.AssertExpectations(t)
	})

	t.Run("OpenVoteSession", func(t *testing.T) {
		id := uint(1)

		mockVoteSessionRepo.On("GetOpenVoteSession").Return(nil, nil)
		mockVoteSessionRepo.On("CreateVoteSession", id).Return(nil)

		err := mockVoteSessionUsecase.OpenVoteSession(id)

		assert.NoError(t, err)
		mockVoteSessionRepo.AssertExpectations(t)
	})

	t.Run("CloseVoteSession", func(t *testing.T) {
		id := uint(1)

		mockVoteSessionRepo.On("CloseVoteSession", id).Return(nil)

		err := mockVoteSessionUsecase.CloseVoteSession(id)

		assert.NoError(t, err)
		mockVoteSessionRepo.AssertExpectations(t)
	})

	t.Run("GetVoteSessionByID", func(t *testing.T) {
		id := uint(1)
		mockVoteSession := &domain.VoteSession{
			ID: id,
		}

		mockVoteSessionRepo.On("GetVoteSessionByID", mock.Anything, id).Return(mockVoteSession, nil)

		voteSession, err := mockVoteSessionUsecase.GetVoteSessionByID(context.Background(), id)

		assert.NoError(t, err)
		assert.NotNil(t, voteSession)
		assert.Equal(t, voteSession, mockVoteSession)
		mockVoteSessionRepo.AssertExpectations(t)
	})
}
