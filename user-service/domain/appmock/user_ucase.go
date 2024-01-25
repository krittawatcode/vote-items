package appmock

import (
	"context"

	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/user-service/domain"

	"github.com/stretchr/testify/mock"
)

// MockUserUseCase is a mock type for domain.UserUseCase
type MockUserUseCase struct {
	mock.Mock
}

// Get is mock of UserUseCase Get
func (m *MockUserUseCase) Get(ctx context.Context, uid uuid.UUID) (*domain.User, error) {
	// args that will be passed to "Return" in the tests, when function
	// is called with a uid. Hence the name "ret"
	ret := m.Called(ctx, uid)

	// first value passed to "Return"
	var r0 *domain.User
	if ret.Get(0) != nil {
		// we can just return this if we know we won't be passing function to "Return"
		r0 = ret.Get(0).(*domain.User)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

// Sign up is a mock of UserUseCase.SignUp
func (m *MockUserUseCase) SignUp(ctx context.Context, u *domain.User) error {
	ret := m.Called(ctx, u)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// SignIn is a mock of UserService.Signin
func (m *MockUserUseCase) SignIn(ctx context.Context, u *domain.User) error {
	ret := m.Called(ctx, u)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
