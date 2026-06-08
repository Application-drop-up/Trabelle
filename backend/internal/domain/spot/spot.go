package spot

import "errors"

var ErrInvalidQuery = errors.New("search query must not be empty")

type Spot struct {
	PlaceID  PlaceID
	Name     string
	Address  string
	Location Location
}
