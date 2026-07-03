package main

import (
	"log"
	"net/http"
)

func main() {
	db, err := openCockroachDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	port := getenvOrDefault("PORT", "8081")
	server := &http.Server{
		Addr:    ":" + port,
		Handler: newRouter(db),
	}

	log.Printf("backend listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}