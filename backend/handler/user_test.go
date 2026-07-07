package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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

func (m *mockUserService) GetAllUsers() ([]model.User, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.users, nil
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
	// #15: return service.ErrNotFound so the handler can map it to 404
	return service.ErrNotFound
}

// ─── Helper ─────────────────────────────────────────────────────────────────

func newTestRouter(svc service.UserService) *mux.Router {
	r := mux.NewRouter()
	h := handler.NewUserHandler(svc)
	handler.RegisterRoutes(r, h)
	return r
}

// decode into map[string]interface{} to handle mixed types (string + time)
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

// ─── GET /users ──────────────────────────────────────────────────────────────

func TestGetAllUsers_Handler(t *testing.T) {
	svc := &mockUserService{
		users: []model.User{
			{ID: "1", Name: "Alice", Email: "alice@example.com"},
			{ID: "2", Name: "Bob", Email: "bob@example.com"},
		},
	}
	router := newTestRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rr.Code)
	}
	var users []map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&users)
	if len(users) != 2 {
		t.Errorf("got %d users, want 2", len(users))
	}
}

func TestGetAllUsers_ServiceError(t *testing.T) {
	router := newTestRouter(&mockUserService{getErr: errors.New("db error")})
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rr.Code)
	}
}

// ─── POST /users ─────────────────────────────────────────────────────────────

func TestCreateUser_Happy(t *testing.T) {
	router := newTestRouter(&mockUserService{})
	body := bytes.NewBufferString(`{"name":"Alice","email":"alice@example.com","password":"secret123"}`)
	req := httptest.NewRequest(http.MethodPost, "/users", body)
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
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{"email":"a@b.com","password":"secret123"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

func TestCreateUser_MissingEmail(t *testing.T) {
	router := newTestRouter(&mockUserService{})
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{"name":"Alice","password":"secret123"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

func TestCreateUser_MissingPassword(t *testing.T) {
	router := newTestRouter(&mockUserService{})
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{"name":"Alice","email":"a@b.com"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

func TestCreateUser_ShortPassword(t *testing.T) {
	router := newTestRouter(&mockUserService{})
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{"name":"Alice","email":"alice@example.com","password":"short"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

func TestCreateUser_InvalidEmail(t *testing.T) {
	router := newTestRouter(&mockUserService{})
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{"name":"Alice","email":"not-an-email","password":"secret123"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	router := newTestRouter(&mockUserService{})
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

func TestCreateUser_ServiceError(t *testing.T) {
	router := newTestRouter(&mockUserService{createErr: errors.New("db error")})
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{"name":"Alice","email":"alice@example.com","password":"secret123"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rr.Code)
	}
}

// ─── GET /users/{id} ─────────────────────────────────────────────────────────

func TestGetUser_Found(t *testing.T) {
	svc := &mockUserService{
		users: []model.User{{ID: "abc123", Name: "Alice", Email: "alice@example.com"}},
	}
	router := newTestRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/users/abc123", nil)
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
	req := httptest.NewRequest(http.MethodGet, "/users/no-such-id", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rr.Code)
	}
}

// ─── PUT /users/{id} ─────────────────────────────────────────────────────────

func TestUpdateUser_Happy(t *testing.T) {
	svc := &mockUserService{
		users: []model.User{{ID: "abc123", Name: "Alice", Email: "alice@example.com"}},
	}
	router := newTestRouter(svc)
	req := httptest.NewRequest(http.MethodPut, "/users/abc123", bytes.NewBufferString(`{"name":"Alice Updated","email":"alice@example.com","password":"newpass1"}`))
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
	req := httptest.NewRequest(http.MethodPut, "/users/ghost-id", bytes.NewBufferString(`{"name":"X","email":"x@x.com","password":"pass1234"}`))
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
	req := httptest.NewRequest(http.MethodPut, "/users/abc123", bytes.NewBufferString("bad-json"))
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
	req := httptest.NewRequest(http.MethodPut, "/users/abc123", bytes.NewBufferString(`{"name":"Alice","email":"alice@example.com","password":"short"}`))
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
	req := httptest.NewRequest(http.MethodPut, "/users/abc123", bytes.NewBufferString(`{"email":"alice@example.com","password":"newpass1"}`))
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
	req := httptest.NewRequest(http.MethodPut, "/users/abc123", bytes.NewBufferString(`{"name":"Alice","email":"bademail","password":"newpass1"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rr.Code)
	}
}

// ─── DELETE /users/{id} ──────────────────────────────────────────────────────

func TestDeleteUser_Happy(t *testing.T) {
	svc := &mockUserService{
		users: []model.User{{ID: "abc123", Name: "Alice", Email: "alice@example.com"}},
	}
	router := newTestRouter(svc)
	req := httptest.NewRequest(http.MethodDelete, "/users/abc123", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rr.Code)
	}
}

// #3 — deleting a non-existent user must return 404, not 200 or 500
func TestDeleteUser_NotFound(t *testing.T) {
	router := newTestRouter(&mockUserService{})
	req := httptest.NewRequest(http.MethodDelete, "/users/no-such-id", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rr.Code)
	}
}

func TestDeleteUser_ServiceError(t *testing.T) {
	router := newTestRouter(&mockUserService{deleteErr: errors.New("db error")})
	req := httptest.NewRequest(http.MethodDelete, "/users/abc123", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rr.Code)
	}
}
