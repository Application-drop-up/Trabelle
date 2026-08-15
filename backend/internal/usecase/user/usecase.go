package user

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/Application-drop-up/Travellle/internal/domain/user"
	"github.com/google/uuid"
)

type UseCase struct {
	repo        domain.Repository
	sessionRepo SessionRepository
	otpRepo     LoginOTPRepository
	emailSender EmailSender
}

func New(
	repo domain.Repository,
	sessionRepo SessionRepository,
	otpRepo LoginOTPRepository,
	emailSender EmailSender,
) *UseCase {
	return &UseCase{
		repo:        repo,
		sessionRepo: sessionRepo,
		otpRepo:     otpRepo,
		emailSender: emailSender,
	}
}

type RegisterCommand struct {
	Email    string
	Password string
	Name     string
}

func (useCase *UseCase) Register(ctx context.Context, command RegisterCommand) (*UserDto, error) {
	email, err := domain.NewEmail(command.Email)
	if err != nil {
		return nil, err
	}

	password, err := domain.NewPassword(command.Password)
	if err != nil {
		return nil, err
	}

	_, err = useCase.repo.FindByEmail(ctx, email.String())
	if err == nil {
		return nil, domain.ErrEmailTaken
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        email.String(),
		PasswordHash: password.Hash(),
		Name:         command.Name,
	}

	if err := useCase.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return NewUserDto(user), nil
}

func (useCase *UseCase) GetUserByID(ctx context.Context, id uuid.UUID) (*UserDto, error) {
	user, err := useCase.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return NewUserDto(user), nil
}

type UpdateCommand struct {
	Name  string
	Email string
}

func (useCase *UseCase) UpdateUser(ctx context.Context, id uuid.UUID, command UpdateCommand) (*UserDto, error) {
	user, err := useCase.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	email, err := domain.NewEmail(command.Email)
	if err != nil {
		return nil, err
	}

	if email.String() != user.Email {
		existing, err := useCase.repo.FindByEmail(ctx, email.String())
		if err == nil && existing.ID != id {
			return nil, domain.ErrEmailTaken
		}
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("find user by email: %w", err)
		}
	}

	user.Name = command.Name
	user.Email = email.String()

	if err := useCase.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return NewUserDto(user), nil
}

func (useCase *UseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return useCase.repo.Delete(ctx, id)
}

// LoginStart verifies email+password and, on success, generates and emails
// a one-time code that must be confirmed via LoginVerify.
func (useCase *UseCase) LoginStart(ctx context.Context, email, password string) error {
	foundUser, err := useCase.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrInvalidCredentials
		}
		return fmt.Errorf("find user by email: %w", err)
	}

	storedPassword := domain.NewPasswordFromHash(foundUser.PasswordHash)
	if !storedPassword.Matches(password) {
		return domain.ErrInvalidCredentials
	}

	otp, err := newOTP(foundUser.ID)
	if err != nil {
		return fmt.Errorf("generate otp: %w", err)
	}
	if err := useCase.otpRepo.Create(ctx, otp); err != nil {
		return fmt.Errorf("create otp: %w", err)
	}

	if err := useCase.emailSender.SendLoginCode(ctx, foundUser.Email, otp.Code); err != nil {
		return fmt.Errorf("send login code: %w", err)
	}

	return nil
}

// LoginVerify checks the OTP code and, on success, issues a new Session.
func (useCase *UseCase) LoginVerify(ctx context.Context, email, code string) (*Session, error) {
	foundUser, err := useCase.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	otp, err := useCase.otpRepo.FindByUserID(ctx, foundUser.ID)
	if err != nil {
		if errors.Is(err, ErrOTPNotFound) {
			return nil, ErrInvalidOTP
		}
		return nil, fmt.Errorf("find otp by user id: %w", err)
	}
	if otp.isExpired() {
		return nil, ErrOTPExpired
	}
	if !otp.matches(code) {
		return nil, ErrInvalidOTP
	}

	session, err := newSession(foundUser.ID)
	if err != nil {
		return nil, fmt.Errorf("generate session: %w", err)
	}
	if err := useCase.sessionRepo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	if err := useCase.otpRepo.Delete(ctx, otp.ID); err != nil {
		return nil, fmt.Errorf("delete otp: %w", err)
	}

	return session, nil
}
