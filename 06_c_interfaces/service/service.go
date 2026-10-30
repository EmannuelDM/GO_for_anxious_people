package service

import (
	"fmt"
	"strings"
	"time"
)

// UserService coordinates user business logic.
// It depends purely on the UserStore interface, decoupling business logic
// from infrastructure and persistence mechanisms.
type UserService struct {
	store UserStore
}

// NewUserService constructs a UserService with its required dependencies.
// IDIOMATIC PATTERN: "Accept interfaces, return structs"
// - Param 'store' is an interface (UserStore) -> maximum flexibility and testability.
// - Return value '*UserService' is a concrete struct -> allows caller to use all methods without artificial abstraction.
func NewUserService(store UserStore) *UserService {
	if store == nil {
		panic("service: store dependency cannot be nil")
	}
	return &UserService{
		store: store,
	}
}

// Register registers a new user after enforcing domain business rules.
func (s *UserService) Register(id, username, email string) (*User, error) {
	id = strings.TrimSpace(id)
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(strings.ToLower(email))

	if id == "" || username == "" || !strings.Contains(email, "@") {
		return nil, ErrInvalidUserData
	}

	// Rule: Email must be unique
	existing, err := s.store.FindByEmail(email)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("%w: %s", ErrUserAlreadyExists, email)
	}

	user := &User{
		ID:        id,
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
	}

	if err := s.store.Save(user); err != nil {
		return nil, fmt.Errorf("service: failed to save user: %w", err)
	}

	return user, nil
}

// GetUser retrieves a user by ID.
func (s *UserService) GetUser(id string) (*User, error) {
	if strings.TrimSpace(id) == "" {
		return nil, ErrInvalidUserData
	}
	return s.store.FindByID(id)
}
