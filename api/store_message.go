package api

import (
	"encoding/json"
	"fmt"
	"message-stats-api/models"
	"message-stats-api/store"
	"net/http"
)

func StoreMessage(store *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg models.Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate sender and receiver phone numbers
		if len(msg.Sender) < 1 || len(msg.Sender) > 20 || len(msg.Receiver) < 10 || len(msg.Receiver) > 15 {
			http.Error(w, "Sender must be between 1 and 20 characters long, and receiver must be between 10 and 15 characters long", http.StatusBadRequest)
			return
		}

		store.AddMessage(msg.Sender, msg.Receiver)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Message stored")
	}
}
