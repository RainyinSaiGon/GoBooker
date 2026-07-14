package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend/handler"
	"backend/model"
	"backend/service"

	"github.com/gorilla/mux"
)

// ─── Mock UserService ────────────────────────────────────────────────────────

type mockUserService struct {
	users     []model.User
	createErr error
	getErr    error
	updateErr error
	deleteErr error
}

func (m *mockUserService) GetAllUsers(query string, page, size int) ([]model.User, int, error) {
	if m.getErr != nil {
		return nil, 0, m.getErr
	}
	var filtered []model.User
	query = strings.ToLower(query)
	for _, u := range m.users {
		if query == "" || strings.Contains(strings.ToLower(u.Name), query) || strings.Contains(strings.ToLower(u.Email), query) {
			filtered = append(filtered, u)
		}
	}
	total := len(filtered)
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	offset := (page - 1) * size
	if offset >= len(filtered) {
		return []model.User{}, total, nil
	}
	end := offset + size
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], total, nil
}

func (m *mockUserService) GetUserByID(id string) (model.User, error) {
	if m.getErr != nil {
		return model.User{}, m.getErr
	}
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return model.User{}, service.ErrNotFound
}

func (m *mockUserService) CreateUser(u model.User) (model.User, error) {
	if m.createErr != nil {
		return model.User{}, m.createErr
	}
	u.ID = "generated-id"
	m.users = append(m.users, u)
	return u, nil
}

func (m *mockUserService) UpdateUser(id string, u model.User) (model.User, error) {
	if m.updateErr != nil {
		return model.User{}, m.updateErr
	}
	for i, existing := range m.users {
		if existing.ID == id {
			u.ID = id
			m.users[i] = u
			return u, nil
		}
	}
	return model.User{}, service.ErrNotFound
}

func (m *mockUserService) DeleteUser(id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	for i, u := range m.users {
		if u.ID == id {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return service.ErrNotFound
}

func newTestRouter(svc service.UserService) *mux.Router {
	r := mux.NewRouter()
	h := handler.NewUserHandler(svc)
	
	// Register health check on root
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)

	// Register user routes with the /api/v1/users prefix
	userRouter := r.PathPrefix("/api/v1/users").Subrouter()
	handler.RegisterUserRoutes(userRouter, h)
	
	return r
}

// decodeBody decodes into map[string]interface{} to handle mixed types.
func decodeBody(t *testing.T, rr *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&m); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	return m
}

// ─── GET /health ─────────────────────────────────────────────────────────────

func TestHealthCheck(t *testing.T) {
	router := newTestRouter(&mockUserService{})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("GET /health status = %d, want 200", rr.Code)
	}
	body := decodeBody(t, rr)
	if body["status"] != "ok" {
		t.Errorf("status = %q, want ok", body["status"])
	}
}

// ─── GET /api/v1/users ───────────────────────────────────────────────────────

func TestGetAllUsers_Handler(t *testing.T) {
	svc := &mockUserService{
		users: []model.User{
			{ID: "1", Name: "Alice", Email: "alice@example.com"},
			{ID: "2", Name: "Bob", Email: "bob@example.com"},
		},
	}
	router := newTestRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rr.Code)
	}
	var resp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&resp)
	users, ok := resp["users"].([]interface{})
	if !ok {
		t.Fatalf("expected users array in response, got: %v", resp)
	}
	if len(users) != 2 {
		t.Errorf("got %d users, want 2", len(users))
	}
	if int(resp["total"].(float64)) != 2 {
		t.Errorf("got total %v, want 2", resp["total"])
	}
}

func TestGetAllUsers_ServiceError(t *testing.T) {
	router := newTestRouter(&mockUserService{getErr: errors.New("db error")})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rr.Code)
	}
}

// ─── POST /api/v1/users ──────────────────────────────────────────────────────

func TestCreateUser_Happy(t *testing.T) {
	router := newTestRouter(&mockUserService{})
	body := bytes.NewBufferString(`{"name":"Alice","email":"alice@example.com","password":"secret123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("status = %d, want 201; body = %s", rr.Code, rr.Body.String())
	}
	resp := decodeBody(t, rr)
	if resp["name"] != "Alice" {
		t.Errorf("name = %q, want Alice", resp["name"])
	}
}

func TestCreateUser_MissingName(t *testing.T) {
	router := newTestRouter(&mockUserService{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(`{"email":"a@b.com","password":"secret123"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

func TestCreateUser_MissingEmail(t *testing.T) {
	router := newTestRouter(&mockUserService{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(`{"name":"Alice","password":"secret123"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

func TestCreateUser_MissingPassword(t *testing.T) {
	router := newTestRouter(&mockUserService{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(`{"name":"Alice","email":"a@b.com"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

func TestCreateUser_ShortPassword(t *testing.T) {
	router := newTestRouter(&mockUserService{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(`{"name":"Alice","email":"alice@example.com","password":"short"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

func TestCreateUser_InvalidEmail(t *testing.T) {
	router := newTestRouter(&mockUserService{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(`{"name":"Alice","email":"not-an-email","password":"secret123"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	router := newTestRouter(&mockUserService{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

func TestCreateUser_ServiceError(t *testing.T) {
	router := newTestRouter(&mockUserService{createErr: errors.New("db error")})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(`{"name":"Alice","email":"alice@example.com","password":"secret123"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rr.Code)
	}
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	router := newTestRouter(&mockUserService{createErr: service.ErrDuplicateEmail})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(`{"name":"Alice","email":"alice@example.com","password":"secret123"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusConflict {
		t.Errorf("status = %d, want 409", rr.Code)
	}
}

// ─── GET /api/v1/users/{id} ──────────────────────────────────────────────────

func TestGetUser_Found(t *testing.T) {
	svc := &mockUserService{
		users: []model.User{{ID: "abc123", Name: "Alice", Email: "alice@example.com"}},
	}
	router := newTestRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/abc123", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rr.Code)
	}
	resp := decodeBody(t, rr)
	if resp["name"] != "Alice" {
		t.Errorf("name = %q, want Alice", resp["name"])
	}
}

func TestGetUser_NotFound(t *testing.T) {
	router := newTestRouter(&mockUserService{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/no-such-id", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rr.Code)
	}
}

// ─── PUT /api/v1/users/{id} ──────────────────────────────────────────────────

func TestUpdateUser_Happy(t *testing.T) {
	svc := &mockUserService{
		users: []model.User{{ID: "abc123", Name: "Alice", Email: "alice@example.com"}},
	}
	router := newTestRouter(svc)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/abc123", bytes.NewBufferString(`{"name":"Alice Updated","email":"alice@example.com","password":"newpass1"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want 200; body = %s", rr.Code, rr.Body.String())
	}
	resp := decodeBody(t, rr)
	if resp["name"] != "Alice Updated" {
		t.Errorf("name = %q, want Alice Updated", resp["name"])
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	router := newTestRouter(&mockUserService{})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/ghost-id", bytes.NewBufferString(`{"name":"X","email":"x@x.com","password":"pass1234"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rr.Code)
	}
}

func TestUpdateUser_InvalidJSON(t *testing.T) {
	svc := &mockUserService{
		users: []model.User{{ID: "abc123", Name: "Alice", Email: "alice@example.com"}},
	}
	router := newTestRouter(svc)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/abc123", bytes.NewBufferString("bad-json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

func TestUpdateUser_ShortPassword(t *testing.T) {
	svc := &mockUserService{
		users: []model.User{{ID: "abc123", Name: "Alice", Email: "alice@example.com"}},
	}
	router := newTestRouter(svc)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/abc123", bytes.NewBufferString(`{"name":"Alice","email":"alice@example.com","password":"short"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

func TestUpdateUser_MissingName(t *testing.T) {
	svc := &mockUserService{
		users: []model.User{{ID: "abc123", Name: "Alice", Email: "alice@example.com"}},
	}
	router := newTestRouter(svc)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/abc123", bytes.NewBufferString(`{"email":"alice@example.com","password":"newpass1"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

func TestUpdateUser_InvalidEmail(t *testing.T) {
	svc := &mockUserService{
		users: []model.User{{ID: "abc123", Name: "Alice", Email: "alice@example.com"}},
	}
	router := newTestRouter(svc)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/abc123", bytes.NewBufferString(`{"name":"Alice","email":"bademail","password":"newpass1"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

// ─── DELETE /api/v1/users/{id} ───────────────────────────────────────────────

func TestDeleteUser_Happy(t *testing.T) {
	svc := &mockUserService{
		users: []model.User{{ID: "abc123", Name: "Alice", Email: "alice@example.com"}},
	}
	router := newTestRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/abc123", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	// Q2: DELETE now returns 204 No Content.
	if rr.Code != http.StatusNoContent {
		t.Errorf("status = %d, want 204", rr.Code)
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	router := newTestRouter(&mockUserService{})
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/no-such-id", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rr.Code)
	}
}

func TestDeleteUser_ServiceError(t *testing.T) {
	router := newTestRouter(&mockUserService{deleteErr: errors.New("db error")})
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/abc123", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rr.Code)
	}
}
