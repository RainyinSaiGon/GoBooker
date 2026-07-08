package handler

import (
	"backend/model"
	"backend/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

const (
    defaultLimit = 20
	defaultPage  = 1
)

// UserResponse is the public representation of a user returned by the API.
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


// GET /users?q=jane&page=1&limit=10 - Get all users with optional query parameters for filtering and pagination.
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	q := query.Get("q")
	page := query.Get("page")
	limit := query.Get("limit")

	
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = defaultPage
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = defaultLimit;
	}

	users, err := h.svc.GetAllUsers(q, pageInt, limitInt)
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

// GET /users/{id} - Get a user by ID
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
	// #7 limit request body to 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

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

	// #8 — req.Email is guaranteed non-empty here; no redundant check needed
	if !isValidEmail(req.Email) {
		writeError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	created, err := h.svc.CreateUser(model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	writeJSON(w, http.StatusCreated, toUserResponse(created))
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.svc.DeleteUser(id); err != nil {
		// #3 — distinguish not-found from internal errors
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "User not found")
		} else {
			writeError(w, http.StatusInternalServerError, "Failed to delete user")
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// #7 limit request body to 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// #9 — apply the same input validation as CreateUser
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "Name is required")
		return
	}
	if req.Email == "" || !isValidEmail(req.Email) {
		writeError(w, http.StatusBadRequest, "Valid email is required")
		return
	}
	// Password is optional on update; validate only when provided
	if req.Password != "" && len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "Password must be at least 8 characters long")
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
