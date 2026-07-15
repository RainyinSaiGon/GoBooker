package service

import (
	"backend/model"
	"backend/repository"
)

type ConcertService interface {
	GetAllConcerts(page int) ([]model.Concert, error)
	GetAllConcertsByLocation(location string) ([]model.Concert, error)
	CreateConcert(concert model.Concert) (model.Concert, error)
	DeleteConcert(concert model.Concert) error
	UpdateConcert(concert model.Concert) (model.Concert, error)
}

type concertService struct {
	repo repository.ConcertRepository
}

// TODO: Implement the methods of the ConcertService interface using the repository layer.
func NewConcertService(repo repository.ConcertRepository) ConcertService {
	return &concertService{repo: repo}
}

func (s *concertService) GetAllConcerts(page int) ([]model.Concert, error) {
	return s.repo.GetAllConcerts(page)
}

func (s *concertService) GetAllConcertsByLocation(location string) ([]model.Concert, error) {
	return s.repo.GetAllConcertsByLocation(location)
}

func (s *concertService) CreateConcert(concert model.Concert) (model.Concert, error) {
	return s.repo.CreateConcert(concert)
}

func (s *concertService) DeleteConcert(concert model.Concert) error {
	return s.repo.DeleteConcert(concert)
}
func (s *concertService) UpdateConcert(concert model.Concert) (model.Concert, error) {
	return s.repo.UpdateConcert(concert)
}
