package spot

import "errors"

var ErrInvalidLocation = errors.New("latitude must be between -90 and 90, longitude between -180 and 180")

type Location struct {
	Latitude  float64
	Longitude float64
}

func NewLocation(lat, lng float64) (Location, error) {
	if lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		return Location{}, ErrInvalidLocation
	}
	return Location{Latitude: lat, Longitude: lng}, nil
}
