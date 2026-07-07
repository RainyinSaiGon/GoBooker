package service

import (
	"backend/model"
	"backend/repository"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// ErrNotFound is returned when a requested resource does not exist.
var ErrNotFound = errors.New("resource not found")

type UserService interface {
	GetAllUsers() ([]model.User, error)
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

func (s *userService) GetAllUsers() ([]model.User, error) {
	return s.repo.GetAllUsers()
}

func (s *userService) GetUserByID(id string) (model.User, error) {
	user, err := s.repo.GetUserByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, ErrNotFound
	}
	return user, err
}

func (s *userService) CreateUser(user model.User) (model.User, error) {
	if user.Role == "" {
		user.Role = model.RoleCustomer
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}
	user.Password = string(hashed)
	return s.repo.CreateUser(user)
}

func (s *userService) DeleteUser(id string) error {
	err := s.repo.DeleteUser(id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func (s *userService) UpdateUser(id string, user model.User) (model.User, error) {
	// Only re-hash if a new password was provided.
	if user.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return model.User{}, err
		}
		user.Password = string(hashed)
	}
	updated, err := s.repo.UpdateUser(id, user)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, ErrNotFound
	}
	return updated, err
}
