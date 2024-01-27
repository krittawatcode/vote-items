package repository

import (
	"context"
	"errors"
	"log"

	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/apperror"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type gormVoteRepository struct {
	Conn *gorm.DB
}

func NewGormVoteRepository(conn *gorm.DB) domain.VoteRepository {
	return &gormVoteRepository{conn}
}

// Create is a method that creates a new vote in the database.
func (r *gormVoteRepository) Create(ctx context.Context, v *domain.Vote) error {
	// check if current session is open or not
	var voteSession domain.VoteSession
	if err := r.Conn.Where("is_open = ?", true).First(&voteSession).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No open vote session found: %v\n", err)
			return apperror.NewNotFound("open vote session", "")
		}
		log.Printf("Error finding open vote session: %v\n", err)
		return apperror.NewInternal()
	}

	var vote domain.Vote
	if err := r.Conn.Where("user_id = ? AND session_id = ?", v.UserID, v.SessionID).First(&vote).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no record found, create a new vote
			if err := r.Conn.Create(v).Error; err != nil {
				if pgErr, ok := err.(*pq.Error); ok {
					// Handle the postgres error here
					log.Printf("Postgres error creating vote: %v\n", pgErr)
					return apperror.NewConflict(pgErr.Message, pgErr.Hint)
				}
				log.Printf("Error creating vote: %v\n", err)
				return apperror.NewInternal()
			}
			log.Printf("Vote created successfully for user ID: %v and session ID: %v\n", v.UserID, v.SessionID)
			return nil
		}
		// If a record is found, it means the user has already voted in this session
		log.Printf("User with ID: %v has already voted in session ID: %v\n", v.UserID, v.SessionID)
		return apperror.NewConflict("User has already voted in this session", "")
	}
	return nil
}

// GetVoteResultsBySession retrieves the vote results for a given session ID.
// It joins the votes and vote_items tables on the vote_item_id field,
// filters the results by the session ID, groups the results by vote_item_id and vote_items.name,
// and orders the results by the count of votes in descending order.
// If no vote results are found for the session ID, it logs the event and returns a not found error.
// If any other error occurs during the execution of the query, it returns the error.
func (r *gormVoteRepository) GetVoteResultsBySession(sessionID uint) ([]domain.VoteResult, error) {
	var voteResults []domain.VoteResult
	if err := r.Conn.Table("votes").
		Select("votes.*, vote_items.name as vote_item_name, count(votes.vote_item_id) as vote_count").
		Joins("JOIN vote_items ON votes.vote_item_id = vote_items.id").
		Where("votes.session_id = ?", sessionID).
		Group("votes.vote_item_id, vote_items.name").
		Order("count(votes.vote_item_id) desc").
		Find(&voteResults).Error; err != nil {
		// handle error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No vote results found for session ID: %v\n", sessionID)
			return nil, apperror.NewNotFound("vote results", "")
		}
		return nil, err
	}
	return voteResults, nil
}
