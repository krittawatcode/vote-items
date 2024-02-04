package appmock

import (
	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/stretchr/testify/mock"
)

type MockVoteResultUsecase struct {
	mock.Mock
}

func (m *MockVoteResultUsecase) GetVoteResultsBySession(sessionID uint) ([]domain.VoteResult, error) {
	args := m.Called(sessionID)
	return args.Get(0).([]domain.VoteResult), args.Error(1)
}
