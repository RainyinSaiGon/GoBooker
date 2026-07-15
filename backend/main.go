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
	concertRepo := repository.NewConcertRepository(db)

	// Services
	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(authRepo, cfg.JWTSecret, cfg.JWTRefreshSecret)
	concertSvc := service.NewConcertService(concertRepo)

	// Handlers
	userHandler := handler.NewUserHandler(userSvc)
	authHandler := handler.NewAuthHandler(authSvc)
	concertHandler := handler.NewConcertHandler(concertSvc)

	// Single root router
	r := mux.NewRouter()

	// API subrouter
	apiV1 := r.PathPrefix("/api/v1").Subrouter()

	// Public Health endpoint
	apiV1.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)

	// Public Auth subrouter
	authRouter := apiV1.PathPrefix("/auth").Subrouter()
	handler.RegisterAuthRoutes(authRouter, authHandler)

	// Protected Concert subrouter
	concertRouter := apiV1.PathPrefix("/concerts").Subrouter()
    concertRouter.Use(middleware.JWTMiddleware(cfg.JWTSecret)) 
	handler.RegisterConcertRoutes(concertRouter, concertHandler)

	// Protected Users subrouter - JWT applied only here
	userRouter := apiV1.PathPrefix("/users").Subrouter()
	userRouter.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	handler.RegisterUserRoutes(userRouter, userHandler)

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