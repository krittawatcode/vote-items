package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/krittawatcode/vote-items/backend-service/database"
	"github.com/krittawatcode/vote-items/backend-service/delivery/handler"
	"github.com/krittawatcode/vote-items/backend-service/repository"
	"github.com/krittawatcode/vote-items/backend-service/usecase"
)

// will initialize a handler starting from data sources
// which inject into repository layer
// which inject into service layer
// which inject into handler layer
func inject(d *database.GormDataSources, r *database.RedisDataSources) (*gin.Engine, error) {
	log.Println("Injecting data sources")

	/*
	 * repository layer
	 */
	userRepository := repository.NewGormUserRepository(d.DB)
	tokenRepository := repository.NewTokenRepository(r.RedisClient)
	voteSessionRepository := repository.NewGormVoteSessionRepository(d.DB)
	voteItemRepository := repository.NewGormVoteItemRepository(d.DB)
	voteRepository := repository.NewGormVoteRepository(d.DB)
	/*
	 * usecase layer
	 */
	userUseCase := usecase.NewUserUseCase(userRepository)
	voteSessionUseCase := usecase.NewVoteSessionUsecase(voteSessionRepository)
	voteItemUseCase := usecase.NewVoteItemUsecase(voteItemRepository)
	voteUseCase := usecase.NewVoteUsecase(voteRepository)

	// load rsa keys
	privKeyFile := os.Getenv("PRIV_KEY_FILE")
	priv, err := os.ReadFile(privKeyFile)

	if err != nil {
		return nil, fmt.Errorf("could not read private key pem file: %w", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)

	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	pubKeyFile := os.Getenv("PUB_KEY_FILE")
	pub, err := os.ReadFile(pubKeyFile)

	if err != nil {
		return nil, fmt.Errorf("could not read public key pem file: %w", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)

	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}

	// load refresh token secret from env variable
	refreshSecret := os.Getenv("REFRESH_SECRET")

	// load expiration lengths from env variables and parse as int
	idTokenExp := os.Getenv("ID_TOKEN_EXP")
	refreshTokenExp := os.Getenv("REFRESH_TOKEN_EXP")

	idExp, err := strconv.ParseInt(idTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse ID_TOKEN_EXP as int: %w", err)
	}

	refreshExp, err := strconv.ParseInt(refreshTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse REFRESH_TOKEN_EXP as int: %w", err)
	}

	tokenUseCase := usecase.NewTokenUseCase(tokenRepository, privKey, pubKey, refreshSecret, idExp, refreshExp)

	// initialize gin.Engine
	router := gin.Default()

	// read in API_URL
	baseURL := os.Getenv("API_URL")
	userPath := os.Getenv("USER_PATH")

	// read in HANDLER_TIMEOUT
	handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
	ht, err := strconv.ParseInt(handlerTimeout, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
	}

	handler.NewUserHandler(router, userUseCase, tokenUseCase, baseURL+userPath, time.Duration(time.Duration(ht)*time.Second))

	/*
	 * setup vote item handler
	 */
	handler.NewVoteItemsHandler(router, voteItemUseCase, voteUseCase, voteSessionUseCase, tokenUseCase, baseURL+"/vote_items", time.Duration(time.Duration(ht)*time.Second))

	return router, nil
}
