package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"backend/repository"
	"backend/router"

	"github.com/gorilla/mux"
)

func main() {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	router := router.NewRouter(userRepo)

	r := mux.NewRouter()

	// Get all users
	r.HandleFunc("/userID", router.GetAllUserHandler).Methods("GET")

	// Get a specific user
	r.HandleFunc("/userID/{id}", router.GetUserHandler).Methods("GET")

	// Create a new user
	r.HandleFunc("/userID", router.CreateUserHandler).Methods("POST")

	// Delete a user
	r.HandleFunc("/userID/{id}", router.DeleteUserHandler).Methods("DELETE")

	// Update information for a user
	r.HandleFunc("/userID/{id}", router.UpdateUserHandler).Methods("PUT")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
