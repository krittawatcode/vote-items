package repository

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/apperror"
	"gorm.io/gorm"
)

type gormVoteSessionRepository struct {
	conn *gorm.DB
}

// NewGormVoteSessionRepository ...
func NewGormVoteSessionRepository(conn *gorm.DB) domain.VoteSessionRepository {
	return &gormVoteSessionRepository{conn}
}

func (r *gormVoteSessionRepository) GetOpenVoteSession() (*domain.VoteSession, error) {
	voteSession := &domain.VoteSession{}
	db := r.conn.Where("is_open = ?", true).First(voteSession)

	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, db.Error
	}

	return voteSession, nil
}

func (r *gormVoteSessionRepository) CreateVoteSession(id uint) error {
	voteSession := domain.VoteSession{
		ID:     id,
		IsOpen: true,
	}
	if err := r.conn.Create(&voteSession).Error; err != nil {
		log.Printf("Could not create a vote session with id: %v. Reason: %v\n", id, err)
		// check unique constraint
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			log.Printf("Could not create a  vote session with id: %v. Reason: %v\n", id, pgErr.Hint)
			return apperror.NewConflict("id", strconv.Itoa(int(id)))
		}
	}
	return nil
}

// CloseVoteSession closes a vote session by its ID.
// It updates the IsOpen field of the vote session to false in the database.
// If the operation fails, it logs the error and returns an appropriate error.
func (r *gormVoteSessionRepository) CloseVoteSession(id uint) error {
	// Create a VoteSession object with the provided ID and IsOpen set to false
	voteSession := domain.VoteSession{
		ID:     id,
		IsOpen: false,
	}

	// Attempt to save the updated vote session in the database
	if err := r.conn.Save(&voteSession).Error; err != nil {
		// If an error occurred, check if it's a postgres error
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			// If it's a postgres error, log the error and return a Conflict error
			log.Printf("Could not close the vote session with id: %v. Reason: %v\n", id, pgErr.Hint)
			return apperror.NewConflict("id", strconv.Itoa(int(id)))
		}
		// For any other error, return an Internal error
		return apperror.NewInternal()
	}

	// If no error occurred, return nil
	return nil
}

// GetVoteSessionByID retrieves a vote session by its ID from the database.
// It returns a pointer to the VoteSession object if found, or an error if not found or any other issue occurred.
func (r *gormVoteSessionRepository) GetVoteSessionByID(ctx context.Context, id uint) (*domain.VoteSession, error) {
	// Initialize a new VoteSession object
	voteSession := &domain.VoteSession{}

	// Query the database for a vote session with the provided ID
	db := r.conn.Where("id = ?", id).First(voteSession)

	// If an error occurred during the query
	if db.Error != nil {
		// If the error is due to the record not being found
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			// Return a new NotFound error
			return nil, apperror.NewNotFound("vote session", strconv.Itoa(int(id)))
		}
		// For any other error, return a new Internal error
		return nil, apperror.NewInternal()
	}

	// If no error occurred, return the retrieved vote session
	return voteSession, nil
}
