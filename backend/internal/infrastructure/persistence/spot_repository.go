package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	plandomain "github.com/Application-drop-up/Travellle/internal/domain/plan"
	domain "github.com/Application-drop-up/Travellle/internal/domain/spot"
	"github.com/lib/pq"
)

type SpotRepository struct {
	db *sql.DB
}

func NewSpotRepository(db *sql.DB) *SpotRepository {
	return &SpotRepository{db: db}
}

func (repo *SpotRepository) Save(ctx context.Context, spot *domain.Spot) error {
	query := `
		INSERT INTO spots (id, place_id, name, address, latitude, longitude, first_plan_id, first_plan_is_public)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := repo.db.ExecContext(
		ctx,
		query,
		spot.ID,
		string(spot.PlaceID),
		spot.Name,
		spot.Address,
		spot.Location.Latitude,
		spot.Location.Longitude,
		spot.FirstPlanID,
		spot.FirstPlanIsPublic,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgFKViolation {
			return plandomain.ErrNotFound
		}
		return fmt.Errorf("insert spot: %w", err)
	}
	return nil
}

func (repo *SpotRepository) FindByPlaceID(ctx context.Context, placeID domain.PlaceID) (*domain.Spot, error) {
	query := `
		SELECT id, place_id, name, address, latitude, longitude, first_plan_id, first_plan_is_public
		FROM spots WHERE place_id = $1`

	spot := &domain.Spot{}

	err := repo.db.QueryRowContext(ctx, query, placeID).
		Scan(
			&spot.ID,
			&spot.PlaceID,
			&spot.Name,
			&spot.Address,
			&spot.Location.Latitude,
			&spot.Location.Longitude,
			&spot.FirstPlanID,
			&spot.FirstPlanIsPublic,
		)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("find spot by place id: %w", err)
	}

	return spot, nil
}

func (repo *SpotRepository) Search(ctx context.Context, query string) ([]*domain.Spot, error) {
	sqlQuery := `
	    SELECT id, place_id, name, address, latitude, longitude, first_plan_id, first_plan_is_public
	    FROM spots WHERE name ILIKE $1`

	pattern := "%" + query + "%"

	rows, err := repo.db.QueryContext(ctx, sqlQuery, pattern)
	if err != nil {
		return nil, fmt.Errorf("find pins by plan id: %w", err)
	}
	defer rows.Close()

	var spots []*domain.Spot
	for rows.Next() {
		spot := &domain.Spot{}
		if err := rows.Scan(
			&spot.ID,
			&spot.PlaceID,
			&spot.Name,
			&spot.Address,
			&spot.Location.Latitude,
			&spot.Location.Longitude,
			&spot.FirstPlanID,
			&spot.FirstPlanIsPublic,
		); err != nil {
			return nil, fmt.Errorf("find spots by id rows: %w", err)
		}

		spots = append(spots, spot)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("search spots rows: %w", err)
	}

	return spots, nil
}
