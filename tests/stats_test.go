package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"message-stats-api/api"
	"message-stats-api/models"
	"message-stats-api/store"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSendMessages(t *testing.T) {
	var wg sync.WaitGroup
	const numMessages = 10000

	var sentCount int32

	rand.Seed(time.Now().UnixNano())

	// Create a new store instance
	newStore, err := store.NewStore()
	if err != nil {
		log.Fatalf("Error initializing message store: %v", err)
	}

	// Create a test server
	ts := httptest.NewServer(api.StoreMessage(newStore))
	defer ts.Close()

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
			case i < 4000:
				sender = fmt.Sprintf("fb34567%05d", rand.Intn(100))
				receiver = fmt.Sprintf("fb76543%05d", rand.Intn(100))
			case i < 7000:
				sender = fmt.Sprintf("v45678%05d", rand.Intn(400))
				receiver = fmt.Sprintf("v87654%09d", rand.Intn(400))
			case i < 9990:
				sender = fmt.Sprintf("56789%05d", rand.Intn(100))
				receiver = fmt.Sprintf("98765%05d", rand.Intn(100))
			default:
				sender = generateNumber()
				receiver = generateNumber()
			}

			text := "Hello"

			if len(sender) >= 1 && len(sender) <= 20 && len(receiver) >= 10 && len(receiver) <= 15 {
				if err := sendMessage(ts.URL, sender, receiver, text); err != nil {
					t.Errorf("failed to send message: %v", err)
				} else {
					atomic.AddInt32(&sentCount, 1)
				}
			} else {
				t.Errorf("sender or receiver number is not of valid length")
			}
		}(i)
	}
	wg.Wait()

	t.Logf("Sent %d messages", sentCount)

	err = newStore.PrintMessageCountBySenderAndRange()
	if err != nil {
		t.Errorf("failed to print message count: %v", err)
	}
}

// generateNumber generates a random 10-digit number as a string
func generateNumber() string {
	return fmt.Sprintf("%010d", rand.Intn(10000000000))
}

func sendMessage(url, sender, receiver, text string) error {
	msg := models.Message{Sender: sender, Receiver: receiver, Text: text}
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(url+"/store-message", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status OK but got %v", resp.StatusCode)
	}
	return nil
}
