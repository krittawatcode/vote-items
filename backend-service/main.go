package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/krittawatcode/vote-items/backend-service/database"
	"github.com/krittawatcode/vote-items/backend-service/domain"
)

// @title Vote Items API
func main() {
	log.Println("Starting server...")

	// initialize data sources
	ds := new(database.GormDataSources)
	err := ds.InitDS()
	if err != nil {
		log.Fatalf("Unable to initialize data sources: %v\n", err)
	}
	ds.DB.AutoMigrate(&domain.User{}, &domain.VoteSession{}, &domain.VoteItem{}, &domain.Vote{})

	err = ds.SeedUsers()
	if err != nil {
		log.Fatalf("Unable to seed users: %v\n", err)
	}

	rc := new(database.RedisDataSources)
	err = rc.InitRC()
	if err != nil {
		log.Fatalf("Unable to initialize data sources: %v\n", err)
	}

	router, err := inject(ds, rc)
	if err != nil {
		log.Fatalf("Unable to inject data sources: %v\n", err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Graceful server shutdown - https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize server: %v\n", err)
		}
	}()

	log.Printf("Listening on port %v\n", srv.Addr)

	// Wait for kill signal of channel
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This blocks until a signal is passed into the quit channel
	<-quit

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown data sources
	if err := ds.Close(); err != nil {
		log.Fatalf("A problem occurred gracefully shutting down data sources: %v\n", err)
	}

	// Shutdown server
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
}
