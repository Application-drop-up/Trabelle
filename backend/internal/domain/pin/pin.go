package pin

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("pin not found")

type Category string

const (
	CategoryRestaurant  Category = "restaurant"
	CategoryHotel       Category = "hotel"
	CategorySightseeing Category = "sightseeing"
	CategoryTransport   Category = "transport"
	CategoryOther       Category = "other"
)

func (c Category) IsValid() bool {
	switch c {
	case CategoryRestaurant, CategoryHotel, CategorySightseeing, CategoryTransport, CategoryOther:
		return true
	}
	return false
}

type Pin struct {
	ID        uuid.UUID
	PlanID    uuid.UUID
	Name      string
	Latitude  float64
	Longitude float64
	Category  Category
	Colour    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
