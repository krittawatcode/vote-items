package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/krittawatcode/vote-items/backend-service/domain"
	"github.com/krittawatcode/vote-items/backend-service/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDataSources struct {
	DB *gorm.DB
}

func (ds *GormDataSources) InitDS() error {
	log.Printf("Initializing data sources\n")
	// load env variables - we could pass these in,
	// but this is sort of just a top-level (main package)
	// helper function, so I'll just read them in here
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")
	pgDB := os.Getenv("PG_DB")
	pgSSL := os.Getenv("PG_SSL")

	pgConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", pgHost, pgPort, pgUser, pgPassword, pgDB, pgSSL)

	log.Printf("Connecting to Postgresql\n")
	db, err := gorm.Open(postgres.Open(pgConnString), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error opening db: %w", err)
	}

	ds.DB = db
	return nil
}

// close to be used in graceful server shutdown
func (ds *GormDataSources) Close() error {
	sqlDB, err := ds.DB.DB()
	if err != nil {
		return fmt.Errorf("error get db connection: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("error closing Postgresql: %w", err)
	}

	return nil
}

var mu sync.Mutex

func (ds *GormDataSources) SeedUsers() error {
	users := []domain.User{
		{Email: "admin@mtl.co.th", Password: "adminPassword"},
		{Email: "krittawat@mercy.gg", Password: "userPassword"},
	}

	for _, user := range users {
		hashedPassword, err := usecase.HashPassword(user.Password)
		if err != nil {
			return err
		}

		user.Password = hashedPassword

		mu.Lock()
		// Check if the user already exists
		var existingUser domain.User
		if err := ds.DB.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				// If the error is not a RecordNotFound error, return the error
				mu.Unlock()
				return err
			}
		} else {
			// If no error occurred, the user already exists, so skip to the next user
			mu.Unlock()
			continue
		}

		// If the user does not exist, create the user
		if err := ds.DB.Create(&user).Error; err != nil {
			mu.Unlock()
			return err
		}
		mu.Unlock()
	}

	return nil
}
