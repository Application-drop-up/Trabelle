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

func (repo *PlanRepository) Create(ctx context.Context, plan *domain.Plan) error {
	query := `
		INSERT INTO plans (id, title, share_token)
		VALUES ($1, $2, $3)
		RETURNING created_at, updated_at`
	if err := repo.db.QueryRowContext(ctx, query, plan.ID, plan.Title, plan.ShareToken).
		Scan(&plan.CreatedAt, &plan.UpdatedAt); err != nil {
		return fmt.Errorf("insert plan: %w", err)
	}
	return nil
}

func (repo *PlanRepository) FindByShareToken(ctx context.Context, token string) (*domain.Plan, error) {
	query := `SELECT id, title, share_token, created_at, updated_at FROM plans WHERE share_token = $1`
	plan := &domain.Plan{}
	err := repo.db.QueryRowContext(ctx, query, token).
		Scan(&plan.ID, &plan.Title, &plan.ShareToken, &plan.CreatedAt, &plan.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find plan by share token: %w", err)
	}
	return plan, nil
}
