package user

import (
	"context"
	"fmt"

	domain "github.com/Application-drop-up/Travellle/internal/domain/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func (useCase *UseCase) Register(context context.Context, command RegisterCommand) (*UserDto, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(command.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        command.Email,
		PasswordHash: string(passwordHash),
		Name:         command.Name,
	}

	if err := useCase.repo.Create(context, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return NewUserDto(user), nil
}
