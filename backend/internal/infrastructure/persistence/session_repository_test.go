package persistence_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	userdomain "github.com/Application-drop-up/Travellle/internal/domain/user"
	"github.com/Application-drop-up/Travellle/internal/infrastructure/persistence"
	"github.com/Application-drop-up/Travellle/internal/testutil"
)

func TestSessionRepository(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	userRepo := persistence.NewUserRepository(conn)
	sessionRepo := persistence.NewSessionRepository(conn)

	createUser := func(t *testing.T) uuid.UUID {
		t.Helper()
		user := &userdomain.User{
			ID:           uuid.New(),
			Email:        uuid.New().String() + "@example.com",
			PasswordHash: "hashed-password",
			Name:         "Taro",
		}
		if err := userRepo.Create(context.Background(), user); err != nil {
			t.Fatalf("Create() user unexpected error: %v", err)
		}
		t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM users WHERE id = $1", user.ID) })
		return user.ID
	}

	newRecord := func(userID uuid.UUID) *persistence.SessionRecord {
		return &persistence.SessionRecord{
			ID:        uuid.New(),
			UserID:    userID,
			Token:     uuid.New().String(),
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
	}

	t.Run("Create and FindByToken", func(t *testing.T) {
		t.Parallel()

		userID := createUser(t)
		record := newRecord(userID)
		t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM sessions WHERE id = $1", record.ID) })

		if err := sessionRepo.Create(context.Background(), record); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}
		if record.CreatedAt.IsZero() {
			t.Error("Create() did not populate CreatedAt")
		}

		got, err := sessionRepo.FindByToken(context.Background(), record.Token)
		if err != nil {
			t.Fatalf("FindByToken() unexpected error: %v", err)
		}
		if got.UserID != userID {
			t.Errorf("UserID = %v, want %v", got.UserID, userID)
		}
	})

	t.Run("FindByToken returns ErrSessionNotFound for an unknown token", func(t *testing.T) {
		t.Parallel()

		_, err := sessionRepo.FindByToken(context.Background(), "unknown-token")
		if !errors.Is(err, persistence.ErrSessionNotFound) {
			t.Errorf("FindByToken() error = %v, want %v", err, persistence.ErrSessionNotFound)
		}
	})

	t.Run("Delete removes the session", func(t *testing.T) {
		t.Parallel()

		userID := createUser(t)
		record := newRecord(userID)
		if err := sessionRepo.Create(context.Background(), record); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}

		if err := sessionRepo.Delete(context.Background(), record.Token); err != nil {
			t.Fatalf("Delete() unexpected error: %v", err)
		}

		_, err := sessionRepo.FindByToken(context.Background(), record.Token)
		if !errors.Is(err, persistence.ErrSessionNotFound) {
			t.Errorf("FindByToken() after delete error = %v, want %v", err, persistence.ErrSessionNotFound)
		}
	})

	t.Run("Delete returns ErrSessionNotFound for an unknown token", func(t *testing.T) {
		t.Parallel()

		err := sessionRepo.Delete(context.Background(), "unknown-token")
		if !errors.Is(err, persistence.ErrSessionNotFound) {
			t.Errorf("Delete() error = %v, want %v", err, persistence.ErrSessionNotFound)
		}
	})
}
