package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// writeJSON serialises data with the given HTTP status code.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("writeJSON encode error: %v", err)
	}
}

// writeError is a convenience wrapper for JSON error responses.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func isValidEmail(email string) bool {
	// Basic email validation: check for presence of "@" and "."
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	return true
}
