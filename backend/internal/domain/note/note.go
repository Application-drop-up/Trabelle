package note

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("note not found")

type Note struct {
	ID        uuid.UUID
	PinID     uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
