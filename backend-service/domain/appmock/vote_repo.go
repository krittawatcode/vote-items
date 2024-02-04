package appmock

import (
	"context"

	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/stretchr/testify/mock"
)

// MockVoteRepository is a mock type for domain.VoteRepository
type MockVoteRepository struct {
	mock.Mock
}

// Create mocks concrete Create
func (m *MockVoteRepository) Create(ctx context.Context, v *domain.Vote) error {
	ret := m.Called(ctx, v)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
