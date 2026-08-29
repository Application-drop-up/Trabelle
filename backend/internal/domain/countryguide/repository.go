package countryguide

import "context"

type Repository interface {
	FindAll(ctx context.Context) ([]*CountryGuide, error)
	FindByCode(ctx context.Context, code string) (*CountryGuide, error)
}
