package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")

	// ErrInvalidCredentials tries to login with an incorrect email address or password.
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// ErrDuplicateEmail tries signup with an email address that's already in use.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a new User type. Notice how the field names and types align
// with the columns in the database `users` table?
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}
