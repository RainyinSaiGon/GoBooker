package repository

import (
	"database/sql"
	"backend/model"
)

type AuthRepository interface {
	GetUserByEmail(email string) (model.User, error)
}

type authRepository struct {
	db *sql.DB
}

func  NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

// Find user by email and return the user details.
func (r *authRepository) GetUserByEmail(email string) (model.User, error) {
	var u model.User
	err := r.db.QueryRow("SELECT id, email, password, role FROM users WHERE email = $1", email).Scan(&u.ID, &u.Email, &u.Password, &u.Role)
	if err != nil {
		return model.User{}, err
	}
	return u, nil
}

