package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func openCockroachDB() (*pgxpool.Pool, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = defaultCockroachURL()
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse cockroachdb url: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("connect to cockroachdb: %w", err)
	}

	return pool, nil
}

func defaultCockroachURL() string {
	host := getenvOrDefault("COCKROACH_HOST", "localhost")
	port := getenvOrDefault("COCKROACH_PORT", "26257")
	user := getenvOrDefault("COCKROACH_USER", "root")
	password := os.Getenv("COCKROACH_PASSWORD")
	database := getenvOrDefault("COCKROACH_DATABASE", "defaultdb")
	sslmode := getenvOrDefault("COCKROACH_SSLMODE", "disable")

	if password != "" {
		return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, database, sslmode)
	}

	return fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=%s", user, host, port, database, sslmode)
}

func getenvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}