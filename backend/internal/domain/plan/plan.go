package plan

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("plan not found")

// IsPublic defaults to false: a Plan is only visible to holders of its
// ShareToken until the owner explicitly shares it. This is the invariant
// that other domains (e.g. Spot, when exposing which Plan first added a
// place) rely on before it's safe to reveal a Plan's ID at all.
type Plan struct {
	ID         uuid.UUID
	Title      string
	ShareToken string
	IsPublic   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
