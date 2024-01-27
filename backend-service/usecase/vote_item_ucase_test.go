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

func TestVoteItemUsecase(t *testing.T) {
	t.Run("FetchActive", func(t *testing.T) {
		mockRepo := new(appmock.MockVoteItemRepository)
		voteItemUsecase := NewVoteItemUsecase(mockRepo)
		mockVoteItems := &[]domain.VoteItem{
			{
				ID: uuid.New(),
			},
		}

		mockRepo.On("FetchActive", mock.Anything).Return(mockVoteItems, nil)

		voteItems, err := voteItemUsecase.FetchActive(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, voteItems)
		assert.Equal(t, voteItems, mockVoteItems)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Create", func(t *testing.T) {
		mockRepo := new(appmock.MockVoteItemRepository)
		voteItemUsecase := NewVoteItemUsecase(mockRepo)
		mockVoteItem := &domain.VoteItem{
			ID: uuid.New(),
		}

		mockRepo.On("Create", mock.Anything, mockVoteItem).Return(nil)

		err := voteItemUsecase.Create(context.Background(), mockVoteItem)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update", func(t *testing.T) {
		mockRepo := new(appmock.MockVoteItemRepository)
		voteItemUsecase := NewVoteItemUsecase(mockRepo)
		mockVoteItem := &domain.VoteItem{
			ID: uuid.New(),
		}

		mockRepo.On("Update", mock.Anything, mockVoteItem).Return(nil)

		err := voteItemUsecase.Update(context.Background(), mockVoteItem)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Delete", func(t *testing.T) {
		mockRepo := new(appmock.MockVoteItemRepository)
		voteItemUsecase := NewVoteItemUsecase(mockRepo)
		vid := uuid.New()

		mockRepo.On("SetActiveVoteItem", mock.Anything, &domain.VoteItem{ID: vid}, false).Return(nil)

		err := voteItemUsecase.Delete(context.Background(), vid)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ClearVoteItem", func(t *testing.T) {
		mockRepo := new(appmock.MockVoteItemRepository)
		voteItemUsecase := NewVoteItemUsecase(mockRepo)
		mockRepo.On("ClearVoteItem", mock.Anything).Return(nil)

		err := voteItemUsecase.ClearVoteItem(context.Background())

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
