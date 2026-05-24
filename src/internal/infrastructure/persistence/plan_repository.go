package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domain "github.com/Application-drop-up/Travellle/internal/domain/plan"
)

type PlanRepository struct {
	db *sql.DB
}

func NewPlanRepository(db *sql.DB) *PlanRepository {
	return &PlanRepository{db: db}
}

func (r *PlanRepository) Create(ctx context.Context, p *domain.Plan) error {
	query := `
		INSERT INTO plans (id, title, share_token)
		VALUES ($1, $2, $3)
		RETURNING created_at, updated_at`
	if err := r.db.QueryRowContext(ctx, query, p.ID, p.Title, p.ShareToken).
		Scan(&p.CreatedAt, &p.UpdatedAt); err != nil {
		return fmt.Errorf("insert plan: %w", err)
	}
	return nil
}

func (r *PlanRepository) FindByShareToken(ctx context.Context, token string) (*domain.Plan, error) {
	query := `SELECT id, title, share_token, created_at, updated_at FROM plans WHERE share_token = $1`
	p := &domain.Plan{}
	err := r.db.QueryRowContext(ctx, query, token).
		Scan(&p.ID, &p.Title, &p.ShareToken, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find plan by share token: %w", err)
	}
	return p, nil
}
