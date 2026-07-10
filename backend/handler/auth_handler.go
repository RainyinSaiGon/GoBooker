package handler

import (
	"backend/service"
	"net/http"
	"encoding/json"
	"backend/model"
)

type AuthHandler struct {
	svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	token, refreshToken, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			writeError(w, http.StatusUnauthorized, "Invalid email or password")
		} else {
			writeError(w, http.StatusInternalServerError, "Failed to login")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"token":       token,
		"refreshToken": refreshToken,
	})
}



