package plan

import "context"

type Repository interface {
	Create(ctx context.Context, plan *Plan) error
	FindByShareToken(ctx context.Context, shareToken string) (*Plan, error)
}
