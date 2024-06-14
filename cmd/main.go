package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"message-stats-api/api"
	"message-stats-api/store"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	loadEnv()

	serveApplication()
}

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// serveApplication starts the HTTP server and sets up graceful shutdown
func serveApplication() {
	// Create a new store instance
	messageStore, err := store.NewStore()
	if err != nil {
		log.Fatalf("Error initializing message store: %v", err)
	}

	// Pass the store instance to the handlers
	http.HandleFunc("/store", api.Store(messageStore))
	http.HandleFunc("/stats", api.Stats(messageStore))

	// Define the server
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Start the server in a new goroutine
	go func() {
		log.Println("Starting server on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on port 8080: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop

	// Create a context with timeout for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
