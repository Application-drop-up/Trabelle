package plan

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("plan not found")

type Plan struct {
	ID         uuid.UUID
	Title      string
	ShareToken string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
