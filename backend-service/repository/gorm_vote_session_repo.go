package repository

import (
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
		ID:     int32(id),
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

func (r *gormVoteSessionRepository) CloseVoteSession(id uint) error {
	voteSession := domain.VoteSession{
		ID:     int32(id),
		IsOpen: false,
	}
	if err := r.conn.Save(&voteSession).Error; err != nil {
		return apperror.NewInternal()
	}
	return nil
}
