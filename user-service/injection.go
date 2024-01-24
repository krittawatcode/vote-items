package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/krittawatcode/vote-items/user-service/database"
	"github.com/krittawatcode/vote-items/user-service/delivery/handler"
	"github.com/krittawatcode/vote-items/user-service/repository"
	"github.com/krittawatcode/vote-items/user-service/usecase"
)

// will initialize a handler starting from data sources
// which inject into repository layer
// which inject into service layer
// which inject into handler layer
func inject(d *database.GormDataSources) (*gin.Engine, error) {
	log.Println("Injecting data sources")

	/*
	 * repository layer
	 */
	userRepository := repository.NewGormUserRepository(d.DB)

	/*
	 * usecase layer
	 */
	userUseCase := usecase.NewUserUseCase(userRepository)

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

	tokenUseCase := usecase.NewTokenUseCase(privKey, pubKey, refreshSecret)

	// initialize gin.Engine
	router := gin.Default()

	handler.NewUserHandler(router, userUseCase, tokenUseCase)

	return router, nil
}
