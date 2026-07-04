package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRoutes wires all handlers onto the given mux.Router.
func RegisterRoutes(r *mux.Router, u *UserHandler) {
	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}).Methods(http.MethodGet)

	// Users
	r.HandleFunc("/users", u.GetAllUsers).Methods(http.MethodGet)
	r.HandleFunc("/users", u.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/users/{id}", u.GetUser).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}", u.UpdateUser).Methods(http.MethodPut)
	r.HandleFunc("/users/{id}", u.DeleteUser).Methods(http.MethodDelete)
}
