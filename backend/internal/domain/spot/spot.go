package spot

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidQuery = errors.New("search query must not be empty")
	ErrNotFound     = errors.New("spot not found")
)

// FirstPlanID identifies the Plan that first added this Spot -- Plans have
// no owning User in this codebase (share-token based), so a Plan reference
// is the closest available notion of "who added it" for attribution.
//
// FirstPlanID must never be exposed in API responses unless
// IsAttributionVisible() is true: Plan.ID is not access-controlled the way
// ShareToken is (e.g. GET /plans/{plan_id}/pins takes the raw plan_id with
// no ShareToken check), so leaking it for a Plan that hasn't been published
// would let a stranger discover and access another user's private Plan.
//
// FirstPlanIsPublic is a denormalized copy of that Plan's IsPublic at the
// time this Spot was saved (see Application layer), so this check doesn't
// require Spot to depend on the plan package or re-query Plan on every
// read.
type Spot struct {
	ID                uuid.UUID
	PlaceID           PlaceID
	Name              string
	Address           string
	Location          Location
	FirstPlanID       uuid.UUID
	FirstPlanIsPublic bool
}

// IsAttributionVisible reports whether FirstPlanID is safe to expose. See
// the FirstPlanID field comment for why this check exists.
func (spot *Spot) IsAttributionVisible() bool {
	return spot.FirstPlanIsPublic
}
