package repository

import (
	"backend/model"
	"database/sql"
	"time"
)

type ConcertRepository interface {
	GetAllConcerts(page int) ([]model.Concert, error)
	GetAllConcertsByLocation(location string) ([]model.Concert, error)
	CreateConcert(concert model.Concert) (model.Concert, error)
	DeleteConcert(concert model.Concert) error
	UpdateConcert(concert model.Concert) (model.Concert, error)
}

type concertRepository struct {
	db *sql.DB
}

func NewConcertRepository(db *sql.DB) ConcertRepository {
	return &concertRepository{db: db}
}

func (r *concertRepository) GetAllConcerts(page int) ([]model.Concert, error) {

	// Get all concerts ordered by date
	rows, err := r.db.Query("SELECT id, name, location, date, start_time, end_time FROM concerts ORDER BY date ASC LIMIT 10 OFFSET $1", (page-1)*10)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan the rows into a slice of Concert structs
	var concerts []model.Concert

	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		var concert model.Concert
		err := rows.Scan(&concert.ID, &concert.Name, &concert.Location, &concert.Date, &concert.StartTime, &concert.EndTime)
		if err != nil {
			return nil, err
		}
		concerts = append(concerts, concert)
	}

	return concerts, nil
}

func (r *concertRepository) GetAllConcertsByLocation(location string) ([]model.Concert, error) {
	// Get all concerts by location ordered by date
	rows, err := r.db.Query("SELECT id, name, location, date, start_time, end_time FROM concerts WHERE location = $1 ORDER BY date ASC", location)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var concerts []model.Concert
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		var concert model.Concert
		err := rows.Scan(&concert.ID, &concert.Name, &concert.Location, &concert.Date, &concert.StartTime, &concert.EndTime)
		if err != nil {
			return nil, err
		}
		concerts = append(concerts, concert)
	}

	return concerts, nil
}

func (r *concertRepository) CreateConcert(concert model.Concert) (model.Concert, error) {

	now := time.Now().UTC()
	err := r.db.QueryRow("INSERT INTO concerts (name, location, date, start_time, end_time, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING name, created_at ON CONFLICT (id) DO NOTHING",
		concert.Name, concert.Location, concert.Date, concert.StartTime, concert.EndTime, now, now).Scan(&concert.Name, &concert.CreatedAt)
	if err != nil {
		return model.Concert{}, err
	}

	return concert, nil
}

func (r *concertRepository) DeleteConcert(concert model.Concert) error {
	_, err := r.db.Exec("DELETE FROM concerts WHERE id = $1", concert.ID)
	return err
}

func (r *concertRepository) UpdateConcert(concert model.Concert) (model.Concert, error) {
	err := r.db.QueryRow("UPDATE concerts SET name = $1, location = $2, date = $3, start_time = $4, end_time = $5, updated_at = now() , version = version + 1 WHERE id = $6 AND version = $7 RETURNING name, updated_at, version",
		concert.Name, concert.Location, concert.Date, concert.StartTime, concert.EndTime, concert.ID, concert.Version).Scan(&concert.Name, &concert.UpdatedAt, &concert.Version)

	if err != nil {
		return model.Concert{}, err
	}

	return concert, nil
}
