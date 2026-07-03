package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func newRouter(db *pgxpool.Pool) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))

	router.Get("/health", healthHandler)
	router.Get("/healthz", healthHandler)
	router.Get("/readyz", readyHandler(db))

	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/bookings", listBookingsHandler)
		r.Post("/bookings", createBookingHandler)
	})

	return router
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "GoBooker backend is running",
	})
}

func readyHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			respondJSON(w, http.StatusServiceUnavailable, map[string]string{
				"status":  "degraded",
				"message": "database unavailable",
			})
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"status":  "ready",
			"message": "backend is healthy",
		})
	}
}