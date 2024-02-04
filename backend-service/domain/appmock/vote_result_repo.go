package appmock

import (
	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/stretchr/testify/mock"
)

// MockGormVoteResultRepository is a mock type for domain.VoteResultRepository
type MockVoteResultRepository struct {
	mock.Mock
}

// GetVoteResultsBySession mocks concrete GetVoteResultsBySession
func (m *MockVoteResultRepository) GetVoteResultsBySession(sessionID uint) ([]domain.VoteResult, error) {
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
