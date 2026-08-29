package spot

import "context"

type Repository interface {
	Save(ctx context.Context, spot *Spot) error
	FindByPlaceID(ctx context.Context, placeID PlaceID) (*Spot, error)
	Search(ctx context.Context, query string) ([]*Spot, error)
}
