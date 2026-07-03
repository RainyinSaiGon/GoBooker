package service

import (
	"backend/repository"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// ErrNotFound is returned by the service when a requested resource does not exist.
var ErrNotFound = errors.New("user not found")

type UserService interface {
	GetAllUsers() ([]repository.User, error)
	GetUserByID(id string) (repository.User, error)
	CreateUser(user repository.User) (repository.User, error)
	DeleteUser(id string) error
	UpdateUser(id string, user repository.User) (repository.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAllUsers() ([]repository.User, error) {
	return s.repo.GetAllUsers()
}

func (s *userService) GetUserByID(id string) (repository.User, error) {
	user, err := s.repo.GetUserByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		return repository.User{}, ErrNotFound
	}
	return user, err
}

func (s *userService) CreateUser(user repository.User) (repository.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return repository.User{}, err
	}
	user.Password = string(hashed)
	return s.repo.CreateUser(user)
}

func (s *userService) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}

func (s *userService) UpdateUser(id string, user repository.User) (repository.User, error) {
	// Only re-hash if a new password is provided.
	if user.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return repository.User{}, err
		}
		user.Password = string(hashed)
	}
	return s.repo.UpdateUser(id, user)
}
