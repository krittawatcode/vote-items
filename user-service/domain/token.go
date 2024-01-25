package domain

import (
	"context"
	"time"
)

// TokenPair used for returning pairs of id and refresh tokens
type TokenPair struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}

// TokenService defines methods the handler layer expects to interact
// with in regards to producing JWT as string
type TokenUseCase interface {
	NewPairFromUser(ctx context.Context, u *User, prevTokenID string) (*TokenPair, error)
}

// TokenRepository defines methods it expects a repository
// it interacts with to implement
type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, prevTokenID string) error
}
