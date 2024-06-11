package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"message-stats-api/api"
	"message-stats-api/models"
	"message-stats-api/store"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Start the server
	http.HandleFunc("/store-message", api.StoreMessage)

	server := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	go func() {
		log.Println("Starting server on port 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on port 8080: %v\n", err)
		}
	}()

	// Wait for server to start
	time.Sleep(time.Second)

	// Run tests
	code := m.Run()

	// Shutdown server after tests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	os.Exit(code)
}

func TestSendMessages(t *testing.T) {
	var wg sync.WaitGroup
	const numMessages = 10000

	var sentCount int32

	rand.Seed(time.Now().UnixNano())

	newStore := store.NewStore()

	for i := 0; i < numMessages; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			var sender, receiver string

			switch {
			case i < 300:
				sender = fmt.Sprintf("12345%05d", rand.Intn(30))
				receiver = fmt.Sprintf("54321%05d", rand.Intn(30))
			case i < 1300:
				sender = fmt.Sprintf("23456%05d", rand.Intn(100))
				receiver = fmt.Sprintf("65432%05d", rand.Intn(100))
			case i < 3000:
				sender = fmt.Sprintf("34567%05d", rand.Intn(100))
				receiver = fmt.Sprintf("76543%05d", rand.Intn(100))
			case i < 7000:
				sender = fmt.Sprintf("45678%05d", rand.Intn(100))
				receiver = fmt.Sprintf("87654%05d", rand.Intn(100))
			case i < 9990:
				sender = fmt.Sprintf("56789%05d", rand.Intn(100))
				receiver = fmt.Sprintf("98765%05d", rand.Intn(100))
			default:
				sender = generateNumber()
				receiver = generateNumber()
			}

			text := "Hello"

			if len(sender) == 10 && len(receiver) == 10 {
				if err := sendMessage(sender, receiver, text); err != nil {
					t.Errorf("failed to send message: %v", err)
				} else {
					atomic.AddInt32(&sentCount, 1)
				}
			} else {
				t.Errorf("sender or receiver number is not 10 digits long")
			}

			// Save message in newStore
			newStore.Lock()
			if newStore.Data == nil {
				newStore.Data = make(map[string]map[string]int)
			}
			senderRange := sender[:len(sender)-5]
			receiverRange := receiver[:len(receiver)-5]
			if newStore.Data[senderRange] == nil {
				newStore.Data[senderRange] = make(map[string]int)
			}
			newStore.Data[senderRange][receiverRange]++
			newStore.Unlock()
		}(i)
	}
	wg.Wait()

	t.Logf("Sent %d messages", sentCount)

	printMessageCountBySenderAndRange(newStore)
}

// generateNumber generates a random 10-digit number as a string
func generateNumber() string {
	return fmt.Sprintf("%010d", rand.Intn(10000000000))
}

func sendMessage(sender, receiver, text string) error {
	msg := models.Message{Sender: sender, Receiver: receiver, Text: text}
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:8080/store-message", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status OK but got %v", resp.StatusCode)
	}
	return nil
}

func printMessageCountBySenderAndRange(store *store.Store) {
	store.RLock()
	defer store.RUnlock()

	for sender, receiverMap := range store.Data {
		for receiver, count := range receiverMap {
			fmt.Printf("Sender range: %s, Receiver range: %s, Message count: %d\n", sender, receiver, count)
		}
	}
}
