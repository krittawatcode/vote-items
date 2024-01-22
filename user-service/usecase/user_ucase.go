package usecase

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/user-service/domain"
	"github.com/krittawatcode/vote-items/user-service/domain/apperror"
)

// UserUseCase acts as a struct for injecting an implementation of UserRepository
// for use in service methods
type UserUseCase struct {
	UserRepository domain.UserRepository
}

// UUConfig will hold repositories that will eventually be injected into this
// this service layer
type UUConfig struct {
	UserRepository domain.UserRepository
}

// NewUserUseCase is a factory function for
// initializing a UserService with its repository layer dependencies
func NewUserUseCase(c *UUConfig) domain.UserUseCase {
	return &UserUseCase{
		UserRepository: c.UserRepository,
	}
}

// Get retrieves a user based on their uuid
func (s *UserUseCase) Get(ctx context.Context, uid uuid.UUID) (*domain.User, error) {
	u, err := s.UserRepository.FindByID(ctx, uid)

	return u, err
}

// Sign up reaches our to a UserRepository to verify the
// email address is available and signs up the user if this is the case
func (s *UserUseCase) SignUp(ctx context.Context, u *domain.User) error {
	pw, err := hashPassword(u.Password)

	if err != nil {
		log.Printf("Unable to sign up user for email: %v\n", u.Email)
		return apperror.NewInternal()
	}

	// now I realize why I originally used SignUp(ctx, email, password)
	// then created a user. It's somewhat un-natural to mutate the user here
	u.Password = pw

	err = s.UserRepository.Create(ctx, u)
	if err != nil {
		return err
	}

	// If we get around to adding events, we'd Publish it here
	// err := s.EventsBroker.PublishUserUpdated(u, true)

	// if err != nil {
	// 	return nil, apperrors.NewInternal()
	// }

	return nil
}
