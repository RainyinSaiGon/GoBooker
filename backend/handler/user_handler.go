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
type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
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

// maxBodySize limits the size of incoming request bodies (1 MB).
const maxBodySize = 1 << 20

func toUserResponse(u model.User) UserResponse {
	return UserResponse{
		ID:    u.ID,
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
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}


	created, err := h.svc.CreateUser(model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrDuplicateEmail):
			writeError(w, http.StatusConflict, "An account with that email already exists")
		case errors.Is(err, service.ErrInvalidName),
			errors.Is(err, service.ErrInvalidEmail),
			errors.Is(err, service.ErrWeakPassword):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "Failed to create user")
		}
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

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
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
		switch {
		case errors.Is(err, service.ErrNotFound):
			writeError(w, http.StatusNotFound, "User not found")
		case errors.Is(err, service.ErrDuplicateEmail):
			writeError(w, http.StatusConflict, "An account with that email already exists")
		case errors.Is(err, service.ErrInvalidName),
			errors.Is(err, service.ErrInvalidEmail),
			errors.Is(err, service.ErrWeakPassword):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "Failed to update user")
		}
		return
	}
	writeJSON(w, http.StatusOK, toUserResponse(updated))
}
