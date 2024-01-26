package repository

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/apperror"
)

type gormUserRepository struct {
	conn *gorm.DB
}

// NewGormUserRepository ...
func NewGormUserRepository(conn *gorm.DB) domain.UserRepository {
	return &gormUserRepository{conn}
}

// Create ...
func (r *gormUserRepository) Create(ctx context.Context, u *domain.User) error {
	// implementation of the Create method
	if err := r.conn.Create(u).Error; err != nil {
		log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, err)
		// check unique constraint
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, pgErr.Hint)
			return apperror.NewConflict("email", u.Email)
		}

		log.Printf("Could not create a user with email and can not assert err to pg.Error: %v. Reason: %v\n", u.Email, err)
		return apperror.NewInternal()
	}
	return nil
}

// FindByID fetches user by id
func (r *gormUserRepository) FindByID(ctx context.Context, uid uuid.UUID) (*domain.User, error) {
	user := &domain.User{}
	if err := r.conn.Where("uid = ?", uid).First(user).Error; err != nil {
		return user, apperror.NewNotFound("uid", uid.String())
	}
	return user, nil
}

// FindByEmail retrieves user row by email address
func (r *gormUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	if err := r.conn.Where("email = ?", email).First(user).Error; err != nil {
		return user, apperror.NewNotFound("email", email)
	}
	return user, nil
}
