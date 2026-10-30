package service

import "errors"

var (
	ErrUserNotFound      = errors.New("service: user not found")
	ErrUserAlreadyExists = errors.New("service: user already exists")
	ErrInvalidUserData   = errors.New("service: invalid user data")
)

// UserStore defines the data storage operations required by the UserService.
//
// IDIOMATIC GO PRINCIPLE: "Define interfaces where they are consumed, not where they are implemented."
// Notice that UserStore is defined right here in package 'service' (the consumer),
// NOT in an external database package!
// The service declares ONLY the exact methods it needs. It does not care whether
// the underlying store is PostgreSQL, MongoDB, Redis, or an In-Memory map.
type UserStore interface {
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	Save(user *User) error
	Delete(id string) error
}
