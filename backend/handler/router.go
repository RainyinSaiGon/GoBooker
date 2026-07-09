package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRoutes wires all handlers onto the given mux.Router.
func RegisterRoutes(r *mux.Router, u *UserHandler) {
	// Health check — lives at root so load-balancers can probe it without
	// needing to know the API version.
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}).Methods(http.MethodGet)

	// B8: versioned sub-router — all resource routes live under /api/v1.
	api := r.PathPrefix("/api/v1").Subrouter()

	// Users
	api.HandleFunc("/users", u.GetAllUsers).Methods(http.MethodGet)
	api.HandleFunc("/users", u.CreateUser).Methods(http.MethodPost)
	api.HandleFunc("/users/{id}", u.GetUser).Methods(http.MethodGet)
	api.HandleFunc("/users/{id}", u.UpdateUser).Methods(http.MethodPut)
	api.HandleFunc("/users/{id}", u.DeleteUser).Methods(http.MethodDelete)
}
