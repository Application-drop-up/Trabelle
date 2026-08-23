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
type Spot struct {
	ID          uuid.UUID
	PlaceID     PlaceID
	Name        string
	Address     string
	Location    Location
	FirstPlanID uuid.UUID
}
