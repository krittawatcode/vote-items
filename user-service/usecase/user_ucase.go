package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/user-service/domain"
)

// UserUseCase acts as a struct for injecting an implementation of UserRepository
// for use in service methods
type UserUseCase struct {
	UserRepository domain.UserRepository
}

// USConfig will hold repositories that will eventually be injected into this
// this service layer
type USConfig struct {
	UserRepository domain.UserRepository
}

// NewUserUseCase is a factory function for
// initializing a UserService with its repository layer dependencies
func NewUserUseCase(c *USConfig) domain.UserUseCase {
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
	panic("Method not implemented")
}
