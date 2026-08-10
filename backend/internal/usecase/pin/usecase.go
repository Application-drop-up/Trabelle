package pin

import (
	"context"
	"fmt"

	domain "github.com/Application-drop-up/Travellle/internal/domain/pin"
	"github.com/google/uuid"
)

type UseCase struct {
	repo domain.Repository
}

func New(repo domain.Repository) *UseCase {
	return &UseCase{repo: repo}
}

type CreateInput struct {
	PlanID    uuid.UUID
	Name      string
	Latitude  float64
	Longitude float64
	Category  domain.Category
	Colour    string
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
	return pin, nil
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
