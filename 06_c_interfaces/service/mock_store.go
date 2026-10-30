package service

// MockStore is a mock implementation of UserStore for testing and simulation.
// By defining interfaces, callers can mock dependencies cleanly without
// external third-party mocking libraries.
type MockStore struct {
	FindByIDFn    func(id string) (*User, error)
	FindByEmailFn func(email string) (*User, error)
	SaveFn        func(user *User) error
	DeleteFn      func(id string) error

	// Call counters for verifying mock interactions
	SavedUsers []*User
}

// Compile-time check that MockStore implements UserStore
var _ UserStore = (*MockStore)(nil)

// NewMockStore initializes a MockStore with default no-op behavior.
func NewMockStore() *MockStore {
	return &MockStore{
		SavedUsers: make([]*User, 0),
	}
}

func (m *MockStore) FindByID(id string) (*User, error) {
	if m.FindByIDFn != nil {
		return m.FindByIDFn(id)
	}
	return nil, ErrUserNotFound
}

func (m *MockStore) FindByEmail(email string) (*User, error) {
	if m.FindByEmailFn != nil {
		return m.FindByEmailFn(email)
	}
	return nil, ErrUserNotFound
}

func (m *MockStore) Save(user *User) error {
	m.SavedUsers = append(m.SavedUsers, user)
	if m.SaveFn != nil {
		return m.SaveFn(user)
	}
	return nil
}

func (m *MockStore) Delete(id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(id)
	}
	return nil
}
