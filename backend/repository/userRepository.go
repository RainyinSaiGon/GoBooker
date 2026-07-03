package repository

import (
	"database/sql"
)

type User struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserRepository interface {
	GetAllUsers() ([]User, error)
	GetUserByID(id string) (User, error)
	CreateUser(user User) (User, error)
	DeleteUser(id string) error
	UpdateUser(id string, user User) (User, error)
}

// userRepository implements UserRepository using a SQL database connection.
type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAllUsers() ([]User, error) {
	rows, err := r.db.Query("SELECT email, name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Email, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) GetUserByID(id string) (User, error) {
	var user User
	err := r.db.QueryRow("SELECT email, name FROM users WHERE email = $1", id).Scan(&user.Email, &user.Name)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *userRepository) CreateUser(user User) (User, error) {
	_, err := r.db.Exec("INSERT INTO users (email, name, password) VALUES ($1, $2, $3)", user.Email, user.Name, user.Password)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *userRepository) DeleteUser(id string) error {
	_, err := r.db.Exec("DELETE FROM users WHERE email = $1", id)
	return err
}

func (r *userRepository) UpdateUser(id string, user User) (User, error) {
	_, err := r.db.Exec("UPDATE users SET name = $1, password = $2 WHERE email = $3", user.Name, user.Password, id)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
