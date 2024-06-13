package api

import (
	"encoding/json"
	"message-stats-api/models"
	"message-stats-api/store"
	"net/http"
)

func StoreMessage(store *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure the request is a POST request
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var msg models.Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Validate sender and receiver phone numbers
		if len(msg.Sender) < 1 || len(msg.Sender) > 20 || len(msg.Receiver) < 10 || len(msg.Receiver) > 15 {
			http.Error(w, "Sender must be between 1 and 20 characters long, and receiver must be between 10 and 15 characters long", http.StatusBadRequest)
			return
		}

		respData, err := store.AddMessage(msg.Sender, msg.Receiver)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := models.Response{
			Data:    respData,
			Message: "Data stored",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
