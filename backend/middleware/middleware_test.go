package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/middleware"
)

// panicHandler is an http.Handler that always panics.
var panicHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	panic("test panic")
})

// okHandler is an http.Handler that always returns 200 OK.
var okHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
})

// ──────────────────────────────────────────────────────────────────────────────
// Recovery middleware
// ──────────────────────────────────────────────────────────────────────────────

func TestRecovery_CatchesPanic(t *testing.T) {
	handler := middleware.Recovery(panicHandler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Recovery status = %d, want 500", rr.Code)
	}
}

func TestRecovery_PassesThrough(t *testing.T) {
	handler := middleware.Recovery(okHandler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Recovery passthrough status = %d, want 200", rr.Code)
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Logger middleware
// ──────────────────────────────────────────────────────────────────────────────

func TestLogger_PassesThrough(t *testing.T) {
	handler := middleware.Logger(okHandler)
	req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Logger passthrough status = %d, want 200", rr.Code)
	}
}

func TestLogger_CapturesNonOKStatus(t *testing.T) {
	notFoundHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	handler := middleware.Logger(notFoundHandler)
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Logger must not swallow the status code.
	if rr.Code != http.StatusNotFound {
		t.Errorf("Logger status = %d, want 404", rr.Code)
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// CORS middleware
// ──────────────────────────────────────────────────────────────────────────────

func TestCORS_SetsHeaders(t *testing.T) {
	handler := middleware.CORSMiddleware("http://localhost:3000")(okHandler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	origin := rr.Header().Get("Access-Control-Allow-Origin")
	if origin == "" {
		t.Error("Access-Control-Allow-Origin header not set")
	}
	methods := rr.Header().Get("Access-Control-Allow-Methods")
	if methods == "" {
		t.Error("Access-Control-Allow-Methods header not set")
	}
	headers := rr.Header().Get("Access-Control-Allow-Headers")
	if headers == "" {
		t.Error("Access-Control-Allow-Headers header not set")
	}
}

func TestCORS_PreflightReturns204(t *testing.T) {
	handler := middleware.CORSMiddleware("http://localhost:3000")(okHandler)
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("OPTIONS preflight status = %d, want 204", rr.Code)
	}
}

func TestCORS_PassesNonPreflightThrough(t *testing.T) {
	handler := middleware.CORSMiddleware("http://localhost:3000")(okHandler)
	req := httptest.NewRequest(http.MethodPost, "/users", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("CORS POST passthrough status = %d, want 200", rr.Code)
	}
}
