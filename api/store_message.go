package api

import (
	"encoding/json"
	"fmt"
	"message-stats-api/models"
	"net/http"
	"sync"
)

var (
	Messages      []models.Message
	MessagesMutex sync.Mutex
)

func StoreMessage(w http.ResponseWriter, r *http.Request) {
	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate sender and receiver phone numbers
	if len(msg.Sender) != 10 || len(msg.Receiver) != 10 {
		http.Error(w, "Sender and receiver phone numbers must be exactly 10 digits long", http.StatusBadRequest)
		return
	}

	MessagesMutex.Lock()
	Messages = append(Messages, msg)
	MessagesMutex.Unlock()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Message stored")
}
