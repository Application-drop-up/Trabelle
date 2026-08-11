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
	createErr      error
	createdUsers   []*domain.User
	existingByMail *domain.User
	findByEmailErr error
	existingByID   *domain.User
	findByIDErr    error
	updateErr      error
	updatedUsers   []*domain.User
}

func (m *mockRepository) Create(_ context.Context, user *domain.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.createdUsers = append(m.createdUsers, user)
	return nil
}

func (m *mockRepository) FindByID(_ context.Context, _ uuid.UUID) (*domain.User, error) {
	if m.existingByID != nil {
		return m.existingByID, nil
	}
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	return nil, domain.ErrNotFound
}

func (m *mockRepository) FindByEmail(_ context.Context, _ string) (*domain.User, error) {
	if m.existingByMail != nil {
		return m.existingByMail, nil
	}
	if m.findByEmailErr != nil {
		return nil, m.findByEmailErr
	}
	return nil, domain.ErrNotFound
}

func (m *mockRepository) Update(_ context.Context, user *domain.User) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.updatedUsers = append(m.updatedUsers, user)
	return nil
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

	t.Run("returns ErrInvalidEmail for a malformed email", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepository{}
		useCase := userUseCase.New(repo)

		command := userUseCase.RegisterCommand{
			Email:    "not-an-email",
			Password: "password123",
			Name:     "Taro",
		}

		_, err := useCase.Register(context.Background(), command)
		if !errors.Is(err, domain.ErrInvalidEmail) {
			t.Fatalf("got error %v, want %v", err, domain.ErrInvalidEmail)
		}
	})

	t.Run("returns ErrPasswordTooShort for a short password", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepository{}
		useCase := userUseCase.New(repo)

		command := userUseCase.RegisterCommand{
			Email:    "taro@example.com",
			Password: "short",
			Name:     "Taro",
		}

		_, err := useCase.Register(context.Background(), command)
		if !errors.Is(err, domain.ErrPasswordTooShort) {
			t.Fatalf("got error %v, want %v", err, domain.ErrPasswordTooShort)
		}
	})

	t.Run("returns ErrEmailTaken when the email is already registered", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepository{existingByMail: &domain.User{Email: "taro@example.com"}}
		useCase := userUseCase.New(repo)

		command := userUseCase.RegisterCommand{
			Email:    "taro@example.com",
			Password: "password123",
			Name:     "Taro",
		}

		_, err := useCase.Register(context.Background(), command)
		if !errors.Is(err, domain.ErrEmailTaken) {
			t.Fatalf("got error %v, want %v", err, domain.ErrEmailTaken)
		}
		if len(repo.createdUsers) != 0 {
			t.Error("Create should not be called when email is taken")
		}
	})
}

func TestUseCase_GetUserByID(t *testing.T) {
	t.Parallel()

	t.Run("returns dto when user exists", func(t *testing.T) {
		t.Parallel()

		id := uuid.New()
		repo := &mockRepository{existingByID: &domain.User{
			ID:    id,
			Email: "taro@example.com",
			Name:  "Taro",
		}}
		useCase := userUseCase.New(repo)

		got, err := useCase.GetUserByID(context.Background(), id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != id {
			t.Errorf("got ID %v, want %v", got.ID, id)
		}
		if got.Email != "taro@example.com" {
			t.Errorf("got Email %q, want %q", got.Email, "taro@example.com")
		}
	})

	t.Run("returns ErrNotFound when user does not exist", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepository{}
		useCase := userUseCase.New(repo)

		_, err := useCase.GetUserByID(context.Background(), uuid.New())
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("got error %v, want %v", err, domain.ErrNotFound)
		}
	})
}

func TestUseCase_UpdateUser(t *testing.T) {
	t.Parallel()

	t.Run("returns dto with updated name on success", func(t *testing.T) {
		t.Parallel()

		id := uuid.New()
		repo := &mockRepository{existingByID: &domain.User{
			ID:    id,
			Email: "taro@example.com",
			Name:  "Taro",
		}}
		useCase := userUseCase.New(repo)

		got, err := useCase.UpdateUser(context.Background(), id, userUseCase.UpdateCommand{Name: "Jiro"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Name != "Jiro" {
			t.Errorf("got Name %q, want %q", got.Name, "Jiro")
		}
		if len(repo.updatedUsers) != 1 {
			t.Fatalf("got %d updated users, want 1", len(repo.updatedUsers))
		}
		if repo.updatedUsers[0].Name != "Jiro" {
			t.Errorf("updated user Name = %q, want %q", repo.updatedUsers[0].Name, "Jiro")
		}
	})

	t.Run("returns ErrNotFound when user does not exist", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepository{}
		useCase := userUseCase.New(repo)

		_, err := useCase.UpdateUser(context.Background(), uuid.New(), userUseCase.UpdateCommand{Name: "Jiro"})
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("got error %v, want %v", err, domain.ErrNotFound)
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		t.Parallel()

		id := uuid.New()
		repo := &mockRepository{
			existingByID: &domain.User{ID: id, Name: "Taro"},
			updateErr:    errors.New("db error"),
		}
		useCase := userUseCase.New(repo)

		_, err := useCase.UpdateUser(context.Background(), id, userUseCase.UpdateCommand{Name: "Jiro"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
