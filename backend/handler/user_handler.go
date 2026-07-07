package handler

import (
	"backend/model"
	"backend/service"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserRequest is the payload accepted on create/update.
type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserHandler holds the service dependency for user HTTP handlers.
type UserHandler struct {
	svc service.UserService
}

// NewUserHandler creates a UserHandler.
func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func toUserResponse(u model.User) UserResponse {
	return UserResponse{
		Name:  u.Name,
		Email: u.Email,
	}
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.GetAllUsers()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}
	responses := make([]UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, toUserResponse(u))
	}
	writeJSON(w, http.StatusOK, responses)
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
	writeJSON(w, http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "Name, email, and password are required")
		return
	}

	if len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "Password must be at least 8 characters long")
		return
	}

	if req.Email != "" && !isValidEmail(req.Email) {
		writeError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	created, err := h.svc.CreateUser(model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		//log.Printf("CRITICAL DATABASE ERROR: %v", err)
		writeError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	writeJSON(w, http.StatusCreated, toUserResponse(created))
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
	writeJSON(w, http.StatusOK, toUserResponse(updated))
}
