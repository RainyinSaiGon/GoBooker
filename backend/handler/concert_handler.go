package handler

import (
	"backend/model"
	"backend/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ConcertHandler struct {
	svc service.ConcertService
}

// TODO: Implement the methods of the ConcertHandler struct using the service layer.
func NewConcertHandler(svc service.ConcertService) *ConcertHandler {
	return &ConcertHandler{svc: svc}
}

func (h *ConcertHandler) GetAllConcerts(w http.ResponseWriter, r *http.Request) {
	// Extract page parameter from the request
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return
	}
	concerts, err := h.svc.GetAllConcerts(pageInt)
	if err != nil {
		return
	}
	writeJSON(w, http.StatusOK, concerts)
}

func (h *ConcertHandler) GetAllConcertsByLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	location := vars["location"]
	concerts, err := h.svc.GetAllConcertsByLocation(location)
	if err != nil {
		return
	}
	writeJSON(w, http.StatusOK, concerts)
}

func (h *ConcertHandler) CreateConcert(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into a Concert struct
	var concert model.Concert
	if err := json.NewDecoder(r.Body).Decode(&concert); err != nil {
		return
	}
	createdConcert, err := h.svc.CreateConcert(concert)
	if err != nil {
		return
	}
	writeJSON(w, http.StatusCreated, createdConcert)
}

func (h *ConcertHandler) DeleteConcert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := h.svc.DeleteConcert(model.Concert{ID: id})
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ConcertHandler) UpdateConcert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	// Parse the request body into a Concert struct
	var concert model.Concert
	if err := json.NewDecoder(r.Body).Decode(&concert); err != nil {
		return
	}
	concert.ID = id
	updatedConcert, err := h.svc.UpdateConcert(concert)
	if err != nil {
		return
	}
	writeJSON(w, http.StatusOK, updatedConcert)
}
