package service

import (
	"sync"
)

// MemoryStore is an in-memory, thread-safe implementation of UserStore.
type MemoryStore struct {
	mu    sync.RWMutex
	users map[string]*User // keyed by User.ID
}

// COMPILE-TIME INTERFACE SATISFACTION CHECK:
// This blank identifier declaration forces the Go compiler to verify that
// *MemoryStore implements the UserStore interface.
// If any method in UserStore is missing or has a mismatched signature,
// compilation fails immediately at this line with a clear error message.
// It costs ZERO runtime overhead because it evaluates during compilation.
var _ UserStore = (*MemoryStore)(nil)

// NewMemoryStore constructs a new in-memory user repository.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users: make(map[string]*User),
	}
}

// FindByID retrieves a user by their unique ID.
func (m *MemoryStore) FindByID(id string) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, ok := m.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}

	// Defensive copy to prevent caller mutation of stored record
	copied := *user
	return &copied, nil
}

// FindByEmail searches for a user by email address.
func (m *MemoryStore) FindByEmail(email string) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, user := range m.users {
		if user.Email == email {
			copied := *user
			return &copied, nil
		}
	}
	return nil, ErrUserNotFound
}

// Save inserts or updates a user in memory.
func (m *MemoryStore) Save(user *User) error {
	if user == nil {
		return ErrInvalidUserData
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Clone to ensure internal state isolation
	copied := *user
	m.users[user.ID] = &copied
	return nil
}

// Delete removes a user by ID.
func (m *MemoryStore) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.users[id]; !ok {
		return ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}
