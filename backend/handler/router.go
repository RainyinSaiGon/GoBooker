package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterUserRoutes registers user CRUD handlers under the assumed prefix ("/users")
func RegisterUserRoutes(r *mux.Router, u *UserHandler) {
	r.HandleFunc("", u.GetAllUsers).Methods(http.MethodGet)
	r.HandleFunc("", u.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/{id}", u.GetUser).Methods(http.MethodGet)
	r.HandleFunc("/{id}", u.UpdateUser).Methods(http.MethodPut)
	r.HandleFunc("/{id}", u.DeleteUser).Methods(http.MethodDelete)
}

// RegisterAuthRoutes registers authentication handlers under the assumed prefix ("/auth")
func RegisterAuthRoutes(r *mux.Router, a *AuthHandler) {
	r.HandleFunc("/login", a.Login).Methods(http.MethodPost)
	r.HandleFunc("/refresh", a.Refresh).Methods(http.MethodPost)
	r.HandleFunc("/logout", a.Logout).Methods(http.MethodPost)
}