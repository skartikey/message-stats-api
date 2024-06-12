package main

import (
	"context"
	"github.com/go-redis/redis/v8"
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

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	rdb := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
	})

	// Check the connection
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")

	serveApplication()
}

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func serveApplication() {
	// Create a new store instance
	messageStore := store.NewStore()

	// Pass the store instance to the StoreMessage handler
	http.HandleFunc("/store-message", api.StoreMessage(messageStore))

	server := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("Starting server on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on port 8080: %v\n", err)
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
