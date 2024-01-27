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

// GetVoteResultsBySession mocks concrete GetVoteResultsBySession
func (m *MockVoteRepository) GetVoteResultsBySession(sessionID uint) ([]domain.VoteResult, error) {
	ret := m.Called(sessionID)

	var r0 []domain.VoteResult
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]domain.VoteResult)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
