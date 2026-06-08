package spot

import "errors"

var ErrEmptyPlaceID = errors.New("place ID must not be empty")

type PlaceID string

func NewPlaceID(id string) (PlaceID, error) {
	if id == "" {
		return "", ErrEmptyPlaceID
	}
	return PlaceID(id), nil
}

func (p PlaceID) String() string {
	return string(p)
}
