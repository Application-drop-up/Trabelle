package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domain "github.com/Application-drop-up/Travellle/internal/domain/planmember"
	plandomain "github.com/Application-drop-up/Travellle/internal/domain/plan"
	userdomain "github.com/Application-drop-up/Travellle/internal/domain/user"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PlanMemberRepository struct {
	db *sql.DB
}

func NewPlanMemberRepository(db *sql.DB) *PlanMemberRepository {
	return &PlanMemberRepository{db: db}
}

func (repo *PlanMemberRepository) Create(ctx context.Context, member *domain.PlanMember) error {
	query := `
		INSERT INTO plan_members (id, plan_id, user_id)
		VALUES ($1, $2, $3)
		RETURNING created_at`
	err := repo.db.QueryRowContext(ctx, query, member.ID, member.PlanID, member.UserID).
		Scan(&member.CreatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch {
			case pqErr.Code == pgUniqueViolation:
				return domain.ErrAlreadyMember
			case pqErr.Code == pgFKViolation && pqErr.Constraint == "plan_members_plan_id_fkey":
				return plandomain.ErrNotFound
			case pqErr.Code == pgFKViolation && pqErr.Constraint == "plan_members_user_id_fkey":
				return userdomain.ErrNotFound
			}
		}
		return fmt.Errorf("insert plan member: %w", err)
	}
	return nil
}

func (repo *PlanMemberRepository) FindByPlanID(ctx context.Context, planID uuid.UUID) ([]*domain.PlanMember, error) {
	query := `
		SELECT id, plan_id, user_id, created_at
		FROM plan_members WHERE plan_id = $1 ORDER BY created_at ASC`
	rows, err := repo.db.QueryContext(ctx, query, planID)
	if err != nil {
		return nil, fmt.Errorf("find plan members by plan id: %w", err)
	}
	defer rows.Close()

	var members []*domain.PlanMember
	for rows.Next() {
		member := &domain.PlanMember{}
		if err := rows.Scan(&member.ID, &member.PlanID, &member.UserID, &member.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan plan member: %w", err)
		}
		members = append(members, member)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("find plan members by plan id rows: %w", err)
	}
	return members, nil
}

func (repo *PlanMemberRepository) Delete(ctx context.Context, planID, userID uuid.UUID) error {
	query := `DELETE FROM plan_members WHERE plan_id = $1 AND user_id = $2`
	result, err := repo.db.ExecContext(ctx, query, planID, userID)
	if err != nil {
		return fmt.Errorf("delete plan member: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete plan member rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
