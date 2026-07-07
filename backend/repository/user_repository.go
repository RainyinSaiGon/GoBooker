package repository

import (
	"backend/model"
	"database/sql"
	"time"
)

// UserRepository defines data access operations for users.
type UserRepository interface {
	GetAllUsers() ([]model.User, error)
	GetUserByID(id string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
	DeleteUser(id string) error
	UpdateUser(id string, user model.User) (model.User, error)
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository returns a UserRepository backed by the given *sql.DB.
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAllUsers() ([]model.User, error) {
	rows, err := r.db.Query(
		`SELECT id, email, name, role, created_at, updated_at FROM users`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *userRepository) GetUserByID(id string) (model.User, error) {
	var u model.User
	err := r.db.QueryRow(
		`SELECT id, email, name, role, created_at, updated_at FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.Name, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func (r *userRepository) CreateUser(u model.User) (model.User, error) {
	now := time.Now().UTC()
	err := r.db.QueryRow(
		`INSERT INTO users (id, email, name, password, role, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $5)
		 RETURNING id, created_at, updated_at`,
		u.Email, u.Name, u.Password, u.Role, now,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func (r *userRepository) DeleteUser(id string) error {
	res, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *userRepository) UpdateUser(id string, u model.User) (model.User, error) {
	now := time.Now().UTC()
	err := r.db.QueryRow(
		`UPDATE users SET name = $1, email = $2, password = $3, updated_at = $4
		 WHERE id = $5
		 RETURNING id, email, name, role, created_at, updated_at`,
		u.Name, u.Email, u.Password, now, id,
	).Scan(&u.ID, &u.Email, &u.Name, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}
