package user_test

import (
	"context"
	"errors"
	"testing"

	domain "github.com/Application-drop-up/Travellle/internal/domain/user"
	userUseCase "github.com/Application-drop-up/Travellle/internal/usecase/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// mockRepository は domain.Repository のテスト用実装
type mockRepository struct {
	createErr    error
	createdUsers []*domain.User
}

func (m *mockRepository) Create(_ context.Context, user *domain.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.createdUsers = append(m.createdUsers, user)
	return nil
}

func (m *mockRepository) FindByID(_ context.Context, _ uuid.UUID) (*domain.User, error) {
	return nil, domain.ErrNotFound
}

func TestUseCase_Register(t *testing.T) {
	t.Parallel()

	t.Run("returns dto and stores hashed password on success", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepository{}
		useCase := userUseCase.New(repo)

		command := userUseCase.RegisterCommand{
			Email:    "taro@example.com",
			Password: "password123",
			Name:     "Taro",
		}

		got, err := useCase.Register(context.Background(), command)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got.Email != command.Email {
			t.Errorf("got Email %q, want %q", got.Email, command.Email)
		}
		if got.Name != command.Name {
			t.Errorf("got Name %q, want %q", got.Name, command.Name)
		}
		if got.ID == uuid.Nil {
			t.Error("got zero ID, want a generated UUID")
		}

		if len(repo.createdUsers) != 1 {
			t.Fatalf("got %d created users, want 1", len(repo.createdUsers))
		}
		stored := repo.createdUsers[0]
		if stored.PasswordHash == command.Password {
			t.Error("PasswordHash was stored as plaintext")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(stored.PasswordHash), []byte(command.Password)); err != nil {
			t.Errorf("stored PasswordHash does not match plaintext password: %v", err)
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepository{createErr: errors.New("db error")}
		useCase := userUseCase.New(repo)

		command := userUseCase.RegisterCommand{
			Email:    "taro@example.com",
			Password: "password123",
			Name:     "Taro",
		}

		_, err := useCase.Register(context.Background(), command)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
