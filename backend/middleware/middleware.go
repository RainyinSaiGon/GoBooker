package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ContextKey is an unexported type for context keys in this package,
// preventing collisions with keys defined in other packages.
type ContextKey string

const (
	// ContextUserID is the context key for the authenticated user's ID.
	ContextUserID ContextKey = "user_id"
	// ContextEmail is the context key for the authenticated user's email.
	ContextEmail ContextKey = "email"
)

// responseWriter wraps http.ResponseWriter to capture the status code written
// by downstream handlers without double-writing the header.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

// Logger logs the HTTP method, path, status code, and latency for every request.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		log.Printf("%s %s %d %s", r.Method, r.RequestURI, wrapped.status, time.Since(start))
	})
}

// Recovery catches any panic from downstream handlers, logs the stack trace,
// and returns a 500.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// Log the full goroutine stack so the panic is actionable.
				log.Printf("panic recovered: %v\n%s", rec, debug.Stack())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware sets CORS headers. allowedOrigin should come from config
// (e.g. "http://localhost:3000" for dev, the real domain in production).
func CORSMiddleware(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ValidateToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err // signature invalid, expired, or malformed
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// JWTMiddleware returns middleware that validates the JWT access token
// from the Authorization header and stores claims in the request context.
func JWTMiddleware(secret string) func(http.Handler) http.Handler {
	secretBytes := []byte(secret)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the token from the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Malformed Authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := authHeader[len("Bearer "):]

			// Validate the token
			claims, err := ValidateToken(tokenString, secretBytes)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Store claims in context for downstream handlers
			ctx := r.Context()
			if userID, ok := claims["user_id"]; ok {
				ctx = context.WithValue(ctx, ContextUserID, userID)
			}
			if email, ok := claims["email"]; ok {
				ctx = context.WithValue(ctx, ContextEmail, email)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
