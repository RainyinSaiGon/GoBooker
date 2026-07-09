package handler

import (
	"backend/model"
	"backend/service"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// UserResponse is the public representation of a user returned by the API.
// Password is intentionally excluded — it is never serialised to JSON.
type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"` // F2: expose role so the frontend can show the real value
}

// PaginatedUsersResponse is the wrapper for user queries.
type PaginatedUsersResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	Size       int            `json:"size"`
	TotalPages int            `json:"totalPages"`
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
		Role:  u.Role,
	}
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	query := q.Get("query")
	
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	
	size, _ := strconv.Atoi(q.Get("size"))
	if size < 1 {
		size = 10
	}

	users, total, err := h.svc.GetAllUsers(query, page, size)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	responses := make([]UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, toUserResponse(u))
	}

	totalPages := 1
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(size)))
	}

	writeJSON(w, http.StatusOK, PaginatedUsersResponse{
		Users:      responses,
		Total:      total,
		Page:       page,
		Size:       size,
		TotalPages: totalPages,
	})
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
	// Limit request body to 1 MB to prevent DoS via large payloads.
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
		// B3: map duplicate-email to 409 Conflict instead of a generic 500.
		if errors.Is(err, service.ErrDuplicateEmail) {
			writeError(w, http.StatusConflict, "An account with that email already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	writeJSON(w, http.StatusCreated, toUserResponse(created))
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.svc.DeleteUser(id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "User not found")
		} else {
			writeError(w, http.StatusInternalServerError, "Failed to delete user")
		}
		return
	}
	// Q2: 204 No Content is the correct REST response for a successful DELETE.
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Limit request body to 1 MB.
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "Name is required")
		return
	}
	if req.Email == "" || !isValidEmail(req.Email) {
		writeError(w, http.StatusBadRequest, "Valid email is required")
		return
	}
	// Password is optional on update; validate only when provided.
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
		} else if errors.Is(err, service.ErrDuplicateEmail) {
			writeError(w, http.StatusConflict, "An account with that email already exists")
		} else {
			writeError(w, http.StatusInternalServerError, "Failed to update user")
		}
		return
	}
	writeJSON(w, http.StatusOK, toUserResponse(updated))
}
