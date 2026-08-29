package spot

import (
	"context"
	"fmt"

	"github.com/Application-drop-up/Travellle/internal/domain/spot"
)

type UseCase struct {
	searcher spot.Searcher
	repo     spot.Repository
}

func New(searcher spot.Searcher, repo spot.Repository) *UseCase {
	return &UseCase{searcher: searcher, repo: repo}
}

// SearchSpots searches previously-saved Spots first, so results already
// known from other users' plans don't cost a call to the external Searcher.
// It only falls back to the Searcher when nothing is known locally yet.
func (useCase *UseCase) SearchSpots(ctx context.Context, query string) ([]*spot.Spot, error) {
	if query == "" {
		return nil, spot.ErrInvalidQuery
	}

	localSpots, err := useCase.repo.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search local spots: %w", err)
	}
	if len(localSpots) > 0 {
		return localSpots, nil
	}

	spots, err := useCase.searcher.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search spots: %w", err)
	}

	return spots, nil
}

func (useCase *UseCase) SaveSpot(ctx context.Context, newSpot *spot.Spot) error {
	if err := useCase.repo.Save(ctx, newSpot); err != nil {
		return fmt.Errorf("save spot: %w", err)
	}

	return nil
}
