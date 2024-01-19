package interfaces

import (
	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/backend/domain/model"
)

// UserService defines methods the handler layer expects
// any service it interacts with to implement
type UserService interface {
	Get(uid uuid.UUID) (*model.User, error)
}

// UserRepository defines methods the service layer expects
// any repository it interacts with to implement
type UserRepository interface {
	FindByID(uid uuid.UUID) (*model.User, error)
}
