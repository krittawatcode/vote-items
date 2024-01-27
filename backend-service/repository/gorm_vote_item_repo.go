package repository

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
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
	// check if current session is open or not
	var voteSession domain.VoteSession
	if err := r.conn.Where("is_open = ?", true).First(&voteSession).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No open vote session found: %v\n", err)
			return apperror.NewNotFound("open vote session", "")
		}
		log.Printf("Error finding open vote session: %v\n", err)
		return apperror.NewInternal()
	}

	v.SessionID = voteSession.ID
	log.Printf("Create vote item with data: %v\n", v)
	result := r.conn.Create(v)
	if result.Error != nil {
		// Log the error
		log.Printf("Could not create a vote item with ID: %v. Reason: %v\n", v.ID, result.Error)

		// Check for unique constraint violation
		var pgErr *pgconn.PgError
		if ok := errors.As(result.Error, &pgErr); ok {
			log.Printf("Could not create a vote item with ID: %v. Reason: %v\n", v.ID, pgErr.Hint)
			return apperror.NewConflict("id", v.ID.String())
		}

		// If the error is not a unique constraint violation, return a generic internal server error
		log.Printf("Could not create a vote item with ID and cannot assert err to pg.Error: %v. Reason: %v\n", v.ID, result.Error)
		return apperror.NewInternal()
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
	result := r.conn.Model(&domain.VoteItem{}).Where("is_active = ?", true).Update("is_active", false)
	if result.Error != nil {
		// Log the error
		log.Printf("Could not clear vote items. Reason: %v\n", result.Error)

		// Check for unique constraint violation
		var pgErr *pgconn.PgError
		if ok := errors.As(result.Error, &pgErr); ok {
			log.Printf("Could not clear vote items. Reason: %v\n", pgErr.Hint)
			return apperror.NewConflict("id", "ClearVoteItem conflict")
		}

		// If the error is not a unique constraint violation, return a generic internal server error
		log.Printf("Could not clear vote items and cannot assert err to pg.Error: %v. Reason: %v\n", "ClearVoteItem", result.Error)
		return apperror.NewInternal()
	}
	return nil
}
