package usecase

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/domain/apperror"
	"github.com/krittawatcode/vote-items/backend-service/domain/appmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func TestNewPairFromUser(t *testing.T) {
	priv, err := os.ReadFile("../cert/rsa_private_test.pem")
	if err != nil {
		t.Fatalf("Failed to read private key: %v", err)
	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		t.Fatalf("Failed to parse private key: %v", err)
	}
	pub, _ := os.ReadFile("../cert/rsa_public_test.pem")
	if err != nil {
		t.Fatalf("Failed to read public key: %v", err)
	}
	pubKey, _ := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		t.Fatalf("Failed to parse public key: %v", err)
	}

	secret := "anotsorandomtestsecret"

	// define token expiration times
	idTokenExp := int64(15 * 60)
	refreshTokenExp := int64(3 * 24 * 2600)

	// include password to make sure it is not serialized
	// since json tag is "-"
	uid, err := uuid.NewRandom()
	if err != nil {
		t.Fatalf("Failed to generate UUID: %v", err)
	}

	u := &domain.User{
		UID:      uid,
		Email:    "bob@bob.com",
		Password: "blarghedymcblarghface",
	}
	prevID := "a_previous_tokenID"

	setSuccessArguments := mock.Arguments{
		mock.Anything,
		u.UID.String(),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
	}

	deleteWithPrevIDArguments := mock.Arguments{
		mock.Anything,
		u.UID.String(),
		prevID,
	}

	t.Run("Returns a token pair with proper values", func(t *testing.T) {
		ctx := context.Background()

		mockTokenRepository := new(appmock.MockTokenRepository)
		tokenUseCase := NewTokenUseCase(mockTokenRepository, privKey, pubKey, secret, idTokenExp, refreshTokenExp)

		mockTokenRepository.On("SetRefreshToken", setSuccessArguments...).Return(nil)
		mockTokenRepository.On("DeleteRefreshToken", deleteWithPrevIDArguments...).Return(nil)
		tokenPair, err := tokenUseCase.NewPairFromUser(ctx, u, prevID)
		assert.NoError(t, err)

		// SetRefreshToken should be called with setSuccessArguments
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		// DeleteRefreshToken should not be called since prevID is ""
		mockTokenRepository.AssertCalled(t, "DeleteRefreshToken", deleteWithPrevIDArguments...)

		var s string
		assert.IsType(t, s, tokenPair.IDToken.SS)

		// decode the Base64URL encoded string
		// simpler to use jwt library which is already imported
		idTokenClaims := &idTokenCustomClaims{}

		_, err = jwt.ParseWithClaims(tokenPair.IDToken.SS, idTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return pubKey, nil
		})

		assert.NoError(t, err)

		// assert claims on idToken
		expectedClaims := []interface{}{
			u.UID,
			u.Email,
		}
		actualIDClaims := []interface{}{
			idTokenClaims.User.UID,
			idTokenClaims.User.Email,
		}

		assert.ElementsMatch(t, expectedClaims, actualIDClaims)
		assert.Empty(t, idTokenClaims.User.Password) // password should never be encoded to json

		expiresAt := time.Unix(idTokenClaims.StandardClaims.ExpiresAt, 0)
		expectedExpiresAt := time.Now().Add(time.Duration(idTokenExp) * time.Second)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)

		refreshTokenClaims := &refreshTokenCustomClaims{}
		_, err = jwt.ParseWithClaims(tokenPair.RefreshToken.SS, refreshTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		assert.IsType(t, s, tokenPair.RefreshToken.SS)

		// assert claims on refresh token
		assert.NoError(t, err)
		assert.Equal(t, u.UID, refreshTokenClaims.UID)

		expiresAt = time.Unix(refreshTokenClaims.StandardClaims.ExpiresAt, 0)
		expectedExpiresAt = time.Now().Add(time.Duration(refreshTokenExp) * time.Second)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)
	})

	t.Run("Error setting refresh token", func(t *testing.T) {
		ctx := context.Background()

		mockTokenRepository := new(appmock.MockTokenRepository)
		tokenUseCase := NewTokenUseCase(mockTokenRepository, privKey, pubKey, secret, idTokenExp, refreshTokenExp)

		// mock call argument/responses
		mockErr := apperror.NewInternal()
		mockTokenRepository.On("SetRefreshToken", setSuccessArguments...).Return(mockErr)

		_, err := tokenUseCase.NewPairFromUser(ctx, u, prevID)

		assert.Error(t, err)
		// SetRefreshToken should be called with setSuccessArguments
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		// DeleteRefreshToken should not be called since prevID is ""
		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")

	})

	// case where prevID is not provided
	// should call SetRefreshToken but not DeleteRefreshToken
	t.Run("Empty string provided for prevID", func(t *testing.T) {
		ctx := context.Background()

		mockTokenRepository := new(appmock.MockTokenRepository)
		tokenUseCase := NewTokenUseCase(mockTokenRepository, privKey, pubKey, secret, idTokenExp, refreshTokenExp)

		mockTokenRepository.On("SetRefreshToken", setSuccessArguments...).Return(nil)

		_, err := tokenUseCase.NewPairFromUser(ctx, u, "")
		assert.NoError(t, err)

		// SetRefreshToken should be called with setSuccessArguments
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		// DeleteRefreshToken should not be called since prevID is ""
		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})

	t.Run("Error deleting refresh token", func(t *testing.T) {
		mockTokenRepository := new(appmock.MockTokenRepository)
		tokenUseCase := NewTokenUseCase(mockTokenRepository, privKey, pubKey, secret, idTokenExp, refreshTokenExp)

		ctx := context.Background()
		// mock call argument/responses
		mockErr := apperror.NewInternal()
		mockTokenRepository.On("SetRefreshToken", setSuccessArguments...).Return(nil)
		mockTokenRepository.On("DeleteRefreshToken", deleteWithPrevIDArguments...).Return(mockErr)

		_, err := tokenUseCase.NewPairFromUser(ctx, u, prevID)

		assert.NoError(t, err)
		// SetRefreshToken should be called with setSuccessArguments
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		// DeleteRefreshToken should be called since prevID is not ""
		mockTokenRepository.AssertCalled(t, "DeleteRefreshToken", deleteWithPrevIDArguments...)
	})
}
