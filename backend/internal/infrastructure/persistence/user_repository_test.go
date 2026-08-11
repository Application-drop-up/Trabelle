package persistence_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	domain "github.com/Application-drop-up/Travellle/internal/domain/user"
	"github.com/Application-drop-up/Travellle/internal/infrastructure/persistence"
	"github.com/Application-drop-up/Travellle/internal/testutil"
)

func newTestUser() *domain.User {
	return &domain.User{
		ID:           uuid.New(),
		Email:        uuid.New().String() + "@example.com",
		PasswordHash: "hashed-password",
		Name:         "Taro",
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	repo := persistence.NewUserRepository(conn)

	t.Run("returns the user when it exists", func(t *testing.T) {
		t.Parallel()

		user := newTestUser()
		t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM users WHERE id = $1", user.ID) })
		if err := repo.Create(context.Background(), user); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}

		got, err := repo.FindByID(context.Background(), user.ID)
		if err != nil {
			t.Fatalf("FindByID() unexpected error: %v", err)
		}
		if got.Email != user.Email || got.Name != user.Name || got.PasswordHash != user.PasswordHash {
			t.Errorf("FindByID() = %+v, want fields matching %+v", got, user)
		}
	})

	t.Run("returns ErrNotFound for an unknown id", func(t *testing.T) {
		t.Parallel()

		_, err := repo.FindByID(context.Background(), uuid.New())
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("FindByID() error = %v, want %v", err, domain.ErrNotFound)
		}
	})
}

func TestUserRepository_Update(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	repo := persistence.NewUserRepository(conn)

	t.Run("updates the name", func(t *testing.T) {
		t.Parallel()

		user := newTestUser()
		t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM users WHERE id = $1", user.ID) })
		if err := repo.Create(context.Background(), user); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}

		user.Name = "Jiro"
		if err := repo.Update(context.Background(), user); err != nil {
			t.Fatalf("Update() unexpected error: %v", err)
		}

		got, err := repo.FindByID(context.Background(), user.ID)
		if err != nil {
			t.Fatalf("FindByID() unexpected error: %v", err)
		}
		if got.Name != "Jiro" {
			t.Errorf("Name = %q, want %q", got.Name, "Jiro")
		}
	})

	t.Run("returns ErrNotFound for an unknown id", func(t *testing.T) {
		t.Parallel()

		user := newTestUser()
		user.ID = uuid.New()

		err := repo.Update(context.Background(), user)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("Update() error = %v, want %v", err, domain.ErrNotFound)
		}
	})
}
