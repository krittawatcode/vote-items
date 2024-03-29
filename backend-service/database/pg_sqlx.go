package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DataSources struct {
	DB *sqlx.DB
}

// InitDS establishes connections to fields in dataSources
func (ds *DataSources) InitDS() error {
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
	db, err := sqlx.Open("postgres", pgConnString)

	if err != nil {
		return fmt.Errorf("error opening db: %w", err)
	}

	// Verify database connection is working
	if err := db.Ping(); err != nil {
		return fmt.Errorf("error connecting to db: %w", err)
	}

	ds.DB = db
	return nil
}

// close to be used in graceful server shutdown
func (ds *DataSources) Close() error {
	if err := ds.DB.Close(); err != nil {
		return fmt.Errorf("error closing Postgresql: %w", err)
	}

	return nil
}
