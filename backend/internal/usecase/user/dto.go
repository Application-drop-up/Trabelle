package user

import (
	"time"

	domain "github.com/Application-drop-up/Travellle/internal/domain/user"
	"github.com/google/uuid"
)

type UserDto struct {
	ID        uuid.UUID
	Email     string
	Name      string
	CreatedAt time.Time
}

func NewUserDto(user *domain.User) *UserDto {
	return &UserDto{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}
}
