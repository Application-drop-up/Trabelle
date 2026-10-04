package planmember

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, member *PlanMember) error
	FindByPlanID(ctx context.Context, planID uuid.UUID) ([]*PlanMember, error)
	Delete(ctx context.Context, planID, userID uuid.UUID) error
}
