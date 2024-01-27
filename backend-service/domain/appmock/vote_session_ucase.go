package appmock

import (
	"context"

	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/stretchr/testify/mock"
)

// MockVoteSessionUseCase is a mock type for domain.VoteSessionUseCase
type MockVoteSessionUseCase struct {
	mock.Mock
}

// GetOpenVoteSession mocks concrete GetOpenVoteSession
func (m *MockVoteSessionUseCase) GetOpenVoteSession() (*domain.VoteSession, error) {
	ret := m.Called()

	var r0 *domain.VoteSession
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.VoteSession)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

// OpenVoteSession mocks concrete OpenVoteSession
func (m *MockVoteSessionUseCase) OpenVoteSession(id uint) error {
	ret := m.Called(id)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// CloseVoteSession mocks concrete CloseVoteSession
func (m *MockVoteSessionUseCase) CloseVoteSession(id uint) error {
	ret := m.Called(id)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// GetVoteSessionByID mocks concrete GetVoteSessionByID
func (m *MockVoteSessionUseCase) GetVoteSessionByID(ctx context.Context, id uint) (*domain.VoteSession, error) {
	ret := m.Called(ctx, id)

	var r0 *domain.VoteSession
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*domain.VoteSession)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}
