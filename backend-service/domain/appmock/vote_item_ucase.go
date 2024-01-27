package appmock

import (
	"context"

	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/stretchr/testify/mock"
)

type MockVoteItemUseCase struct {
	mock.Mock
}

func (m *MockVoteItemUseCase) FetchActive(ctx context.Context) (*[]domain.VoteItem, error) {
	ret := m.Called(ctx)
	var r0 *[]domain.VoteItem
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*[]domain.VoteItem)
	}
	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}
	return r0, r1
}

func (m *MockVoteItemUseCase) Create(ctx context.Context, v *domain.VoteItem) error {
	ret := m.Called(ctx, v)
	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}
	return r0
}

func (m *MockVoteItemUseCase) Update(ctx context.Context, v *domain.VoteItem) error {
	ret := m.Called(ctx, v)
	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}
	return r0
}

func (m *MockVoteItemUseCase) Delete(ctx context.Context, vid uuid.UUID) error {
	ret := m.Called(ctx, vid)
	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}
	return r0
}

func (m *MockVoteItemUseCase) ClearVoteItem(ctx context.Context) error {
	ret := m.Called(ctx)
	var r0 error
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(error)
	}
	return r0
}
