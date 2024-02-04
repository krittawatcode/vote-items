package repository

import (
	"errors"
	"log"

	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/apperror"
	"gorm.io/gorm"
)

type gormVoteResultRepository struct {
	conn *gorm.DB
}

func NewGormVoteResultRepository(conn *gorm.DB) domain.VoteResultRepository {
	return &gormVoteResultRepository{conn}
}

func (r *gormVoteResultRepository) GetVoteResultsBySession(sessionID uint) ([]domain.VoteResult, error) {
	var results []domain.VoteResult

	// Join votes and vote_items tables, filter by session ID, group by vote_item_id and vote_items.name, and order by vote count
	err := r.conn.Table("votes").
		Select("vote_items.id as vote_item_id, vote_items.name as vote_item_name, COUNT(votes.id) as vote_count").
		Joins("JOIN vote_items ON votes.vote_item_id = vote_items.id").
		Where("votes.session_id = ?", sessionID).
		Group("vote_items.id, vote_items.name").
		Order("vote_count DESC").
		Scan(&results).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Log the event and return a not found error if no vote results are found
			log.Printf("No vote results found for session ID: %v\n", sessionID)
			return nil, apperror.NewNotFound("vote results", "")
		}
		// Return the error if any other error occurs during the execution of the query
		log.Printf("Error retrieving vote results: %v\n", err)
		return nil, apperror.NewInternal()
	}

	return results, nil
}
