package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/user-service/domain"
	"github.com/krittawatcode/vote-items/user-service/domain/appmock"

	"github.com/krittawatcode/vote-items/user-service/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {
	mockRepo := new(appmock.MockUserRepository)
	userUseCase := &usecase.UserUseCase{
		UserRepository: mockRepo,
	}

	t.Run("success", func(t *testing.T) {
		uid := uuid.New()

		mockUserResp := &domain.User{
			UID:   uid,
			Email: "bob@bob.com",
		}

		mockRepo.On("FindByID", mock.Anything, uid).Return(mockUserResp, nil)

		u, err := userUseCase.Get(context.Background(), uid)

		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, u, mockUserResp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		uid := uuid.New()
		mockRepo.On("FindByID", mock.Anything, uid).Return(nil, fmt.Errorf("Some error down the call chain"))

		u, err := userUseCase.Get(context.Background(), uid)

		assert.Nil(t, u)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}