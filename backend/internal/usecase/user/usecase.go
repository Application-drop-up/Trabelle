package user

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/Application-drop-up/Travellle/internal/domain/user"
	"github.com/google/uuid"
)

type UseCase struct {
	repo domain.Repository
}

func New(repo domain.Repository) *UseCase {
	return &UseCase{repo: repo}
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
