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

func TestLoginOTPRepository(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	userRepo := persistence.NewUserRepository(conn)
	otpRepo := persistence.NewLoginOTPRepository(conn)

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

	newRecord := func(userID uuid.UUID) *persistence.LoginOTPRecord {
		return &persistence.LoginOTPRecord{
			ID:        uuid.New(),
			UserID:    userID,
			Code:      "123456",
			ExpiresAt: time.Now().Add(10 * time.Minute),
		}
	}

	t.Run("Create and FindByUserID", func(t *testing.T) {
		t.Parallel()

		userID := createUser(t)
		record := newRecord(userID)
		t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM login_otps WHERE id = $1", record.ID) })

		if err := otpRepo.Create(context.Background(), record); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}
		if record.CreatedAt.IsZero() {
			t.Error("Create() did not populate CreatedAt")
		}

		got, err := otpRepo.FindByUserID(context.Background(), userID)
		if err != nil {
			t.Fatalf("FindByUserID() unexpected error: %v", err)
		}
		if got.Code != record.Code {
			t.Errorf("Code = %q, want %q", got.Code, record.Code)
		}
	})

	t.Run("FindByUserID returns ErrLoginOTPNotFound when no otp exists", func(t *testing.T) {
		t.Parallel()

		_, err := otpRepo.FindByUserID(context.Background(), uuid.New())
		if !errors.Is(err, persistence.ErrLoginOTPNotFound) {
			t.Errorf("FindByUserID() error = %v, want %v", err, persistence.ErrLoginOTPNotFound)
		}
	})

	t.Run("FindByUserID returns the most recent otp", func(t *testing.T) {
		t.Parallel()

		userID := createUser(t)

		older := newRecord(userID)
		t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM login_otps WHERE id = $1", older.ID) })
		if err := otpRepo.Create(context.Background(), older); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}

		newer := newRecord(userID)
		t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM login_otps WHERE id = $1", newer.ID) })
		if err := otpRepo.Create(context.Background(), newer); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}

		got, err := otpRepo.FindByUserID(context.Background(), userID)
		if err != nil {
			t.Fatalf("FindByUserID() unexpected error: %v", err)
		}
		if got.ID != newer.ID {
			t.Errorf("FindByUserID() returned otp %v, want the most recent %v", got.ID, newer.ID)
		}
	})

	t.Run("Delete removes the otp", func(t *testing.T) {
		t.Parallel()

		userID := createUser(t)
		record := newRecord(userID)
		if err := otpRepo.Create(context.Background(), record); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}

		if err := otpRepo.Delete(context.Background(), record.ID); err != nil {
			t.Fatalf("Delete() unexpected error: %v", err)
		}

		_, err := otpRepo.FindByUserID(context.Background(), userID)
		if !errors.Is(err, persistence.ErrLoginOTPNotFound) {
			t.Errorf("FindByUserID() after delete error = %v, want %v", err, persistence.ErrLoginOTPNotFound)
		}
	})

	t.Run("Delete returns ErrLoginOTPNotFound for an unknown id", func(t *testing.T) {
		t.Parallel()

		err := otpRepo.Delete(context.Background(), uuid.New())
		if !errors.Is(err, persistence.ErrLoginOTPNotFound) {
			t.Errorf("Delete() error = %v, want %v", err, persistence.ErrLoginOTPNotFound)
		}
	})
}
