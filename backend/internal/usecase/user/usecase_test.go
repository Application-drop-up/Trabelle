package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

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

// mockSessionRepository は userUseCase.SessionRepository のテスト用実装
type mockSessionRepository struct {
	createErr       error
	createdSession  *userUseCase.Session
	existingSession *userUseCase.Session
	findByTokenErr  error
	deleteErr       error
	deletedTokens   []string
}

func (m *mockSessionRepository) Create(_ context.Context, session *userUseCase.Session) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.createdSession = session
	return nil
}

func (m *mockSessionRepository) FindByToken(_ context.Context, _ string) (*userUseCase.Session, error) {
	if m.existingSession != nil {
		return m.existingSession, nil
	}
	if m.findByTokenErr != nil {
		return nil, m.findByTokenErr
	}
	return nil, userUseCase.ErrSessionNotFound
}

func (m *mockSessionRepository) Delete(_ context.Context, token string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	m.deletedTokens = append(m.deletedTokens, token)
	return nil
}

// mockLoginOTPRepository は userUseCase.LoginOTPRepository のテスト用実装
type mockLoginOTPRepository struct {
	createErr   error
	createdOTPs []*userUseCase.OTP
	existingOTP *userUseCase.OTP
	findErr     error
	deleteErr   error
	deletedIDs  []uuid.UUID
}

func (m *mockLoginOTPRepository) Create(_ context.Context, otp *userUseCase.OTP) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.createdOTPs = append(m.createdOTPs, otp)
	return nil
}

func (m *mockLoginOTPRepository) FindByUserID(_ context.Context, _ uuid.UUID) (*userUseCase.OTP, error) {
	if m.existingOTP != nil {
		return m.existingOTP, nil
	}
	if m.findErr != nil {
		return nil, m.findErr
	}
	return nil, userUseCase.ErrOTPNotFound
}

func (m *mockLoginOTPRepository) Delete(_ context.Context, id uuid.UUID) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	m.deletedIDs = append(m.deletedIDs, id)
	return nil
}

// mockEmailSender は userUseCase.EmailSender のテスト用実装
type mockEmailSender struct {
	sendErr  error
	sentTo   string
	sentCode string
}

func (m *mockEmailSender) SendLoginCode(_ context.Context, to, code string) error {
	if m.sendErr != nil {
		return m.sendErr
	}
	m.sentTo = to
	m.sentCode = code
	return nil
}

// newUseCase builds a UseCase with default (no-op) Session/OTP/Email mocks,
// for tests that only exercise Register/GetUserByID/UpdateUser/DeleteUser.
func newUseCase(repo *mockRepository) *userUseCase.UseCase {
	return userUseCase.New(repo, &mockSessionRepository{}, &mockLoginOTPRepository{}, &mockEmailSender{})
}

func TestUseCase_Register(t *testing.T) {
	t.Parallel()

	t.Run("returns dto and stores hashed password on success", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepository{}
		useCase := newUseCase(repo)

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
		useCase := newUseCase(repo)

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
		useCase := newUseCase(repo)

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
		useCase := newUseCase(repo)

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
		useCase := newUseCase(repo)

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
		useCase := newUseCase(repo)

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
		useCase := newUseCase(repo)

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
		useCase := newUseCase(repo)

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
		useCase := newUseCase(repo)

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
		useCase := newUseCase(repo)

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
		useCase := newUseCase(repo)

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
		useCase := newUseCase(repo)

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
		useCase := newUseCase(repo)

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
		useCase := newUseCase(repo)

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
		useCase := newUseCase(repo)

		err := useCase.DeleteUser(context.Background(), uuid.New())
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("got error %v, want %v", err, domain.ErrNotFound)
		}
	})
}

func hashPassword(t *testing.T, plaintext string) string {
	t.Helper()
	password, err := domain.NewPassword(plaintext)
	if err != nil {
		t.Fatalf("hashPassword: unexpected error: %v", err)
	}
	return password.Hash()
}

func TestUseCase_LoginStart(t *testing.T) {
	t.Parallel()

	t.Run("generates and sends an otp on success", func(t *testing.T) {
		t.Parallel()

		user := &domain.User{
			ID:           uuid.New(),
			Email:        "taro@example.com",
			PasswordHash: hashPassword(t, "password123"),
			Name:         "Taro",
		}
		repo := &mockRepository{existingByMail: user}
		otpRepo := &mockLoginOTPRepository{}
		emailSender := &mockEmailSender{}
		useCase := userUseCase.New(repo, &mockSessionRepository{}, otpRepo, emailSender)

		code, err := useCase.LoginStart(context.Background(), user.Email, "password123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(otpRepo.createdOTPs) != 1 {
			t.Fatalf("createdOTPs = %d, want 1", len(otpRepo.createdOTPs))
		}
		if otpRepo.createdOTPs[0].UserID != user.ID {
			t.Errorf("otp UserID = %v, want %v", otpRepo.createdOTPs[0].UserID, user.ID)
		}
		if code != otpRepo.createdOTPs[0].Code {
			t.Errorf("code = %q, want %q", code, otpRepo.createdOTPs[0].Code)
		}
		if emailSender.sentTo != user.Email {
			t.Errorf("sentTo = %q, want %q", emailSender.sentTo, user.Email)
		}
		if emailSender.sentCode != otpRepo.createdOTPs[0].Code {
			t.Errorf("sentCode = %q, want %q", emailSender.sentCode, otpRepo.createdOTPs[0].Code)
		}
	})

	t.Run("returns ErrInvalidCredentials when the email is unknown", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepository{}
		useCase := newUseCase(repo)

		_, err := useCase.LoginStart(context.Background(), "unknown@example.com", "password123")
		if !errors.Is(err, domain.ErrInvalidCredentials) {
			t.Fatalf("got error %v, want %v", err, domain.ErrInvalidCredentials)
		}
	})

	t.Run("returns ErrInvalidCredentials when the password does not match", func(t *testing.T) {
		t.Parallel()

		user := &domain.User{
			ID:           uuid.New(),
			Email:        "taro@example.com",
			PasswordHash: hashPassword(t, "password123"),
			Name:         "Taro",
		}
		repo := &mockRepository{existingByMail: user}
		useCase := newUseCase(repo)

		_, err := useCase.LoginStart(context.Background(), user.Email, "wrong-password")
		if !errors.Is(err, domain.ErrInvalidCredentials) {
			t.Fatalf("got error %v, want %v", err, domain.ErrInvalidCredentials)
		}
	})

	t.Run("propagates email sender error", func(t *testing.T) {
		t.Parallel()

		user := &domain.User{
			ID:           uuid.New(),
			Email:        "taro@example.com",
			PasswordHash: hashPassword(t, "password123"),
			Name:         "Taro",
		}
		repo := &mockRepository{existingByMail: user}
		emailSender := &mockEmailSender{sendErr: errors.New("smtp error")}
		useCase := userUseCase.New(repo, &mockSessionRepository{}, &mockLoginOTPRepository{}, emailSender)

		_, err := useCase.LoginStart(context.Background(), user.Email, "password123")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestUseCase_LoginVerify(t *testing.T) {
	t.Parallel()

	t.Run("issues a session on success", func(t *testing.T) {
		t.Parallel()

		user := &domain.User{ID: uuid.New(), Email: "taro@example.com", Name: "Taro"}
		otp := &userUseCase.OTP{ID: uuid.New(), UserID: user.ID, Code: "123456", ExpiresAt: time.Now().Add(10 * time.Minute)}
		repo := &mockRepository{existingByMail: user}
		otpRepo := &mockLoginOTPRepository{existingOTP: otp}
		sessionRepo := &mockSessionRepository{}
		useCase := userUseCase.New(repo, sessionRepo, otpRepo, &mockEmailSender{})

		session, err := useCase.LoginVerify(context.Background(), user.Email, "123456")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if session.UserID != user.ID {
			t.Errorf("session UserID = %v, want %v", session.UserID, user.ID)
		}
		if sessionRepo.createdSession != session {
			t.Error("session was not persisted via SessionRepository.Create")
		}
		if len(otpRepo.deletedIDs) != 1 || otpRepo.deletedIDs[0] != otp.ID {
			t.Errorf("deletedIDs = %v, want [%v]", otpRepo.deletedIDs, otp.ID)
		}
	})

	t.Run("returns ErrInvalidCredentials when the email is unknown", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepository{}
		useCase := newUseCase(repo)

		_, err := useCase.LoginVerify(context.Background(), "unknown@example.com", "123456")
		if !errors.Is(err, domain.ErrInvalidCredentials) {
			t.Fatalf("got error %v, want %v", err, domain.ErrInvalidCredentials)
		}
	})

	t.Run("returns ErrInvalidOTP when no otp exists", func(t *testing.T) {
		t.Parallel()

		user := &domain.User{ID: uuid.New(), Email: "taro@example.com", Name: "Taro"}
		repo := &mockRepository{existingByMail: user}
		useCase := newUseCase(repo)

		_, err := useCase.LoginVerify(context.Background(), user.Email, "123456")
		if !errors.Is(err, userUseCase.ErrInvalidOTP) {
			t.Fatalf("got error %v, want %v", err, userUseCase.ErrInvalidOTP)
		}
	})

	t.Run("returns ErrOTPExpired when the otp has expired", func(t *testing.T) {
		t.Parallel()

		user := &domain.User{ID: uuid.New(), Email: "taro@example.com", Name: "Taro"}
		otp := &userUseCase.OTP{ID: uuid.New(), UserID: user.ID, Code: "123456", ExpiresAt: time.Now().Add(-time.Minute)}
		repo := &mockRepository{existingByMail: user}
		otpRepo := &mockLoginOTPRepository{existingOTP: otp}
		useCase := userUseCase.New(repo, &mockSessionRepository{}, otpRepo, &mockEmailSender{})

		_, err := useCase.LoginVerify(context.Background(), user.Email, "123456")
		if !errors.Is(err, userUseCase.ErrOTPExpired) {
			t.Fatalf("got error %v, want %v", err, userUseCase.ErrOTPExpired)
		}
	})

	t.Run("returns ErrInvalidOTP when the code does not match", func(t *testing.T) {
		t.Parallel()

		user := &domain.User{ID: uuid.New(), Email: "taro@example.com", Name: "Taro"}
		otp := &userUseCase.OTP{ID: uuid.New(), UserID: user.ID, Code: "123456", ExpiresAt: time.Now().Add(10 * time.Minute)}
		repo := &mockRepository{existingByMail: user}
		otpRepo := &mockLoginOTPRepository{existingOTP: otp}
		useCase := userUseCase.New(repo, &mockSessionRepository{}, otpRepo, &mockEmailSender{})

		_, err := useCase.LoginVerify(context.Background(), user.Email, "999999")
		if !errors.Is(err, userUseCase.ErrInvalidOTP) {
			t.Fatalf("got error %v, want %v", err, userUseCase.ErrInvalidOTP)
		}
	})

	t.Run("propagates session repository error", func(t *testing.T) {
		t.Parallel()

		user := &domain.User{ID: uuid.New(), Email: "taro@example.com", Name: "Taro"}
		otp := &userUseCase.OTP{ID: uuid.New(), UserID: user.ID, Code: "123456", ExpiresAt: time.Now().Add(10 * time.Minute)}
		repo := &mockRepository{existingByMail: user}
		otpRepo := &mockLoginOTPRepository{existingOTP: otp}
		sessionRepo := &mockSessionRepository{createErr: errors.New("db error")}
		useCase := userUseCase.New(repo, sessionRepo, otpRepo, &mockEmailSender{})

		_, err := useCase.LoginVerify(context.Background(), user.Email, "123456")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestUseCase_CurrentUser(t *testing.T) {
	t.Parallel()

	t.Run("returns the user for a valid session", func(t *testing.T) {
		t.Parallel()

		user := &domain.User{ID: uuid.New(), Email: "taro@example.com", Name: "Taro"}
		session := &userUseCase.Session{
			ID:        uuid.New(),
			UserID:    user.ID,
			Token:     "valid-token",
			ExpiresAt: time.Now().Add(time.Hour),
		}
		repo := &mockRepository{existingByID: user}
		sessionRepo := &mockSessionRepository{existingSession: session}
		useCase := userUseCase.New(repo, sessionRepo, &mockLoginOTPRepository{}, &mockEmailSender{})

		dto, err := useCase.CurrentUser(context.Background(), session.Token)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dto.ID != user.ID {
			t.Errorf("ID = %v, want %v", dto.ID, user.ID)
		}
	})

	t.Run("returns ErrSessionNotFound when the session does not exist", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepository{}
		useCase := newUseCase(repo)

		_, err := useCase.CurrentUser(context.Background(), "unknown-token")
		if !errors.Is(err, userUseCase.ErrSessionNotFound) {
			t.Fatalf("got error %v, want %v", err, userUseCase.ErrSessionNotFound)
		}
	})

	t.Run("returns ErrSessionNotFound when the session has expired", func(t *testing.T) {
		t.Parallel()

		session := &userUseCase.Session{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			Token:     "expired-token",
			ExpiresAt: time.Now().Add(-time.Hour),
		}
		repo := &mockRepository{}
		sessionRepo := &mockSessionRepository{existingSession: session}
		useCase := userUseCase.New(repo, sessionRepo, &mockLoginOTPRepository{}, &mockEmailSender{})

		_, err := useCase.CurrentUser(context.Background(), session.Token)
		if !errors.Is(err, userUseCase.ErrSessionNotFound) {
			t.Fatalf("got error %v, want %v", err, userUseCase.ErrSessionNotFound)
		}
	})

	t.Run("propagates repository error when finding the user fails", func(t *testing.T) {
		t.Parallel()

		session := &userUseCase.Session{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			Token:     "valid-token",
			ExpiresAt: time.Now().Add(time.Hour),
		}
		repo := &mockRepository{findByIDErr: errors.New("db error")}
		sessionRepo := &mockSessionRepository{existingSession: session}
		useCase := userUseCase.New(repo, sessionRepo, &mockLoginOTPRepository{}, &mockEmailSender{})

		_, err := useCase.CurrentUser(context.Background(), session.Token)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestUseCase_Logout(t *testing.T) {
	t.Parallel()

	t.Run("deletes the session", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepository{}
		sessionRepo := &mockSessionRepository{}
		useCase := userUseCase.New(repo, sessionRepo, &mockLoginOTPRepository{}, &mockEmailSender{})

		if err := useCase.Logout(context.Background(), "some-token"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(sessionRepo.deletedTokens) != 1 || sessionRepo.deletedTokens[0] != "some-token" {
			t.Errorf("deletedTokens = %v, want [some-token]", sessionRepo.deletedTokens)
		}
	})

	t.Run("succeeds when the session does not exist", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepository{}
		sessionRepo := &mockSessionRepository{deleteErr: userUseCase.ErrSessionNotFound}
		useCase := userUseCase.New(repo, sessionRepo, &mockLoginOTPRepository{}, &mockEmailSender{})

		if err := useCase.Logout(context.Background(), "unknown-token"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		t.Parallel()

		repo := &mockRepository{}
		sessionRepo := &mockSessionRepository{deleteErr: errors.New("db error")}
		useCase := userUseCase.New(repo, sessionRepo, &mockLoginOTPRepository{}, &mockEmailSender{})

		if err := useCase.Logout(context.Background(), "some-token"); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
