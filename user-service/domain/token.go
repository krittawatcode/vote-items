package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// RefreshToken stores token properties that
// are accessed in multiple application layers
type RefreshToken struct {
	ID  uuid.UUID `json:"-"`
	UID uuid.UUID `json:"-"`
	SS  string    `json:"refreshToken"`
}

// IDToken stores token properties that
// are accessed in multiple application layers
type IDToken struct {
	SS string `json:"idToken"`
}

// TokenPair used for returning pairs of id and refresh tokens
type TokenPair struct {
	IDToken
	RefreshToken
}

// TokenService defines methods the handler layer expects to interact
// with in regards to producing JWT as string
type TokenUseCase interface {
	NewPairFromUser(ctx context.Context, u *User, prevTokenID string) (*TokenPair, error)
	ValidateIDToken(tokenString string) (*User, error)
	ValidateRefreshToken(refreshTokenString string) (*RefreshToken, error)
}

// TokenRepository defines methods it expects a repository
// it interacts with to implement
type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID string, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, prevTokenID string) error
}
