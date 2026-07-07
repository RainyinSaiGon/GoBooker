package service_test

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"backend/model"
	"backend/service"
)


type fakeUserRepository struct {
	usersByID    map[string]model.User
	usersByEmail map[string]model.User
	nextID       int

	// Optional error injection for testing failure paths.
	createErr error
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		usersByID:    make(map[string]model.User),
		usersByEmail: make(map[string]model.User),
	}
}

func (f *fakeUserRepository) GetAllUsers() ([]model.User, error) {
	all := make([]model.User, 0, len(f.usersByID))
	for _, u := range f.usersByID {
		all = append(all, u)
	}
	return all, nil
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
	user.ID = string(rune('0' + f.nextID)) // simple deterministic fake ID
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

	testUser := model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	got, err := svc.CreateUser(testUser)
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	if got.Name != testUser.Name {
		t.Errorf("Name = %q, want %q", got.Name, testUser.Name)
	}
	if got.Email != testUser.Email {
		t.Errorf("Email = %q, want %q", got.Email, testUser.Email)
	}
	if got.ID == "" {
		t.Error("CreateUser() did not assign an ID")
	}
}

func TestCreateUser_DefaultsToCustomerRole(t *testing.T) {
	repo := newFakeUserRepository()
	svc := service.NewUserService(repo)

	got, err := svc.CreateUser(model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	if got.Role != model.RoleCustomer {
		t.Errorf("Role = %q, want %q", got.Role, model.RoleCustomer)
	}
}

func TestCreateUser_PreservesExplicitRole(t *testing.T) {
	repo := newFakeUserRepository()
	svc := service.NewUserService(repo)

	got, err := svc.CreateUser(model.User{
		Name:     "Admin User",
		Email:    "admin@example.com",
		Password: "password123",
		Role:     model.RoleAdmin,
	})
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	if got.Role != model.RoleAdmin {
		t.Errorf("Role = %q, want %q", got.Role, model.RoleAdmin)
	}
}

func TestCreateUser_HashesPassword(t *testing.T) {
	repo := newFakeUserRepository()
	svc := service.NewUserService(repo)

	plaintext := "password123"
	got, err := svc.CreateUser(model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: plaintext,
	})
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	if got.Password == plaintext {
		t.Fatal("CreateUser() stored password in plaintext, want bcrypt hash")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(got.Password), []byte(plaintext)); err != nil {
		t.Errorf("stored password is not a valid bcrypt hash of the original: %v", err)
	}
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	repo := newFakeUserRepository()
	svc := service.NewUserService(repo)

	testUser := model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	if _, err := svc.CreateUser(testUser); err != nil {
		t.Fatalf("first CreateUser() error = %v", err)
	}

	if _, err := svc.CreateUser(testUser); err == nil {
		t.Error("expected error on duplicate email, got nil")
	}
}

func TestCreateUser_RepositoryError(t *testing.T) {
	repo := newFakeUserRepository()
	repo.createErr = errors.New("db connection lost")
	svc := service.NewUserService(repo)

	_, err := svc.CreateUser(model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	})

	if err == nil {
		t.Fatal("expected error when repository fails, got nil")
	}
}

func TestCreateUser_EmptyPasswordStillHashes(t *testing.T) {
	repo := newFakeUserRepository()
	svc := service.NewUserService(repo)

	got, err := svc.CreateUser(model.User{
		Name:  "Test User",
		Email: "test@example.com",
	})
	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}
	if got.Password == "" {
		t.Error("expected a hashed value even for empty input password")
	}
}