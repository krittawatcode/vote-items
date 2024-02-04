package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/krittawatcode/vote-items/backend-service/database"
	"github.com/krittawatcode/vote-items/backend-service/delivery/handler"
	"github.com/krittawatcode/vote-items/backend-service/docs"
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
	voteResultRepository := repository.NewGormVoteResultRepository(d.DB)
	/*
	 * usecase layer
	 */
	userUseCase := usecase.NewUserUseCase(userRepository)
	voteSessionUseCase := usecase.NewVoteSessionUsecase(voteSessionRepository)
	voteItemUseCase := usecase.NewVoteItemUsecase(voteItemRepository)
	voteUseCase := usecase.NewVoteUsecase(voteRepository)
	voteResultUseCase := usecase.NewVoteResultUsecase(voteResultRepository)
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
	// Add a health check endpoint

	// set up swagger
	// docs.SwaggerInfo.BasePath = "/api/v1"
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// read in API_URL
	baseURL := os.Getenv("API_URL")
	userPath := os.Getenv("USER_PATH")
	voteSessionPath := os.Getenv("VOTE_SESSION_PATH")
	voteItemPath := os.Getenv("VOTE_ITEM_PATH")
	votePath := os.Getenv("VOTE_PATH")
	voteResultPath := os.Getenv("VOTE_RESULT_PATH")

	// read in HANDLER_TIMEOUT
	handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
	ht, err := strconv.ParseInt(handlerTimeout, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
	}
	timeout := time.Duration(time.Duration(ht) * time.Second)

	/*
	 * setup handler
	 */
	handler.NewUserHandler(router, userUseCase, tokenUseCase, baseURL+userPath, timeout)
	handler.NewVoteSessionsHandler(router, voteSessionUseCase, tokenUseCase, baseURL+voteSessionPath, timeout)
	handler.NewVoteItemsHandler(router, voteItemUseCase, tokenUseCase, baseURL+voteItemPath, timeout)
	handler.NewVotesHandler(router, voteUseCase, tokenUseCase, baseURL+votePath, timeout)
	handler.NewVoteResultsHandler(router, voteResultUseCase, tokenUseCase, baseURL+voteResultPath, timeout)

	// set up swagger
	docs.SwaggerInfo.BasePath = baseURL
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "running"})
	})

	return router, nil
}
