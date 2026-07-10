package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRoutes wires all handlers onto the given mux.Router.
func RegisterRoutes(r *mux.Router, u *UserHandler) {
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}).Methods(http.MethodGet)

	// API subrouter
	api := r.PathPrefix("/api/v1").Subrouter()

	// Usersx
	api.HandleFunc("/users", u.GetAllUsers).Methods(http.MethodGet)
	api.HandleFunc("/users", u.CreateUser).Methods(http.MethodPost)
	api.HandleFunc("/users/{id}", u.GetUser).Methods(http.MethodGet)
	api.HandleFunc("/users/{id}", u.UpdateUser).Methods(http.MethodPut)
	api.HandleFunc("/users/{id}", u.DeleteUser).Methods(http.MethodDelete)


}

func RegisterAuthRoutes(r *mux.Router, a *AuthHandler) {
	auth := r.PathPrefix("/api/v1").Subrouter()
	auth.HandleFunc("/auth/login", a.Login).Methods(http.MethodPost)
}