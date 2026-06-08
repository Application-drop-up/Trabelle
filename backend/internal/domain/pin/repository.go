package pin

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, pin *Pin) error
	FindByID(ctx context.Context, planID, pinID uuid.UUID) (*Pin, error)
	FindByPlanID(ctx context.Context, planID uuid.UUID) ([]*Pin, error)
	Update(ctx context.Context, pin *Pin) error
	Delete(ctx context.Context, planID, pinID uuid.UUID) error
}
