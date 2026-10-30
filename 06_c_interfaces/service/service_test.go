package service_test

import (
	"errors"
	"testing"

	"06_c_interfaces/service"
)

func TestUserService_Register_Success(t *testing.T) {
	mockStore := service.NewMockStore()
	mockStore.FindByEmailFn = func(email string) (*service.User, error) {
		return nil, service.ErrUserNotFound
	}

	svc := service.NewUserService(mockStore)
	user, err := svc.Register("user-1", "alice", "alice@example.com")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if user.Username != "alice" {
		t.Errorf("expected username 'alice', got: '%s'", user.Username)
	}

	if len(mockStore.SavedUsers) != 1 {
		t.Errorf("expected 1 saved user in mock, got: %d", len(mockStore.SavedUsers))
	}
}

func TestUserService_Register_DuplicateEmail(t *testing.T) {
	mockStore := service.NewMockStore()
	mockStore.FindByEmailFn = func(email string) (*service.User, error) {
		// Simulate existing user with same email
		return &service.User{ID: "existing-1", Email: email}, nil
	}

	svc := service.NewUserService(mockStore)
	_, err := svc.Register("user-2", "bob", "bob@example.com")

	if !errors.Is(err, service.ErrUserAlreadyExists) {
		t.Fatalf("expected ErrUserAlreadyExists, got: %v", err)
	}

	if len(mockStore.SavedUsers) != 0 {
		t.Errorf("expected 0 saved users, got: %d", len(mockStore.SavedUsers))
	}
}

func TestUserService_GetUser_NotFound(t *testing.T) {
	mockStore := service.NewMockStore()
	mockStore.FindByIDFn = func(id string) (*service.User, error) {
		return nil, service.ErrUserNotFound
	}

	svc := service.NewUserService(mockStore)
	_, err := svc.GetUser("unknown-id")

	if !errors.Is(err, service.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}
