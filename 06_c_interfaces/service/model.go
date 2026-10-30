package service

import "time"

// User represents a user entity in our business domain.
type User struct {
	ID        string
	Username  string
	Email     string
	CreatedAt time.Time
}
