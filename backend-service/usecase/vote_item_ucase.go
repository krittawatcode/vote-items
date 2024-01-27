package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/backend-service/domain"
)

type voteItemUsecase struct {
	voteItemRepo domain.VoteItemRepository
}

func NewVoteItemUsecase(v domain.VoteItemRepository) domain.VoteItemUseCase {
	return &voteItemUsecase{
		voteItemRepo: v,
	}
}

func (u *voteItemUsecase) FetchActive(ctx context.Context) (*[]domain.VoteItem, error) {
	voteItems, err := u.voteItemRepo.FetchActive(ctx)
	if err != nil {
		return nil, err
	}
	return voteItems, nil
}

func (u *voteItemUsecase) Create(ctx context.Context, v *domain.VoteItem) error {
	err := u.voteItemRepo.Create(ctx, v)
	if err != nil {
		return err
	}
	return nil
}

func (u *voteItemUsecase) Update(ctx context.Context, v *domain.VoteItem) error {
	err := u.voteItemRepo.Update(ctx, v)
	if err != nil {
		return err
	}
	return nil
}

func (u *voteItemUsecase) Delete(ctx context.Context, vid uuid.UUID) error {
	v := &domain.VoteItem{ID: vid}
	err := u.voteItemRepo.SetActiveVoteItem(ctx, v, false)
	if err != nil {
		return err
	}
	return nil
}

// will set all voteItem to inactive
func (u *voteItemUsecase) ClearVoteItem(ctx context.Context) error {
	err := u.voteItemRepo.ClearVoteItem(ctx)
	if err != nil {
		return err
	}
	return nil
}
