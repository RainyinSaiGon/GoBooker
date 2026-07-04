package config

import (
	"log"
	"os"
)

// Config holds all application-level configuration.
type Config struct {
	Port        string
	DatabaseURL string
}

// Load reads configuration from environment variables and returns a Config.
// It fatally logs if any required value is missing.
func Load() Config {
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
