package service

import (
	"backend/model"
	"backend/repository"
	"database/sql"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// ErrNotFound is returned when a requested resource does not exist.
var ErrNotFound = errors.New("resource not found")

// ErrDuplicateEmail is returned when creating or updating a user with an
// email address that is already registered.
var ErrDuplicateEmail = errors.New("email already in use")

// UserService defines business-logic operations for users.
type UserService interface {
	GetAllUsers(query string, page, size int) ([]model.User, int, error)
	GetUserByID(id string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
	DeleteUser(id string) error
	UpdateUser(id string, user model.User) (model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService returns a UserService backed by the given repository.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// isDuplicateKey reports whether err is a unique-constraint violation from
// CockroachDB or PostgreSQL.
func isDuplicateKey(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "SQLSTATE 23505")
}

// GetAllUsers handles page/size limits and calculates limits and offsets
// before calling the database repository layer.
func (s *userService) GetAllUsers(query string, page, size int) ([]model.User, int, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	offset := (page - 1) * size
	return s.repo.GetAllUsers(query, size, offset)
}

// GetUserByID returns the user with the given ID, or ErrNotFound if no such
// user exists.
func (s *userService) GetUserByID(id string) (model.User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrNotFound
		}
		return model.User{}, err
	}
	return user, nil
}

// CreateUser hashes the password, assigns the default role, and persists the
// user. Returns ErrDuplicateEmail if the email is already taken.
func (s *userService) CreateUser(user model.User) (model.User, error) {
	if user.Role == "" {
		user.Role = model.RoleCustomer
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}
	user.Password = string(hashed)

	created, err := s.repo.CreateUser(user)
	if err != nil {
		if isDuplicateKey(err) {
			return model.User{}, ErrDuplicateEmail
		}
		return model.User{}, err
	}
	return created, nil
}

// DeleteUser removes the user with the given ID. Returns ErrNotFound when no
// such user exists.
func (s *userService) DeleteUser(id string) error {
	err := s.repo.DeleteUser(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

// UpdateUser replaces mutable fields for the given user. A blank Password
// field means "keep the existing hash". Returns ErrNotFound or
// ErrDuplicateEmail as appropriate.
func (s *userService) UpdateUser(id string, user model.User) (model.User, error) {
	if user.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return model.User{}, err
		}
		user.Password = string(hashed)
	}

	updated, err := s.repo.UpdateUser(id, user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrNotFound
		}
		if isDuplicateKey(err) {
			return model.User{}, ErrDuplicateEmail
		}
		return model.User{}, err
	}
	return updated, nil
}
