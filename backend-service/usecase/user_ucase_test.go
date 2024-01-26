package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/apperror"
	"github.com/krittawatcode/vote-items/backend-service/domain/appmock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {
	mockRepo := new(appmock.MockUserRepository)
	userUseCase := &userUseCase{
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

func TestSignUp(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &domain.User{
			Email:    "bob@bob.com",
			Password: "howdyhoneighbor!",
		}

		mockUserRepository := new(appmock.MockUserRepository)
		userService := NewUserUseCase(mockUserRepository)

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockUserRepository.
			On("Create", mock.Anything, mockUser).
			Run(func(args mock.Arguments) {
				userArg := args.Get(1).(*domain.User) // arg 0 is context, arg 1 is *User
				userArg.UID = uid
			}).Return(nil)

		ctx := context.TODO()
		err := userService.SignUp(ctx, mockUser)

		assert.NoError(t, err)

		// assert user now has a userID
		assert.Equal(t, uid, mockUser.UID)

		mockUserRepository.AssertExpectations(t)
	})
	t.Run("Error", func(t *testing.T) {
		mockUser := &domain.User{
			Email:    "bob@bob.com",
			Password: "howdyhoneighbor!",
		}

		mockUserRepository := new(appmock.MockUserRepository)
		userService := NewUserUseCase(mockUserRepository)

		mockErr := apperror.NewConflict("email", mockUser.Email)

		// We can use Run method to modify the user when the Create method is called.
		//  We can then chain on a Return method to return no error
		mockUserRepository.
			On("Create", mock.Anything, mockUser).
			Return(mockErr)

		ctx := context.TODO()
		err := userService.SignUp(ctx, mockUser)

		// assert error is error we response with in mock
		assert.EqualError(t, err, mockErr.Error())

		mockUserRepository.AssertExpectations(t)
	})
}
