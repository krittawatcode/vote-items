package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/krittawatcode/vote-items/user-service/domain"
	"github.com/krittawatcode/vote-items/user-service/domain/apperror"
)

// redisTokenRepository is data/repository implementation
// of service layer TokenRepository
type redisTokenRepository struct {
	Redis *redis.Client
}

// NewTokenRepository is a factory for initializing User Repositories
func NewTokenRepository(redisClient *redis.Client) domain.TokenRepository {
	return &redisTokenRepository{
		Redis: redisClient,
	}
}

// SetRefreshToken stores a refresh token with an expiry time
func (r *redisTokenRepository) SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error {
	// We'll store userID with token id so we can scan (non-blocking)
	// over the user's tokens and delete them in case of token leakage
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := r.Redis.Set(key, 0, expiresIn).Err(); err != nil {
		log.Printf("Could not SET refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return apperror.NewInternal()
	}
	return nil
}

// DeleteRefreshToken used to delete old  refresh tokens
// Services my access this to revolve tokens
func (r *redisTokenRepository) DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)

	result := r.Redis.Del(key)

	if err := result.Err(); err != nil {
		log.Printf("Could not delete refresh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return apperror.NewInternal()
	}

	// Val returns count of deleted keys.
	// If no key was deleted, the refresh token is invalid
	if result.Val() < 1 {
		log.Printf("Refresh token to redis for userID/tokenID: %s/%s does not exist\n", userID, tokenID)
		return apperror.NewAuthorization("Invalid refresh token")
	}

	return nil
}
