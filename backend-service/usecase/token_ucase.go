package usecase

import (
	"context"
	"crypto/rsa"
	"log"

	"github.com/google/uuid"
	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/apperror"
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
	if prevTokenID != "" {
		if err := s.TokenRepository.DeleteRefreshToken(ctx, u.UID.String(), prevTokenID); err != nil {
			log.Printf("Could not delete previous refreshToken for uid: %v, tokenID: %v\n", u.UID.String(), prevTokenID)

			return nil, err
		}
	}

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
	if err := s.TokenRepository.SetRefreshToken(ctx, u.UID.String(), refreshToken.ID.String(), refreshToken.ExpiresIn); err != nil {
		log.Printf("Error storing tokenID for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperror.NewInternal()
	}

	return &domain.TokenPair{
		IDToken:      domain.IDToken{SS: idToken},
		RefreshToken: domain.RefreshToken{SS: refreshToken.SS, ID: refreshToken.ID, UID: u.UID},
	}, nil
}

// ValidateIDToken validates the id token jwt string
// It returns the user extract from the IDTokenCustomClaims
func (s *tokenUseCase) ValidateIDToken(tokenString string) (*domain.User, error) {
	claims, err := validateIDToken(tokenString, s.PubKey) // uses public RSA key

	// We'll just return unauthorized error in all instances of failing to verify user
	if err != nil {
		log.Printf("Unable to validate or parse idToken - Error: %v\n", err)
		return nil, apperror.NewAuthorization("Unable to verify user from idToken")
	}

	return claims.User, nil
}

// ValidateRefreshToken checks to make sure the JWT provided by a string is valid
// and returns a RefreshToken if valid
func (s *tokenUseCase) ValidateRefreshToken(tokenString string) (*domain.RefreshToken, error) {
	// validate actual JWT with string a secret
	claims, err := validateRefreshToken(tokenString, s.RefreshSecret)

	// We'll just return unauthorized error in all instances of failing to verify user
	if err != nil {
		log.Printf("Unable to validate or parse refreshToken for token string: %s\n%v\n", tokenString, err)
		return nil, apperror.NewAuthorization("Unable to verify user from refresh token")
	}

	// Standard claims store ID as a string. I want "model" to be clear our string
	// is a UUID. So we parse claims.Id as UUID
	tokenUUID, err := uuid.Parse(claims.Id)

	if err != nil {
		log.Printf("Claims ID could not be parsed as UUID: %s\n%v\n", claims.Id, err)
		return nil, apperror.NewAuthorization("Unable to verify user from refresh token")
	}

	return &domain.RefreshToken{
		SS:  tokenString,
		ID:  tokenUUID,
		UID: claims.UID,
	}, nil
}
