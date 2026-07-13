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

	// Services
	userSvc := service.NewUserService(userRepo)

	// Handlers
	userHandler := handler.NewUserHandler(userSvc)

	authRepo := repository.NewAuthRepository(db)
	authSvc := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authSvc)

	// Router
	r := mux.NewRouter()
	handler.RegisterRoutes(r, userHandler)
	handler.RegisterAuthRoutes(r, authHandler)

	// Middleware (applied outermost → innermost)
	chain := middleware.CORSMiddleware(cfg.AllowedOrigin)(
	middleware.Recovery(
		middleware.Logger(
			middleware.JWTMiddleware(r),
		),
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