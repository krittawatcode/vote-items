package usecase

import (
	"context"

	"github.com/krittawatcode/vote-items/backend-service/domain"
)

type voteUsecase struct {
	voteRepo domain.VoteRepository
}

func NewVoteUsecase(v domain.VoteRepository) domain.VoteUseCase {
	return &voteUsecase{
		voteRepo: v,
	}
}

func (u *voteUsecase) Create(ctx context.Context, v *domain.Vote) error {
	err := u.voteRepo.Create(ctx, v)
	if err != nil {
		return err
	}
	return nil
}
