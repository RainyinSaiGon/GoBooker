package config

import (
	"errors"
	"io/fs"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application-level configuration.
type Config struct {
	Port        string
	DatabaseURL string
}

// Load reads configuration from environment variables and returns a Config.
// It fatally logs if any required value is missing.
func Load() Config {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Fatalf("error loading .env file: %v", err)
	}
	if errors.Is(err, fs.ErrNotExist) {
		log.Println(".env file not found, using existing environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	return Config{
		Port:        port,
		DatabaseURL: dbURL,
	}
}
