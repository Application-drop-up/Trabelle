package spot

import (
	"context"
	"fmt"

	"github.com/Application-drop-up/Travellle/internal/domain/spot"
)

type UseCase struct {
	searcher spot.Searcher
}

func New(searcher spot.Searcher) *UseCase {
	return &UseCase{searcher: searcher}
}

func (uc *UseCase) SearchSpots(ctx context.Context, query string) ([]*spot.Spot, error) {
	if query == "" {
		return nil, spot.ErrInvalidQuery
	}

	spots, err := uc.searcher.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search spots: %w", err)
	}

	return spots, nil
}
