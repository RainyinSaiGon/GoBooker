package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

// emailRE is an RFC 5321-compatible email validator.
// It requires a local-part, @, a domain with at least one dot, and a TLD of 2+ chars.
// B2: replaces the naive strings.Contains check that accepted "a@b" or "@foo.com".
var emailRE = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

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

// isValidEmail reports whether email matches a basic RFC 5321 pattern.
func isValidEmail(email string) bool {
	return emailRE.MatchString(email)
}
