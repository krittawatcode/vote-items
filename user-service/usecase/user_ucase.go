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
type userUseCase struct {
	UserRepository domain.UserRepository
}

// NewUserUseCase is a factory function for
// initializing a NewUserUseCase with its usecase layer dependencies
func NewUserUseCase(r domain.UserRepository) domain.UserUseCase {
	return &userUseCase{
		UserRepository: r,
	}
}

// Get retrieves a user based on their uuid
func (s *userUseCase) Get(ctx context.Context, uid uuid.UUID) (*domain.User, error) {
	u, err := s.UserRepository.FindByID(ctx, uid)

	return u, err
}

// Sign up reaches our to a UserRepository to verify the
// email address is available and signs up the user if this is the case
func (s *userUseCase) SignUp(ctx context.Context, u *domain.User) error {
	pw, err := hashPassword(u.Password)

	if err != nil {
		log.Printf("Unable to sign up user for email: %v\n", u.Email)
		return apperror.NewInternal()
	}

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
