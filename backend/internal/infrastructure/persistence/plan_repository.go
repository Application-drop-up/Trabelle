package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domain "github.com/Application-drop-up/Travellle/internal/domain/plan"
	"github.com/google/uuid"
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
	query := `SELECT id, title, share_token, is_public, created_at, updated_at FROM plans WHERE share_token = $1`
	plan := &domain.Plan{}
	err := repo.db.QueryRowContext(ctx, query, token).
		Scan(&plan.ID, &plan.Title, &plan.ShareToken, &plan.IsPublic, &plan.CreatedAt, &plan.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find plan by share token: %w", err)
	}
	return plan, nil
}

func (repo *PlanRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Plan, error) {
	query := `SELECT id, title, share_token, is_public, created_at, updated_at FROM plans WHERE id = $1`
	plan := &domain.Plan{}
	err := repo.db.QueryRowContext(ctx, query, id).
		Scan(&plan.ID, &plan.Title, &plan.ShareToken, &plan.IsPublic, &plan.CreatedAt, &plan.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find plan by id: %w", err)
	}
	return plan, nil
}

func (repo *PlanRepository) FindByMemberID(ctx context.Context, userID uuid.UUID) ([]*domain.Plan, error) {
	query := `
		SELECT p.id, p.title, p.share_token, p.is_public, p.created_at, p.updated_at
		FROM plans p
		JOIN plan_members pm ON pm.plan_id = p.id
		WHERE pm.user_id = $1
		ORDER BY pm.created_at ASC`
	rows, err := repo.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("find plans by member id: %w", err)
	}
	defer rows.Close()

	var plans []*domain.Plan
	for rows.Next() {
		plan := &domain.Plan{}
		if err := rows.Scan(&plan.ID, &plan.Title, &plan.ShareToken, &plan.IsPublic, &plan.CreatedAt, &plan.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan plan: %w", err)
		}
		plans = append(plans, plan)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("find plans by member id rows: %w", err)
	}
	return plans, nil
}

func (repo *PlanRepository) UpdateVisibility(ctx context.Context, plan *domain.Plan) error {
	query := `
		UPDATE plans SET is_public = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING updated_at`
	err := repo.db.QueryRowContext(ctx, query, plan.IsPublic, plan.ID).Scan(&plan.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("update plan visibility: %w", err)
	}
	return nil
}
