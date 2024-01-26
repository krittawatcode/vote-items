package repository

import (
	"context"

	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/apperror"
	"gorm.io/gorm"
)

type gormVoteItemRepository struct {
	conn *gorm.DB
}

func NewGormVoteItemRepository(conn *gorm.DB) domain.VoteItemRepository {
	return &gormVoteItemRepository{conn}
}

func (r *gormVoteItemRepository) FetchActive(ctx context.Context) (*[]domain.VoteItem, error) {
	var voteItems []domain.VoteItem
	if err := r.conn.Where("is_active = ?", true).Find(&voteItems).Error; err != nil {
		return nil, apperror.NewInternal()
	}
	return &voteItems, nil
}

func (r *gormVoteItemRepository) Create(ctx context.Context, v *domain.VoteItem) error {
	result := r.conn.Create(v)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *gormVoteItemRepository) Update(ctx context.Context, v *domain.VoteItem) error {
	var currentVoteItem domain.VoteItem
	if err := r.conn.First(&currentVoteItem, v.ID).Error; err != nil || currentVoteItem.VoteCount != 0 {
		return apperror.NewConflict("Cannot update vote item: Vote count is not zero or item not found", "")
	}
	return r.conn.Save(v).Error
}

func (r *gormVoteItemRepository) SetActiveVoteItem(ctx context.Context, v *domain.VoteItem, isActive bool) error {
	var currentVoteItem domain.VoteItem
	if err := r.conn.First(&currentVoteItem, v.ID).Error; err != nil {
		return apperror.NewNotFound("Vote item not found", "")
	}
	if currentVoteItem.VoteCount != 0 {
		return apperror.NewConflict("Cannot set active vote item: Vote count is not zero", "")
	}
	currentVoteItem.IsActive = isActive
	return r.conn.Save(&currentVoteItem).Error
}

// will set all voteItem to inactive
func (r *gormVoteItemRepository) ClearVoteItem(ctx context.Context) error {
	return r.conn.Model(&domain.VoteItem{}).Where("is_active = ?", true).Update("is_active", false).Error
}
