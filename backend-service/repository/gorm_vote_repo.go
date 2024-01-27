package repository

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/jackc/pgx/v5/pgconn"
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

func (r *gormVoteRepository) GetVoteResultsBySession(sessionID uint) ([]domain.VoteResult, error) {
	var voteResults []domain.VoteResult
	if err := r.Conn.Joins("JOIN vote_items ON votes.vote_item_id = vote_items.id").
		Where("session_id = ?", sessionID).
		Order("sum(vote_item_id) desc").
		Select("votes.*, vote_items.name as vote_item_name").
		Find(&voteResults).Error; err != nil {
		// Log the error
		log.Printf("Could not get vote results for session ID: %v. Reason: %v\n", sessionID, err)

		// Check for unique constraint violation
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			log.Printf("Could not get vote results for session ID: %v. Reason: %v\n", sessionID, pgErr.Hint)
			return nil, apperror.NewConflict("session_id", strconv.Itoa(int(sessionID)))
		}

		// If the error is not a unique constraint violation, return a generic internal server error
		log.Printf("Could not get vote results for session ID and cannot assert err to pg.Error: %v. Reason: %v\n", sessionID, err)
		return nil, apperror.NewInternal()
	}
	return voteResults, nil
}
