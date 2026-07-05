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

func (repo *PinRepository) Create(ctx context.Context, pin *domain.Pin) error {
	query := `
		INSERT INTO pins (id, plan_id, name, latitude, longitude, category, colour)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at`
	err := repo.db.QueryRowContext(ctx, query,
		pin.ID, pin.PlanID, pin.Name, pin.Latitude, pin.Longitude, string(pin.Category), pin.Colour).
		Scan(&pin.CreatedAt, &pin.UpdatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgFKViolation {
			return plandomain.ErrNotFound
		}
		return fmt.Errorf("insert pin: %w", err)
	}
	return nil
}

func (repo *PinRepository) FindByID(ctx context.Context, planID, pinID uuid.UUID) (*domain.Pin, error) {
	query := `
		SELECT id, plan_id, name, latitude, longitude, category, colour, created_at, updated_at
		FROM pins WHERE id = $1 AND plan_id = $2`
	pin := &domain.Pin{}
	var category string
	err := repo.db.QueryRowContext(ctx, query, pinID, planID).
		Scan(&pin.ID, &pin.PlanID, &pin.Name, &pin.Latitude, &pin.Longitude, &category, &pin.Colour, &pin.CreatedAt, &pin.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find pin by id: %w", err)
	}
	pin.Category = domain.Category(category)
	return pin, nil
}

func (repo *PinRepository) Update(ctx context.Context, pin *domain.Pin) error {
	query := `
		UPDATE pins SET category = $1, colour = $2, updated_at = NOW()
		WHERE id = $3 AND plan_id = $4
		RETURNING updated_at`
	err := repo.db.QueryRowContext(ctx, query, string(pin.Category), pin.Colour, pin.ID, pin.PlanID).
		Scan(&pin.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("update pin: %w", err)
	}
	return nil
}

func (repo *PinRepository) FindByPlanID(ctx context.Context, planID uuid.UUID) ([]*domain.Pin, error) {
	query := `
		SELECT id, plan_id, name, latitude, longitude, category, colour, created_at, updated_at
		FROM pins WHERE plan_id = $1 ORDER BY created_at ASC`
	rows, err := repo.db.QueryContext(ctx, query, planID)
	if err != nil {
		return nil, fmt.Errorf("find pins by plan id: %w", err)
	}
	defer rows.Close()

	var pins []*domain.Pin
	for rows.Next() {
		pin := &domain.Pin{}
		var category string
		if err := rows.Scan(&pin.ID, &pin.PlanID, &pin.Name, &pin.Latitude, &pin.Longitude, &category, &pin.Colour, &pin.CreatedAt, &pin.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan pin: %w", err)
		}
		pin.Category = domain.Category(category)
		pins = append(pins, pin)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("find pins by plan id rows: %w", err)
	}
	return pins, nil
}

func (repo *PinRepository) Delete(ctx context.Context, planID, pinID uuid.UUID) error {
	query := `DELETE FROM pins WHERE id = $1 AND plan_id = $2`
	result, err := repo.db.ExecContext(ctx, query, pinID, planID)
	if err != nil {
		return fmt.Errorf("delete pin: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete pin rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
