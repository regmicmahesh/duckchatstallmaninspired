package server_test

import (
	"testing"

	"github.com/regmicmahesh/term-chat/internal/server"
)

func TestRegisterUser(t *testing.T) {

	t.Run("should register a user", func(t *testing.T) {
		s := server.NewServer()

		if len(s.RegisteredUsers) != 0 {
			t.Errorf("Expected 0 registered users, got %d", len(s.RegisteredUsers))
		}

		s.RegisterUser("user", "user")

		if len(s.RegisteredUsers) != 1 {
			t.Errorf("Expected 1 user, got %d", len(s.RegisteredUsers))
		}

	})

}

func TestIsUserCredentialsValid(t *testing.T) {

	t.Run("should return true if user credentials are valid", func(t *testing.T) {
		s := server.NewServer()

		s.RegisterUser("user", "user")

		if !s.IsUserCredentialsValid("user", "user") {
			t.Errorf("Expected true")
		}
	})

	t.Run("should return false if user credentials are invalid", func(t *testing.T) {
		s := server.NewServer()

		s.RegisterUser("user", "user")

		if s.IsUserCredentialsValid("user", "user2") {
			t.Errorf("Expected false")
		}
	})

}

func TestIsUserRegistered(t *testing.T) {

	t.Run("should return true if user is registered", func(t *testing.T) {
		s := server.NewServer()

		s.RegisterUser("user", "user")
		if !s.IsUserRegistered("user") {
			t.Errorf("Expected true")
		}
	})

	t.Run("should return false if user is not registered", func(t *testing.T) {
		s := server.NewServer()

		if s.IsUserRegistered("user") {
			t.Errorf("Expected false")
		}
	})

}
