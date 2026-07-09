package repository_test

import (
	"database/sql"
	"os"
	"testing"

	"backend/model"
	"backend/repository"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// openTestDB opens a database connection from the TEST_DATABASE_URL environment
// variable.  If the variable is not set, the test is skipped so that the suite
// can run without a live database (e.g. in unit-only CI).
func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set – skipping integration tests")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatalf("db.Ping: %v", err)
	}
	return db
}

// seedUser inserts a user directly and registers a cleanup to delete it.
func seedUser(t *testing.T, repo repository.UserRepository, u model.User) model.User {
	t.Helper()
	created, err := repo.CreateUser(u)
	if err != nil {
		t.Fatalf("seedUser CreateUser: %v", err)
	}
	t.Cleanup(func() { repo.DeleteUser(created.ID) })
	return created
}

func TestUserRepository_CreateAndGetByID(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	repo := repository.NewUserRepository(db)

	u := model.User{
		Name:     "Integration User",
		Email:    "integ@example.com",
		Password: "$2a$10$fakehash",
		Role:     model.RoleCustomer,
	}
	created := seedUser(t, repo, u)

	if created.ID == "" {
		t.Fatal("CreateUser did not return an ID")
	}
	if created.Name != u.Name {
		t.Errorf("Name = %q, want %q", created.Name, u.Name)
	}

	got, err := repo.GetUserByID(created.ID)
	if err != nil {
		t.Fatalf("GetUserByID: %v", err)
	}
	if got.Email != u.Email {
		t.Errorf("Email = %q, want %q", got.Email, u.Email)
	}
}

func TestUserRepository_GetAllUsers(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	repo := repository.NewUserRepository(db)

	seedUser(t, repo, model.User{Name: "A", Email: "a_integ@example.com", Password: "hash", Role: model.RoleCustomer})
	seedUser(t, repo, model.User{Name: "B", Email: "b_integ@example.com", Password: "hash", Role: model.RoleCustomer})

	users, total, err := repo.GetAllUsers("", 10, 0)
	if err != nil {
		t.Fatalf("GetAllUsers: %v", err)
	}
	if len(users) < 2 {
		t.Errorf("GetAllUsers returned %d users, want at least 2", len(users))
	}
	if total < 2 {
		t.Errorf("GetAllUsers total = %d, want at least 2", total)
	}
}

func TestUserRepository_GetUserByID_NotFound(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	repo := repository.NewUserRepository(db)

	_, err := repo.GetUserByID("00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Fatal("expected error for non-existent user, got nil")
	}
}

func TestUserRepository_UpdateUser(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	repo := repository.NewUserRepository(db)

	created := seedUser(t, repo, model.User{
		Name:     "Before",
		Email:    "update_integ@example.com",
		Password: "hash",
		Role:     model.RoleCustomer,
	})

	updated, err := repo.UpdateUser(created.ID, model.User{
		Name:     "After",
		Password: "newhash",
	})
	if err != nil {
		t.Fatalf("UpdateUser: %v", err)
	}
	if updated.Name != "After" {
		t.Errorf("Name = %q, want After", updated.Name)
	}
}

func TestUserRepository_DeleteUser(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()
	repo := repository.NewUserRepository(db)

	u := model.User{
		Name:     "ToDelete",
		Email:    "delete_integ@example.com",
		Password: "hash",
		Role:     model.RoleCustomer,
	}
	created, err := repo.CreateUser(u)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	if err := repo.DeleteUser(created.ID); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}

	_, err = repo.GetUserByID(created.ID)
	if err == nil {
		t.Error("user still exists after DeleteUser")
	}
}
