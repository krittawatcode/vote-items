package usecase

import (
	"context"
	"crypto/rsa"
	"log"

	"github.com/krittawatcode/vote-items/user-service/domain"
	"github.com/krittawatcode/vote-items/user-service/domain/apperror"
)

// tokenUseCase acts as a struct for injecting an implementation of TokenRepository
// for use in service methods
type tokenUseCase struct {
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

// NewTokenUseCase is a factory function for
// initializing a TokenUseCase with its usecase layer dependencies
func NewTokenUseCase(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey, refreshSecret string, idExpirationSecs int64, refreshExpirationSecs int64) domain.TokenUseCase {
	return &tokenUseCase{
		PrivKey:               privKey,
		PubKey:                pubKey,
		RefreshSecret:         refreshSecret,
		IDExpirationSecs:      idExpirationSecs,
		RefreshExpirationSecs: refreshExpirationSecs,
	}
}

// NewPairFromUser creates fresh id and refresh tokens for the current user
// If a previous token is included, the previous token is removed from
// the tokens repository
func (s *tokenUseCase) NewPairFromUser(ctx context.Context, u *domain.User, prevTokenID string) (*domain.TokenPair, error) {
	// No need to use a repository for idToken as it is unrelated to any data source
	idToken, err := generateIDToken(u, s.PrivKey, s.IDExpirationSecs)

	if err != nil {
		log.Printf("Error generating idToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperror.NewInternal()
	}

	refreshToken, err := generateRefreshToken(u.UID, s.RefreshSecret, s.RefreshExpirationSecs)

	if err != nil {
		log.Printf("Error generating refreshToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperror.NewInternal()
	}

	// TODO: store refresh tokens by calling TokenRepository methods

	return &domain.TokenPair{
		IDToken:      idToken,
		RefreshToken: refreshToken.SS,
	}, nil
}
