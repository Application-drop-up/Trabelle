package planmember

import (
	"context"
	"fmt"

	domain "github.com/Application-drop-up/Travellle/internal/domain/planmember"
	"github.com/google/uuid"
)

type UseCase struct {
	repo domain.Repository
}

func New(repo domain.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (useCase *UseCase) AddMember(ctx context.Context, planID, userID uuid.UUID) (*domain.PlanMember, error) {
	member := &domain.PlanMember{
		ID:     uuid.New(),
		PlanID: planID,
		UserID: userID,
	}
	if err := useCase.repo.Create(ctx, member); err != nil {
		return nil, fmt.Errorf("add plan member: %w", err)
	}
	return member, nil
}

func (useCase *UseCase) ListMembers(ctx context.Context, planID uuid.UUID) ([]*domain.PlanMember, error) {
	return useCase.repo.FindByPlanID(ctx, planID)
}

func (useCase *UseCase) RemoveMember(ctx context.Context, planID, userID uuid.UUID) error {
	if err := useCase.repo.Delete(ctx, planID, userID); err != nil {
		return fmt.Errorf("remove plan member: %w", err)
	}
	return nil
}
