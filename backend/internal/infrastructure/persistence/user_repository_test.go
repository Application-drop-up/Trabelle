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

	t.Run("updates the name and email", func(t *testing.T) {
		t.Parallel()

		user := newTestUser()
		t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM users WHERE id = $1", user.ID) })
		if err := repo.Create(context.Background(), user); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}

		user.Name = "Jiro"
		user.Email = uuid.New().String() + "@example.com"
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
		if got.Email != user.Email {
			t.Errorf("Email = %q, want %q", got.Email, user.Email)
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

	t.Run("returns ErrEmailTaken when the email is already used by another user", func(t *testing.T) {
		t.Parallel()

		userA := newTestUser()
		t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM users WHERE id = $1", userA.ID) })
		if err := repo.Create(context.Background(), userA); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}

		userB := newTestUser()
		t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM users WHERE id = $1", userB.ID) })
		if err := repo.Create(context.Background(), userB); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}

		userB.Email = userA.Email
		err := repo.Update(context.Background(), userB)
		if !errors.Is(err, domain.ErrEmailTaken) {
			t.Errorf("Update() error = %v, want %v", err, domain.ErrEmailTaken)
		}
	})
}

func TestUserRepository_Delete(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	repo := persistence.NewUserRepository(conn)

	t.Run("deletes an existing user", func(t *testing.T) {
		t.Parallel()

		user := newTestUser()
		if err := repo.Create(context.Background(), user); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}

		if err := repo.Delete(context.Background(), user.ID); err != nil {
			t.Fatalf("Delete() unexpected error: %v", err)
		}

		_, err := repo.FindByID(context.Background(), user.ID)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("FindByID() after delete error = %v, want %v", err, domain.ErrNotFound)
		}
	})

	t.Run("returns ErrNotFound for an unknown id", func(t *testing.T) {
		t.Parallel()

		err := repo.Delete(context.Background(), uuid.New())
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("Delete() error = %v, want %v", err, domain.ErrNotFound)
		}
	})
}
