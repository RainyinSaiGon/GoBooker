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
	ID        string    `db:"id"         json:"id"`
	Name      string    `db:"name"        json:"name"`
	Email     string    `db:"email"       json:"email"`
	Password  string    `db:"password"    json:"-"`
	Role      string    `db:"role"        json:"role"`
	CreatedAt time.Time `db:"created_at"  json:"created_at"`
	UpdatedAt time.Time `db:"updated_at"  json:"updated_at"`
}
