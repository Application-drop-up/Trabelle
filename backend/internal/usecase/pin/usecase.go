package pin

import (
	"context"
	"fmt"

	domain "github.com/Application-drop-up/Travellle/internal/domain/pin"
	spotdomain "github.com/Application-drop-up/Travellle/internal/domain/spot"
	spotuc "github.com/Application-drop-up/Travellle/internal/usecase/spot"
	"github.com/google/uuid"
)

type UseCase struct {
	repo        domain.Repository
	spotUseCase *spotuc.UseCase
}

func New(repo domain.Repository, spotUseCase *spotuc.UseCase) *UseCase {
	return &UseCase{repo: repo, spotUseCase: spotUseCase}
}

type CreateInput struct {
	PlanID    uuid.UUID
	Name      string
	Latitude  float64
	Longitude float64
	Category  domain.Category
	Colour    string
	// PlaceID and Address are optional: they're only set when the Pin was
	// created from a Spot search result, and are used to cache that Spot
	// (see saveSpot) so future searches can reuse it without calling
	// TomTom again.
	PlaceID string
	Address string
}

type UpdateInput struct {
	Category *domain.Category
	Colour   *string
}

func (useCase *UseCase) CreatePin(ctx context.Context, input CreateInput) (*domain.Pin, error) {
	pin := &domain.Pin{
		ID:        uuid.New(),
		PlanID:    input.PlanID,
		Name:      input.Name,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
		Category:  input.Category,
		Colour:    input.Colour,
	}

	if err := useCase.repo.Create(ctx, pin); err != nil {
		return nil, fmt.Errorf("create pin: %w", err)
	}

	if input.PlaceID != "" {
		useCase.saveSpot(ctx, input)
	}

	return pin, nil
}

// saveSpot best-effort persists the Spot behind a newly created Pin so
// future searches can reuse it. Errors are intentionally swallowed --
// caching a Spot (e.g. it may already be known) is a side effect of
// creating a Pin, not the operation the caller asked for, so it must never
// fail Pin creation.
func (useCase *UseCase) saveSpot(ctx context.Context, input CreateInput) {
	placeID, err := spotdomain.NewPlaceID(input.PlaceID)
	if err != nil {
		return
	}
	location, err := spotdomain.NewLocation(input.Latitude, input.Longitude)
	if err != nil {
		return
	}

	_ = useCase.spotUseCase.SaveSpot(ctx, &spotdomain.Spot{
		ID:          uuid.New(),
		PlaceID:     placeID,
		Name:        input.Name,
		Address:     input.Address,
		Location:    location,
		FirstPlanID: input.PlanID,
	})
}

func (useCase *UseCase) UpdatePin(ctx context.Context, planID, pinID uuid.UUID, input UpdateInput) (*domain.Pin, error) {
	pin, err := useCase.repo.FindByID(ctx, planID, pinID)
	if err != nil {
		return nil, err
	}

	if input.Category != nil {
		pin.Category = *input.Category
	}
	if input.Colour != nil {
		pin.Colour = *input.Colour
	}

	if err := useCase.repo.Update(ctx, pin); err != nil {
		return nil, fmt.Errorf("update pin: %w", err)
	}
	return pin, nil
}

func (useCase *UseCase) DeletePin(ctx context.Context, planID, pinID uuid.UUID) error {
	return useCase.repo.Delete(ctx, planID, pinID)
}

func (useCase *UseCase) ListPins(ctx context.Context, planID uuid.UUID) ([]*domain.Pin, error) {
	return useCase.repo.FindByPlanID(ctx, planID)
}
