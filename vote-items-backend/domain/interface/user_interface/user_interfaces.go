package user_interface

import (
	"context"

	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/backend/domain/model"
)

// UserUseCase defines methods the handler layer expects
// any service it interacts with to implement
type UserUseCase interface {
	Get(ctx context.Context, uid uuid.UUID) (*model.User, error)
}

// UserRepository defines methods the service layer expects
// any repository it interacts with to implement
type UserRepository interface {
	FindByID(ctx context.Context, uid uuid.UUID) (*model.User, error)
}
