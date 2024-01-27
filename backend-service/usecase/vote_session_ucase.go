package usecase

import (
	"context"
	"errors"

	"github.com/krittawatcode/vote-items/backend-service/domain"
)

type VoteSessionUsecase struct {
	VoteSessionRepository domain.VoteSessionRepository
}

func NewVoteSessionUsecase(r domain.VoteSessionRepository) domain.VoteSessionUseCase {
	return &VoteSessionUsecase{
		VoteSessionRepository: r,
	}
}

func (u *VoteSessionUsecase) GetOpenVoteSession() (*domain.VoteSession, error) {
	voteSession, err := u.VoteSessionRepository.GetOpenVoteSession()
	if err != nil {
		return voteSession, err
	}

	return voteSession, nil
}

func (u *VoteSessionUsecase) OpenVoteSession(id uint) error {
	if _, err := u.VoteSessionRepository.GetOpenVoteSession(); err != nil {
		return err
	}

	err := u.VoteSessionRepository.CreateVoteSession(id)
	if err != nil {
		return err
	}
	return nil
}

func (u *VoteSessionUsecase) CloseVoteSession(id uint) error {
	err := u.VoteSessionRepository.CloseVoteSession(id)
	if err != nil {
		return errors.New("failed to close vote session: " + err.Error())
	}
	return nil
}

func (u *VoteSessionUsecase) GetVoteSessionByID(ctx context.Context, id uint) (*domain.VoteSession, error) {
	voteSession, err := u.VoteSessionRepository.GetVoteSessionByID(ctx, id)
	if err != nil {
		return voteSession, err
	}

	return voteSession, nil
}
