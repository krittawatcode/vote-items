package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/backend/domain/model"

	"github.com/stretchr/testify/mock"
)

// MockUserUseCase is a mock type for model.UserUseCase
type MockUserUseCase struct {
	mock.Mock
}

// Get is mock of UserUseCase Get
func (m *MockUserUseCase) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	// args that will be passed to "Return" in the tests, when function
	// is called with a uid. Hence the name "ret"
	ret := m.Called(ctx, uid)

	// first value passed to "Return"
	var r0 *model.User
	if ret.Get(0) != nil {
		// we can just return this if we know we won't be passing function to "Return"
		r0 = ret.Get(0).(*model.User)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
