package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"backend/repository"
	"backend/router"
	"backend/service"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/gorilla/mux"
)

func main() {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Verify the connection is live before accepting traffic.
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	h := router.NewRouter(userSvc)

	r := mux.NewRouter()

	// Users resource
	r.HandleFunc("/users", h.GetAllUserHandler).Methods(http.MethodGet)
	r.HandleFunc("/users", h.CreateUserHandler).Methods(http.MethodPost)
	r.HandleFunc("/users/{id}", h.GetUserHandler).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}", h.UpdateUserHandler).Methods(http.MethodPut)
	r.HandleFunc("/users/{id}", h.DeleteUserHandler).Methods(http.MethodDelete)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
