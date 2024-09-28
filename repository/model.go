package repository

import (
	"github.com/gofrs/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `db:"id"`         // UUID for each user
	Name      string    `db:"name"`       // User's full name
	Email     string    `db:"email"`      // User's email (must be unique)
	Password  string    `db:"password"`   // Hashed password for security
	IsActive  bool      `db:"is_active"`  // Indicates if the user is activated
	CreatedAt time.Time `db:"created_at"` // Timestamp of when the user was created
	UpdatedAt time.Time `db:"updated_at"` // Timestamp of the last update
}
