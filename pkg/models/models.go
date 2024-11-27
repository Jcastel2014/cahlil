package models

import "errors"

var (
	ErrRecordNotFound     = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateUsername  = errors.New("models: duplicate username")
)

type User struct {
	ID       int
	Username string
	Password string
}
