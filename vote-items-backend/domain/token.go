package domain

import "context"

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
