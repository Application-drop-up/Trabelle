package plan

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, plan *Plan) error
	FindByShareToken(ctx context.Context, shareToken string) (*Plan, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Plan, error)
	UpdateVisibility(ctx context.Context, plan *Plan) error
}
