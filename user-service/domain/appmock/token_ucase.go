package appmock

import (
	"context"

	"github.com/krittawatcode/vote-items/user-service/domain"
	"github.com/stretchr/testify/mock"
)

// MockTokenUseCase is a mock type for domain.TokenUseCase
type MockTokenUseCase struct {
	mock.Mock
}

// NewPairFromUser mocks concrete NewPairFromUser
func (m *MockTokenUseCase) NewPairFromUser(ctx context.Context, u *domain.User, prevTokenID string) (*domain.TokenPair, error) {
	ret := m.Called(ctx, u, prevTokenID)

	// first value passed to "Return"
	var r0 *domain.TokenPair
	if ret.Get(0) != nil {
		// we can just return this if we know we won't be passing function to "Return"
		r0 = ret.Get(0).(*domain.TokenPair)
	}

	var r1 error

	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}