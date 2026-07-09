package service_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"backend/model"
	"backend/service"
)

type fakeUserRepository struct {
	usersByID    map[string]model.User
	usersByEmail map[string]model.User
	nextID       int

	createErr error
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		usersByID:    make(map[string]model.User),
		usersByEmail: make(map[string]model.User),
	}
}

func (f *fakeUserRepository) GetAllUsers(query string, limit, offset int) ([]model.User, int, error) {
	var filtered []model.User
	query = strings.ToLower(query)
	for _, u := range f.usersByID {
		if query == "" || strings.Contains(strings.ToLower(u.Name), query) || strings.Contains(strings.ToLower(u.Email), query) {
			filtered = append(filtered, u)
		}
	}
	total := len(filtered)
	if offset >= len(filtered) {
		return []model.User{}, total, nil
	}
	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], total, nil
}

func (f *fakeUserRepository) GetUserByID(id string) (model.User, error) {
	u, ok := f.usersByID[id]
	if !ok {
		return model.User{}, errors.New("not found")
	}
	return u, nil
}

func (f *fakeUserRepository) CreateUser(user model.User) (model.User, error) {
	if f.createErr != nil {
		return model.User{}, f.createErr
	}
	if _, exists := f.usersByEmail[user.Email]; exists {
		return model.User{}, errors.New("email already exists")
	}
	f.nextID++
	user.ID = fmt.Sprintf("%d", f.nextID)
	f.usersByID[user.ID] = user
	f.usersByEmail[user.Email] = user
	return user, nil
}

func (f *fakeUserRepository) DeleteUser(id string) error {
	u, ok := f.usersByID[id]
	if !ok {
		return errors.New("not found")
	}
	delete(f.usersByID, id)
	delete(f.usersByEmail, u.Email)
	return nil
}

func (f *fakeUserRepository) UpdateUser(id string, user model.User) (model.User, error) {
	existing, ok := f.usersByID[id]
	if !ok {
		return model.User{}, errors.New("not found")
	}
	user.ID = existing.ID
	f.usersByID[id] = user
	return user, nil
}

func TestCreateUser(t *testing.T) {
	repo := newFakeUserRepository()
	svc := service.NewUserService(repo)

	testUser := model.User{Name: "Test User", Email: "test@example.com", Password: "password123"}
	got, err := svc.CreateUser(testUser)
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}
	if got.Name != testUser.Name || got.Email != testUser.Email || got.ID == "" {
		t.Errorf("CreateUser() returned incomplete user: %+v", got)
	}
}

func TestGetAllUsers(t *testing.T) {
	repo := newFakeUserRepository()
	svc := service.NewUserService(repo)
	svc.CreateUser(model.User{Name: "U1", Email: "u1@ex.com"})
	svc.CreateUser(model.User{Name: "U2", Email: "u2@ex.com"})

	users, total, err := svc.GetAllUsers("", 1, 10)
	if err != nil {
		t.Fatalf("GetAllUsers() error = %v", err)
	}
	if len(users) != 2 {
		t.Errorf("got %d users, want 2", len(users))
	}
	if total != 2 {
		t.Errorf("got total %d, want 2", total)
	}
}

func TestGetUserByID(t *testing.T) {
	repo := newFakeUserRepository()
	svc := service.NewUserService(repo)
	u, _ := svc.CreateUser(model.User{Name: "Test", Email: "t@ex.com"})

	got, err := svc.GetUserByID(u.ID)
	if err != nil {
		t.Fatalf("GetUserByID() error = %v", err)
	}
	if got.ID != u.ID {
		t.Errorf("got ID %s, want %s", got.ID, u.ID)
	}

	_, err = svc.GetUserByID("nonexistent")
	if err == nil {
		t.Error("expected error for missing user, got nil")
	}
}

func TestUpdateUser(t *testing.T) {
	repo := newFakeUserRepository()
	svc := service.NewUserService(repo)
	u, _ := svc.CreateUser(model.User{Name: "Old", Email: "old@ex.com"})

	updated := model.User{Name: "New", Email: "old@ex.com"}
	got, err := svc.UpdateUser(u.ID, updated)
	if err != nil {
		t.Fatalf("UpdateUser() error = %v", err)
	}
	if got.Name != "New" {
		t.Errorf("got name %s, want New", got.Name)
	}
}

func TestDeleteUser(t *testing.T) {
	repo := newFakeUserRepository()
	svc := service.NewUserService(repo)
	u, _ := svc.CreateUser(model.User{Name: "Del", Email: "del@ex.com"})

	err := svc.DeleteUser(u.ID)
	if err != nil {
		t.Fatalf("DeleteUser() error = %v", err)
	}

	_, err = svc.GetUserByID(u.ID)
	if err == nil {
		t.Error("user still exists after deletion")
	}
}

func TestCreateUser_HashesPassword(t *testing.T) {
	repo := newFakeUserRepository()
	svc := service.NewUserService(repo)
	plaintext := "password123"
	got, _ := svc.CreateUser(model.User{Name: "T", Email: "t@e.com", Password: plaintext})
	if err := bcrypt.CompareHashAndPassword([]byte(got.Password), []byte(plaintext)); err != nil {
		t.Error("password not hashed correctly")
	}
}