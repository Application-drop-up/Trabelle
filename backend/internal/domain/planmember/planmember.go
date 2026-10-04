package planmember

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound      = errors.New("plan member not found")
	ErrAlreadyMember = errors.New("user is already a member of this plan")
)

type PlanMember struct {
	ID        uuid.UUID
	PlanID    uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
}
