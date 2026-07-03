package service

import (
	"backend/repository"
	
)

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
	return s.repo.GetUserByID(id)
}

func (s *userService) CreateUser(user repository.User) (repository.User, error) {
	return s.repo.CreateUser(user)
}

func (s *userService) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}

func (s *userService) UpdateUser(id string, user repository.User) (repository.User, error) {
	return s.repo.UpdateUser(id, user)
}

