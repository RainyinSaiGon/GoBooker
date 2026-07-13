package main

import (
	"backend/config"
	"backend/handler"
	"backend/middleware"
	"backend/repository"
	"backend/service"
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("database connection established")

	// Repositories
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)

	// Services
	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(authRepo)

	// Handlers
	userHandler := handler.NewUserHandler(userSvc)
	authHandler := handler.NewAuthHandler(authSvc)

	// Single root router
	r := mux.NewRouter()

	// Public subrouter - no JWT
	publicRouter := r.PathPrefix("/api").Subrouter()
	handler.RegisterAuthRoutes(publicRouter, authHandler)

	// Protected subrouter - JWT applied only here
	protectedRouter := r.PathPrefix("/api").Subrouter()
	protectedRouter.Use(middleware.JWTMiddleware)
	handler.RegisterRoutes(protectedRouter, userHandler)

	// Middleware applied to the whole server (outermost → innermost)
	chain := middleware.CORSMiddleware(cfg.AllowedOrigin)(
		middleware.Recovery(
			middleware.Logger(r),
		),
	)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      chain,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("listening on :%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}