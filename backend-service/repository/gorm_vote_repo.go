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
	conn *gorm.DB
}

func NewGormVoteRepository(conn *gorm.DB) domain.VoteRepository {
	return &gormVoteRepository{conn}
}

// Create is a method that creates a new vote in the database.
func (r *gormVoteRepository) Create(ctx context.Context, v *domain.Vote) error {
	// log request data
	log.Printf("Creating vote : %v\n", v)
	// check if current session is open or not
	var voteSession domain.VoteSession
	if err := r.conn.Where("is_open = ?", true).First(&voteSession).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No open vote session found: %v\n", err)
			return apperror.NewNotFound("vote session", "OPEN")
		}
		log.Printf("Error finding open vote session: %v\n", err)
		return apperror.NewInternal()
	}

	v.SessionID = voteSession.ID

	// Check if the user has already voted for an item in this session
	var existingVote domain.Vote
	if err := r.conn.Where("user_id = ? AND session_id = ?", v.UserID, voteSession.ID).First(&existingVote).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("User with ID: %v has already voted for an item in session ID: %v\n", v.UserID, v.SessionID)
		return apperror.NewConflict("User has already voted for an item in this session", string(existingVote.VoteItemID.String()))
	}

	// Create a new vote
	if err := r.conn.Create(v).Error; err != nil {
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
