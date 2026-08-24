package common

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteCustomError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": msg}); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func WriteError(w http.ResponseWriter, err error) {
	status, msg := MapError(err)
	WriteCustomError(w, status, msg)
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if status == http.StatusNoContent {
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
