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

	// Connection pool tuning
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Repositories
	userRepo := repository.NewUserRepository(db)

	// Services
	userSvc := service.NewUserService(userRepo)

	// Handlers
	userHandler := handler.NewUserHandler(userSvc)

	// Router
	r := mux.NewRouter()
	handler.RegisterRoutes(r, userHandler)

	// Middleware (applied outermost → innermost)

	chain := middleware.CORSMiddleware(cfg.AllowedOrigin)(middleware.Recovery(middleware.Logger(r)))

	log.Printf("listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, chain); err != nil {
		log.Fatal(err)
	}
}
