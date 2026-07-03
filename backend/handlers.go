package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type booking struct {
	ID           string    `json:"id"`
	EventID      string    `json:"eventId"`
	CustomerName string    `json:"customerName"`
	Quantity     int       `json:"quantity"`
	CreatedAt    time.Time `json:"createdAt"`
}

type createBookingRequest struct {
	EventID      string `json:"eventId"`
	CustomerName string `json:"customerName"`
	Quantity     int    `json:"quantity"`
}

func listBookingsHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]any{
		"items": []booking{},
	})
}

func createBookingHandler(w http.ResponseWriter, r *http.Request) {
	var payload createBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON payload",
		})
		return
	}

	if payload.Quantity <= 0 {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "quantity must be greater than zero",
		})
		return
	}

	respondJSON(w, http.StatusCreated, booking{
		ID:           "draft-booking",
		EventID:      payload.EventID,
		CustomerName: payload.CustomerName,
		Quantity:     payload.Quantity,
		CreatedAt:    time.Now().UTC(),
	})
}

func respondJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}