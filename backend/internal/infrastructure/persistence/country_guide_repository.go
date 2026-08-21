package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domain "github.com/Application-drop-up/Travellle/internal/domain/countryguide"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CountryGuideRepository struct {
	db *sql.DB
}

func NewCountryGuideRepository(db *sql.DB) *CountryGuideRepository {
	return &CountryGuideRepository{db: db}
}

func (repo *CountryGuideRepository) FindAll(ctx context.Context) ([]*domain.CountryGuide, error) {
	query := `
		SELECT id, country_code, country_name
		FROM country_guides ORDER BY country_name ASC`
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("find all country guides: %w", err)
	}
	defer rows.Close()

	var guides []*domain.CountryGuide
	for rows.Next() {
		guide := &domain.CountryGuide{}
		if err := rows.Scan(&guide.ID, &guide.CountryCode, &guide.CountryName); err != nil {
			return nil, fmt.Errorf("scan country guide: %w", err)
		}
		guides = append(guides, guide)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("find all country guides rows: %w", err)
	}

	items, err := repo.findItemsByGuideIDs(ctx, guideIDs(guides))
	if err != nil {
		return nil, err
	}
	for _, guide := range guides {
		guide.Items = items[guide.ID]
	}

	return guides, nil
}

func (repo *CountryGuideRepository) FindByCode(ctx context.Context, code string) (*domain.CountryGuide, error) {
	query := `SELECT id, country_code, country_name FROM country_guides WHERE country_code = $1`
	guide := &domain.CountryGuide{}
	err := repo.db.QueryRowContext(ctx, query, code).
		Scan(&guide.ID, &guide.CountryCode, &guide.CountryName)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find country guide by code: %w", err)
	}

	items, err := repo.findItemsByGuideIDs(ctx, []uuid.UUID{guide.ID})
	if err != nil {
		return nil, err
	}
	guide.Items = items[guide.ID]

	return guide, nil
}

func (repo *CountryGuideRepository) findItemsByGuideIDs(ctx context.Context, guideIDs []uuid.UUID) (map[uuid.UUID][]domain.Item, error) {
	items := make(map[uuid.UUID][]domain.Item)
	if len(guideIDs) == 0 {
		return items, nil
	}

	query := `
		SELECT id, country_guide_id, category, title, description, url, is_mandatory
		FROM country_guide_items
		WHERE country_guide_id = ANY($1::uuid[])
		ORDER BY sort_order ASC`
	ids := make([]string, len(guideIDs))
	for i, guideID := range guideIDs {
		ids[i] = guideID.String()
	}
	rows, err := repo.db.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("find country guide items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var guideID uuid.UUID
		var item domain.Item
		var category string
		var description, url sql.NullString
		if err := rows.Scan(&item.ID, &guideID, &category, &item.Title, &description, &url, &item.IsMandatory); err != nil {
			return nil, fmt.Errorf("scan country guide item: %w", err)
		}
		item.Category = domain.ItemCategory(category)
		item.Description = description.String
		item.URL = url.String
		items[guideID] = append(items[guideID], item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("find country guide items rows: %w", err)
	}

	return items, nil
}

func guideIDs(guides []*domain.CountryGuide) []uuid.UUID {
	ids := make([]uuid.UUID, len(guides))
	for i, guide := range guides {
		ids[i] = guide.ID
	}
	return ids
}
