package api

import (
	"encoding/json"
	"message-stats-api/models"
	"message-stats-api/store"
	"net/http"
	"regexp"
	"strings"
)

func Store(store *store.Store) http.HandlerFunc {
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

		// Sanitise input
		// Trim '+' prefix if any
		if strings.HasPrefix(msg.Receiver, "+") {
			msg.Receiver = msg.Receiver[1:]
		}
		// Remove leading zeros
		msg.Receiver = strings.TrimLeft(msg.Receiver, "0")

		if !validateSender(msg.Sender) {
			http.Error(w, "Sender must be between 1 and 20 characters long", http.StatusBadRequest)
			return
		}

		if !validateReceiver(msg.Receiver) {
			http.Error(w, "Receiver must be numeric between 10 and 15 digit long", http.StatusBadRequest)
			return
		}

		respData, err := store.AddMessage(msg.Sender, msg.Receiver)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := models.Response{
			Data:    *respData,
			Message: "Data stored",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Stats(store *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure the request is a GET request
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var statsReq models.StatsRequest
		if err := json.NewDecoder(r.Body).Decode(&statsReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		respData, err := store.GetMessage(statsReq.Sender, statsReq.Receiver)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := models.Response{
			Data:    *respData,
			Message: "Success",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func validateSender(sender string) bool {
	return len(sender) >= 1 && len(sender) <= 20
}

func validateReceiver(receiver string) bool {
	regex := regexp.MustCompile(`^[1-9]\d{9,14}$`)
	return regex.MatchString(receiver)
}
