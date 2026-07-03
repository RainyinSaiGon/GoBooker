package router

import (
	"backend/repository"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type UserDTO struct {
	ID       string `json:"id"`
	UserName string `json:"name"`
	Email    string `json:"email"`
}

type UserRequest struct {
	UserName string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Handler struct {
	userRepo repository.UserRepository
}

func NewRouter(userRepo repository.UserRepository) *Handler {
	return &Handler{userRepo: userRepo}
}

func toDTO(user repository.User) UserDTO {
	return UserDTO{
		UserName: user.Name,
		Email:    user.Email,
	}
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSONResponse(w, status, map[string]string{"error": msg})
}

func (h *Handler) GetAllUserHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepo.GetAllUsers()
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
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.userRepo.GetUserByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve user")
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

	user := repository.User{
		Name:     req.UserName,
		Email:    req.Email,
		Password: req.Password,
	}

	createdUser, err := h.userRepo.CreateUser(user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	writeJSONResponse(w, http.StatusCreated, toDTO(createdUser))
}

func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.userRepo.DeleteUser(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}
	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user := repository.User{
		Name:     req.UserName,
		Email:    req.Email,
		Password: req.Password,
	}

	updatedUser, err := h.userRepo.UpdateUser(id, user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}
	writeJSONResponse(w, http.StatusOK, toDTO(updatedUser))
}
