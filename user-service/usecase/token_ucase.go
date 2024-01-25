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
	TokenRepository       domain.TokenRepository
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

// NewTokenUseCase is a factory function for
// initializing a TokenUseCase with its usecase layer dependencies
func NewTokenUseCase(r domain.TokenRepository, privKey *rsa.PrivateKey, pubKey *rsa.PublicKey, refreshSecret string, idExpirationSecs int64, refreshExpirationSecs int64) domain.TokenUseCase {
	return &tokenUseCase{
		TokenRepository:       r,
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

	// set freshly minted refresh token to valid list
	if err := s.TokenRepository.SetRefreshToken(ctx, u.UID.String(), refreshToken.ID, refreshToken.ExpiresIn); err != nil {
		log.Printf("Error storing tokenID for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperror.NewInternal()
	}

	// delete user's current refresh token (used when refreshing idToken)
	if prevTokenID != "" {
		if err := s.TokenRepository.DeleteRefreshToken(ctx, u.UID.String(), prevTokenID); err != nil {
			log.Printf("Could not delete previous refreshToken for uid: %v, tokenID: %v\n", u.UID.String(), prevTokenID)
		}
	}

	return &domain.TokenPair{
		IDToken:      idToken,
		RefreshToken: refreshToken.SS,
	}, nil
}
