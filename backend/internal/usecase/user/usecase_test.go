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
	deleteErr      error
	deletedIDs     []uuid.UUID
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

func (m *mockRepository) Delete(_ context.Context, id uuid.UUID) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	m.deletedIDs = append(m.deletedIDs, id)
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

	t.Run("returns dto with updated name and email on success", func(t *testing.T) {
		t.Parallel()

		id := uuid.New()
		repo := &mockRepository{existingByID: &domain.User{
			ID:    id,
			Email: "taro@example.com",
			Name:  "Taro",
		}}
		useCase := userUseCase.New(repo)

		command := userUseCase.UpdateCommand{Name: "Jiro", Email: "jiro@example.com"}
		got, err := useCase.UpdateUser(context.Background(), id, command)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.Name != "Jiro" {
			t.Errorf("got Name %q, want %q", got.Name, "Jiro")
		}
		if got.Email != "jiro@example.com" {
			t.Errorf("got Email %q, want %q", got.Email, "jiro@example.com")
		}
		if len(repo.updatedUsers) != 1 {
			t.Fatalf("got %d updated users, want 1", len(repo.updatedUsers))
		}
	})

	t.Run("returns ErrNotFound when user does not exist", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepository{}
		useCase := userUseCase.New(repo)

		command := userUseCase.UpdateCommand{Name: "Jiro", Email: "jiro@example.com"}
		_, err := useCase.UpdateUser(context.Background(), uuid.New(), command)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("got error %v, want %v", err, domain.ErrNotFound)
		}
	})

	t.Run("returns ErrInvalidEmail for a malformed email", func(t *testing.T) {
		t.Parallel()

		id := uuid.New()
		repo := &mockRepository{existingByID: &domain.User{ID: id, Email: "taro@example.com", Name: "Taro"}}
		useCase := userUseCase.New(repo)

		command := userUseCase.UpdateCommand{Name: "Jiro", Email: "not-an-email"}
		_, err := useCase.UpdateUser(context.Background(), id, command)
		if !errors.Is(err, domain.ErrInvalidEmail) {
			t.Fatalf("got error %v, want %v", err, domain.ErrInvalidEmail)
		}
	})

	t.Run("returns ErrEmailTaken when the email is already used by another user", func(t *testing.T) {
		t.Parallel()

		id := uuid.New()
		repo := &mockRepository{
			existingByID:   &domain.User{ID: id, Email: "taro@example.com", Name: "Taro"},
			existingByMail: &domain.User{ID: uuid.New(), Email: "jiro@example.com"},
		}
		useCase := userUseCase.New(repo)

		command := userUseCase.UpdateCommand{Name: "Taro", Email: "jiro@example.com"}
		_, err := useCase.UpdateUser(context.Background(), id, command)
		if !errors.Is(err, domain.ErrEmailTaken) {
			t.Fatalf("got error %v, want %v", err, domain.ErrEmailTaken)
		}
		if len(repo.updatedUsers) != 0 {
			t.Error("Update should not be called when email is taken")
		}
	})

	t.Run("allows keeping the same email", func(t *testing.T) {
		t.Parallel()

		id := uuid.New()
		repo := &mockRepository{existingByID: &domain.User{ID: id, Email: "taro@example.com", Name: "Taro"}}
		useCase := userUseCase.New(repo)

		command := userUseCase.UpdateCommand{Name: "Jiro", Email: "taro@example.com"}
		_, err := useCase.UpdateUser(context.Background(), id, command)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		t.Parallel()

		id := uuid.New()
		repo := &mockRepository{
			existingByID: &domain.User{ID: id, Email: "taro@example.com", Name: "Taro"},
			updateErr:    errors.New("db error"),
		}
		useCase := userUseCase.New(repo)

		command := userUseCase.UpdateCommand{Name: "Jiro", Email: "taro@example.com"}
		_, err := useCase.UpdateUser(context.Background(), id, command)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestUseCase_DeleteUser(t *testing.T) {
	t.Parallel()

	t.Run("deletes the user on success", func(t *testing.T) {
		t.Parallel()

		id := uuid.New()
		repo := &mockRepository{}
		useCase := userUseCase.New(repo)

		if err := useCase.DeleteUser(context.Background(), id); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(repo.deletedIDs) != 1 || repo.deletedIDs[0] != id {
			t.Errorf("deletedIDs = %v, want [%v]", repo.deletedIDs, id)
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepository{deleteErr: domain.ErrNotFound}
		useCase := userUseCase.New(repo)

		err := useCase.DeleteUser(context.Background(), uuid.New())
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("got error %v, want %v", err, domain.ErrNotFound)
		}
	})
}
