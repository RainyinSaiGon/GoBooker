package handler

import (
	"backend/model"
	"backend/service"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

// UserDTO is the outward-facing representation of a user (no password).
type UserDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

// UserRequest is the payload accepted on create/update.
type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// UserHandler holds the service dependency for user HTTP handlers.
type UserHandler struct {
	svc service.UserService
}

// NewUserHandler creates a UserHandler.
func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func toUserDTO(u model.User) UserDTO {
	return UserDTO{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.GetAllUsers()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}
	dtos := make([]UserDTO, 0, len(users))
	for _, u := range users {
		dtos = append(dtos, toUserDTO(u))
	}
	writeJSON(w, http.StatusOK, dtos)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	user, err := h.svc.GetUserByID(id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "User not found")
		} else {
			writeError(w, http.StatusInternalServerError, "Failed to retrieve user")
		}
		return
	}
	writeJSON(w, http.StatusOK, toUserDTO(user))
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	created, err := h.svc.CreateUser(model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	writeJSON(w, http.StatusCreated, toUserDTO(created))
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.svc.DeleteUser(id); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	updated, err := h.svc.UpdateUser(id, model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "User not found")
		} else {
			writeError(w, http.StatusInternalServerError, "Failed to update user")
		}
		return
	}
	writeJSON(w, http.StatusOK, toUserDTO(updated))
}
