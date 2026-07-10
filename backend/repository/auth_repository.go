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

// Find user by email and return the user ID and hashed password.
func (r *authRepository) GetUserByEmail(email string) (model.User, error) {
	var id, password string
	err := r.db.QueryRow("SELECT id, password FROM users WHERE email = $1", email).Scan(&id, &password)
	if err != nil {
		return model.User{}, err
	}
	return model.User{ID: id, Password: password}, nil
}

