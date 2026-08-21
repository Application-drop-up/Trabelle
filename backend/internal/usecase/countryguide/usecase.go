package countryguide

import (
	"context"
	"fmt"

	domain "github.com/Application-drop-up/Travellle/internal/domain/countryguide"
)

type UseCase struct {
	repo domain.Repository
}

func New(repo domain.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (useCase *UseCase) ListCountryGuides(ctx context.Context) (CountryGuideDtoCollection, error) {
	guides, err := useCase.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list country guides: %w", err)
	}
	return NewCountryGuideDtoCollection(guides), nil
}

func (useCase *UseCase) GetCountryGuide(ctx context.Context, code string) (*CountryGuideDto, error) {
	guide, err := useCase.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return NewCountryGuideDto(guide), nil
}
