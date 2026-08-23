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

// IsViewableWithoutToken reports whether this Plan can be read without its
// ShareToken. ShareToken-based access (for an in-progress trip shared with
// specific people) and public access (a finished trip published for anyone,
// including strangers, to view) are separate concerns belonging to
// different phases of a Plan's lifecycle -- this method exists so callers
// check Plan's own rule about the latter, rather than reading IsPublic
// directly.
func (plan *Plan) IsViewableWithoutToken() bool {
	return plan.IsPublic
}
