package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"message-stats-api/models"
	"os"
)

// Store represents a Redis-backed store for messages
type Store struct {
	client *redis.Client
	ctx    context.Context
}

// NewStore initializes a new Store with a Redis client
func NewStore() (*Store, error) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()

	// Ping the Redis server to check the connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Connected to Redis")

	return &Store{
		client: client,
		ctx:    ctx,
	}, nil
}

// AddMessage increments the message count for a given sender and receiver
func (s *Store) AddMessage(sender, receiver string) (*models.ResData, error) {
	senderRange := sender                       // Save full sender
	receiverRange := receiver[:len(receiver)-5] // Save prefix, trimming last 5 char

	// Increment the message count in Redis hash
	curCount, err := s.client.HIncrBy(s.ctx, senderRange, receiverRange, 1).Result()
	if err != nil {
		return &models.ResData{}, fmt.Errorf("failed to increment message count: %w", err)
	}

	// Prepare the response data
	resData := models.ResData{
		Sender:   sender,
		Receiver: receiver,
		Count:    curCount,
	}

	log.Printf("Message count incremented for sender: %s, receiver: %s, count: %d\n", senderRange, receiverRange, curCount)
	return &resData, nil
}

func (s *Store) GetMessage(sender, receiver string) (*models.ResData, error) {
	// Define sender range and receiver range safely
	senderRange := sender
	var receiverRange string
	if len(receiver) > 5 {
		receiverRange = receiver[:len(receiver)-5]
	} else {
		receiverRange = receiver
	}

	// Get count for a given sender and receiver
	curCount, err := s.client.Get(s.ctx, fmt.Sprintf("message_count:%s:%s", sender, receiver)).Int64()
	if errors.Is(err, redis.Nil) {
		curCount = 0 // Handle the case where the key does not exist
	} else if err != nil {
		return &models.ResData{}, fmt.Errorf("failed to get message count: %w", err)
	}

	// Prepare the response data
	resData := &models.ResData{
		Sender:   sender,
		Receiver: receiver,
		Count:    curCount,
	}

	log.Printf("Message count returned for sender: %s, receiver: %s, count: %d\n", senderRange, receiverRange, curCount)
	return resData, nil
}

// PrintMessageCountBySenderAndRange prints the message counts for all senders and receiver ranges
func (s *Store) PrintMessageCountBySenderAndRange() error {
	// Retrieve all keys from Redis
	keys, err := s.client.Keys(s.ctx, "*").Result()
	if err != nil {
		return fmt.Errorf("failed to retrieve keys from Redis: %w", err)
	}

	// Iterate over keys and fetch message counts
	for _, sender := range keys {
		data, err := s.client.HGetAll(s.ctx, sender).Result()
		if err != nil {
			return fmt.Errorf("failed to retrieve data for sender %s: %w", sender, err)
		}

		for receiver, count := range data {
			fmt.Printf("Sender range: %s, Receiver range: %s, Message count: %s\n", sender, receiver, count)
		}
	}

	return nil
}
