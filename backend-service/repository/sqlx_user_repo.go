package repository

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/apperror"
	"github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

// pgUserRepository is data/repository implementation
// of service layer UserRepository
type pgUserRepository struct {
	DB *sqlx.DB
}

// NewUserRepository is a factory for initializing User Repositories
func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &pgUserRepository{
		DB: db,
	}
}

// Create reaches out to database SQLX api
func (r *pgUserRepository) Create(ctx context.Context, u *domain.User) error {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *"

	if err := r.DB.GetContext(ctx, u, query, u.Email, u.Password); err != nil {
		// check unique constraint
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, err.Code.Name())
			return apperror.NewConflict("email", u.Email)
		}

		log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, err)
		return apperror.NewInternal()
	}
	return nil
}

// FindByID fetches user by id
func (r *pgUserRepository) FindByID(ctx context.Context, uid uuid.UUID) (*domain.User, error) {
	user := &domain.User{}

	query := "SELECT * FROM users WHERE uid=$1"

	// we need to actually check errors as it could be something other than not found
	if err := r.DB.GetContext(ctx, user, query, uid); err != nil {
		return user, apperror.NewNotFound("uid", uid.String())
	}

	return user, nil
}

// FindByEmail retrieves user row by email address
func (r *pgUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}

	query := "SELECT * FROM users WHERE email=$1"

	if err := r.DB.GetContext(ctx, user, query, email); err != nil {
		log.Printf("Unable to get user with email address: %v. Err: %v\n", email, err)
		return user, apperror.NewNotFound("email", email)
	}

	return user, nil
}
