package model

import "time"

// Role constants for User.
const (
	RoleCustomer = "customer"
	RoleAdmin    = "admin"
)

// User represents an application user.
// Password is bcrypt-hashed and must never be serialised to JSON.
type User struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
