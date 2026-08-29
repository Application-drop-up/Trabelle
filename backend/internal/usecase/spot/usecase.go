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

func (useCase *UseCase) SearchSpots(ctx context.Context, query string) ([]*spot.Spot, error) {
	if query == "" {
		return nil, spot.ErrInvalidQuery
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
