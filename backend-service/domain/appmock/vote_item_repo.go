package appmock

import (
	"context"

	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/stretchr/testify/mock"
)

// MockVoteItemRepository is a mock type for domain.VoteItemRepository
type MockVoteItemRepository struct {
	mock.Mock
}

// FetchActive mocks concrete FetchActive
func (m *MockVoteItemRepository) FetchActive(ctx context.Context) (*[]domain.VoteItem, error) {
	ret := m.Called(ctx)

	var r0 *[]domain.VoteItem
	if ret.Get(0) != nil {
		// Ensure you're asserting to the correct type (*[]domain.VoteItem)
		r0 = ret.Get(0).(*[]domain.VoteItem)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1
}

// Create mocks concrete Create
func (m *MockVoteItemRepository) Create(ctx context.Context, v *domain.VoteItem) error {
	ret := m.Called(ctx, v)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// Update mocks concrete Update
func (m *MockVoteItemRepository) Update(ctx context.Context, v *domain.VoteItem) error {
	ret := m.Called(ctx, v)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// SetActiveVoteItem mocks concrete SetActiveVoteItem
func (m *MockVoteItemRepository) SetActiveVoteItem(ctx context.Context, v *domain.VoteItem, isActive bool) error {
	ret := m.Called(ctx, v, isActive)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}

// ClearVoteItem mocks concrete ClearVoteItem
func (m *MockVoteItemRepository) ClearVoteItem(ctx context.Context) error {
	ret := m.Called(ctx)

	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}

	return r0
}
