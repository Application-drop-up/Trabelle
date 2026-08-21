package persistence_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"

	domain "github.com/Application-drop-up/Travellle/internal/domain/countryguide"
	"github.com/Application-drop-up/Travellle/internal/infrastructure/persistence"
	"github.com/Application-drop-up/Travellle/internal/testutil"
)

func newTestCountryGuide(t *testing.T, conn *sql.DB, countryCode, countryName string) uuid.UUID {
	t.Helper()

	guideID := uuid.New()
	_, err := conn.Exec(`
		INSERT INTO country_guides (id, country_code, country_name)
		VALUES ($1, $2, $3)`,
		guideID, countryCode, countryName)
	if err != nil {
		t.Fatalf("failed to insert prerequisite country guide: %v", err)
	}
	t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM country_guides WHERE id = $1", guideID) })
	return guideID
}

func newTestCountryGuideItem(t *testing.T, conn *sql.DB, guideID uuid.UUID, category, title string, sortOrder int) {
	t.Helper()

	_, err := conn.Exec(`
		INSERT INTO country_guide_items (country_guide_id, category, title, description, url, is_mandatory, sort_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		guideID, category, title, "some description", "https://example.com", true, sortOrder)
	if err != nil {
		t.Fatalf("failed to insert prerequisite country guide item: %v", err)
	}
}

func TestCountryGuideRepository_FindAll(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	repo := persistence.NewCountryGuideRepository(conn)

	guideID := newTestCountryGuide(t, conn, "ZZ", "Zzedland")
	newTestCountryGuideItem(t, conn, guideID, string(domain.CategoryEntryCard), "Arrival card", 1)
	newTestCountryGuideItem(t, conn, guideID, string(domain.CategoryPackingTip), "Bring an adapter", 0)

	got, err := repo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("FindAll() unexpected error: %v", err)
	}

	var found *domain.CountryGuide
	for _, guide := range got {
		if guide.ID == guideID {
			found = guide
			break
		}
	}
	if found == nil {
		t.Fatalf("FindAll() did not return seeded guide %s", guideID)
	}
	if found.CountryCode != "ZZ" || found.CountryName != "Zzedland" {
		t.Errorf("FindAll() guide = %+v, want CountryCode=ZZ CountryName=Zzedland", found)
	}
	if len(found.Items) != 2 {
		t.Fatalf("FindAll() guide has %d items, want 2", len(found.Items))
	}
	if found.Items[0].Title != "Bring an adapter" || found.Items[1].Title != "Arrival card" {
		t.Errorf("FindAll() items not ordered by sort_order: %+v", found.Items)
	}
	if !found.Items[0].IsMandatory {
		t.Error("FindAll() item IsMandatory not populated")
	}
}

func TestCountryGuideRepository_FindByCode(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	repo := persistence.NewCountryGuideRepository(conn)

	t.Run("returns the guide when it exists", func(t *testing.T) {
		t.Parallel()

		guideID := newTestCountryGuide(t, conn, "YY", "Yyland")
		newTestCountryGuideItem(t, conn, guideID, string(domain.CategorySIMRecommendation), "Buy a SIM at the airport", 0)

		got, err := repo.FindByCode(context.Background(), "YY")
		if err != nil {
			t.Fatalf("FindByCode() unexpected error: %v", err)
		}
		if got.ID != guideID || got.CountryName != "Yyland" {
			t.Errorf("FindByCode() = %+v, want ID=%s CountryName=Yyland", got, guideID)
		}
		if len(got.Items) != 1 || got.Items[0].Category != domain.CategorySIMRecommendation {
			t.Errorf("FindByCode() items = %+v, want one sim_recommendation item", got.Items)
		}
	})

	t.Run("returns ErrNotFound for an unknown code", func(t *testing.T) {
		t.Parallel()

		_, err := repo.FindByCode(context.Background(), "XX")
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("FindByCode() error = %v, want %v", err, domain.ErrNotFound)
		}
	})
}
