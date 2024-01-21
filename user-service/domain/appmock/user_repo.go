package appmock

import (
	"context"

	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/user-service/domain"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock type for domain.UserRepository
type MockUserRepository struct {
	mock.Mock
}

// FindByID is mock of UserRepository FindByID
func (m *MockUserRepository) FindByID(ctx context.Context, uid uuid.UUID) (*domain.User, error) {
	ret := m.Called(ctx, uid)

	var r0 *domain.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.User)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
