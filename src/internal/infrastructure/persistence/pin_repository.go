package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	plandomain "github.com/Application-drop-up/Travellle/internal/domain/plan"
	domain "github.com/Application-drop-up/Travellle/internal/domain/pin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const pgFKViolation = "23503"

type PinRepository struct {
	db *sql.DB
}

func NewPinRepository(db *sql.DB) *PinRepository {
	return &PinRepository{db: db}
}

func (r *PinRepository) Create(ctx context.Context, p *domain.Pin) error {
	query := `
		INSERT INTO pins (id, plan_id, name, latitude, longitude, category, colour)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at`
	err := r.db.QueryRowContext(ctx, query,
		p.ID, p.PlanID, p.Name, p.Latitude, p.Longitude, string(p.Category), p.Colour).
		Scan(&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgFKViolation {
			return plandomain.ErrNotFound
		}
		return fmt.Errorf("insert pin: %w", err)
	}
	return nil
}

func (r *PinRepository) FindByID(ctx context.Context, planID, pinID uuid.UUID) (*domain.Pin, error) {
	query := `
		SELECT id, plan_id, name, latitude, longitude, category, colour, created_at, updated_at
		FROM pins WHERE id = $1 AND plan_id = $2`
	p := &domain.Pin{}
	var category string
	err := r.db.QueryRowContext(ctx, query, pinID, planID).
		Scan(&p.ID, &p.PlanID, &p.Name, &p.Latitude, &p.Longitude, &category, &p.Colour, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find pin by id: %w", err)
	}
	p.Category = domain.Category(category)
	return p, nil
}

func (r *PinRepository) Update(ctx context.Context, p *domain.Pin) error {
	query := `
		UPDATE pins SET category = $1, colour = $2, updated_at = NOW()
		WHERE id = $3 AND plan_id = $4
		RETURNING updated_at`
	err := r.db.QueryRowContext(ctx, query, string(p.Category), p.Colour, p.ID, p.PlanID).
		Scan(&p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("update pin: %w", err)
	}
	return nil
}

func (r *PinRepository) Delete(ctx context.Context, planID, pinID uuid.UUID) error {
	query := `DELETE FROM pins WHERE id = $1 AND plan_id = $2`
	result, err := r.db.ExecContext(ctx, query, pinID, planID)
	if err != nil {
		return fmt.Errorf("delete pin: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete pin rows affected: %w", err)
	}
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
