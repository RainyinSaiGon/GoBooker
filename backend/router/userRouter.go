package router

import (
	"backend/repository"
	"backend/service"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

// UserDTO is the outward-facing representation of a user (no password).
type UserDTO struct {
	ID       string `json:"id"`
	UserName string `json:"name"`
	Email    string `json:"email"`
}

// UserRequest is the payload accepted on create/update.
type UserRequest struct {
	UserName string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Handler holds the service dependency for all user-related HTTP handlers.
type Handler struct {
	userSvc service.UserService
}

// NewRouter wires a Handler around the given UserService.
func NewRouter(userSvc service.UserService) *Handler {
	return &Handler{userSvc: userSvc}
}

// toDTO maps a repository.User to a UserDTO.
// Email doubles as the public identifier since no separate UUID exists yet.
func toDTO(u repository.User) UserDTO {
	return UserDTO{
		ID:       u.Email,
		UserName: u.Name,
		Email:    u.Email,
	}
}

// writeJSONResponse serialises data and sends it with the given status code.
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeError is a convenience wrapper for JSON error responses.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSONResponse(w, status, map[string]string{"error": msg})
}

func (h *Handler) GetAllUserHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.userSvc.GetAllUsers()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	dtos := make([]UserDTO, 0, len(users))
	for _, u := range users {
		dtos = append(dtos, toDTO(u))
	}
	writeJSONResponse(w, http.StatusOK, dtos)
}

func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	user, err := h.userSvc.GetUserByID(id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "User not found")
		} else {
			writeError(w, http.StatusInternalServerError, "Failed to retrieve user")
		}
		return
	}
	writeJSONResponse(w, http.StatusOK, toDTO(user))
}

func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	created, err := h.userSvc.CreateUser(repository.User{
		Name:     req.UserName,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	writeJSONResponse(w, http.StatusCreated, toDTO(created))
}

func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.userSvc.DeleteUser(id); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}
	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updated, err := h.userSvc.UpdateUser(id, repository.User{
		Name:     req.UserName,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}
	writeJSONResponse(w, http.StatusOK, toDTO(updated))
}
