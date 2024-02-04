package usecase

import "github.com/krittawatcode/vote-items/backend-service/domain"

type voteResultUsecase struct {
	voteResultRepo domain.VoteResultRepository
}

func NewVoteResultUsecase(v domain.VoteResultRepository) domain.VoteResultUseCase {
	return &voteResultUsecase{
		voteResultRepo: v,
	}
}

func (u *voteResultUsecase) GetVoteResultsBySession(sessionID uint) ([]domain.VoteResult, error) {
	return u.voteResultRepo.GetVoteResultsBySession(sessionID)
}
