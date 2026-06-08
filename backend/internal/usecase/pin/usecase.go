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

func (uc *UseCase) CreatePin(ctx context.Context, input CreateInput) (*domain.Pin, error) {
	p := &domain.Pin{
		ID:        uuid.New(),
		PlanID:    input.PlanID,
		Name:      input.Name,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
		Category:  input.Category,
		Colour:    input.Colour,
	}

	if err := uc.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("create pin: %w", err)
	}
	return p, nil
}

func (uc *UseCase) UpdatePin(ctx context.Context, planID, pinID uuid.UUID, input UpdateInput) (*domain.Pin, error) {
	p, err := uc.repo.FindByID(ctx, planID, pinID)
	if err != nil {
		return nil, err
	}

	if input.Category != nil {
		p.Category = *input.Category
	}
	if input.Colour != nil {
		p.Colour = *input.Colour
	}

	if err := uc.repo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("update pin: %w", err)
	}
	return p, nil
}

func (uc *UseCase) DeletePin(ctx context.Context, planID, pinID uuid.UUID) error {
	return uc.repo.Delete(ctx, planID, pinID)
}

func (uc *UseCase) ListPins(ctx context.Context, planID uuid.UUID) ([]*domain.Pin, error) {
	return uc.repo.FindByPlanID(ctx, planID)
}
