package router_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/Application-drop-up/Travellle/internal/router"
	"github.com/Application-drop-up/Travellle/internal/testutil"
)

type countryGuideItemResponse struct {
	ID          string `json:"id"`
	Category    string `json:"category"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	IsMandatory bool   `json:"is_mandatory"`
}

type countryGuideResponse struct {
	ID          string                     `json:"id"`
	CountryCode string                     `json:"country_code"`
	CountryName string                     `json:"country_name"`
	Items       []countryGuideItemResponse `json:"items"`
}

func newTestCountryGuideForHandler(t *testing.T, db *sql.DB, countryCode, countryName string) uuid.UUID {
	t.Helper()

	guideID := uuid.New()
	_, err := db.Exec(`
		INSERT INTO country_guides (id, country_code, country_name)
		VALUES ($1, $2, $3)`,
		guideID, countryCode, countryName)
	if err != nil {
		t.Fatalf("failed to insert prerequisite country guide: %v", err)
	}
	t.Cleanup(func() { _, _ = db.Exec("DELETE FROM country_guides WHERE id = $1", guideID) })

	_, err = db.Exec(`
		INSERT INTO country_guide_items (country_guide_id, category, title, description, url, is_mandatory, sort_order)
		VALUES ($1, 'entry_card', 'Arrival card', 'Apply online', 'https://example.com', true, 0)`,
		guideID)
	if err != nil {
		t.Fatalf("failed to insert prerequisite country guide item: %v", err)
	}

	return guideID
}

func TestCountryGuideHandler_List(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"}, false)

	guideID := newTestCountryGuideForHandler(t, db, "ZZ", "Zzedland")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/country-guides", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /api/v1/country-guides status = %d, want %d, body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	var got []countryGuideResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	var found *countryGuideResponse
	for i, guide := range got {
		if guide.ID == guideID.String() {
			found = &got[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("GET /api/v1/country-guides did not return seeded guide %s", guideID)
	}
	if found.CountryCode != "ZZ" || found.CountryName != "Zzedland" {
		t.Errorf("guide = %+v, want CountryCode=ZZ CountryName=Zzedland", found)
	}
	if len(found.Items) != 1 || found.Items[0].Category != "entry_card" || !found.Items[0].IsMandatory {
		t.Errorf("guide items = %+v, want one mandatory entry_card item", found.Items)
	}
}

func TestCountryGuideHandler_GetByCode(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"}, false)

	t.Run("returns the guide when it exists", func(t *testing.T) {
		t.Parallel()

		guideID := newTestCountryGuideForHandler(t, db, "YY", "Yyland")

		req := httptest.NewRequest(http.MethodGet, "/api/v1/country-guides/YY", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("GET /api/v1/country-guides/YY status = %d, want %d, body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		var got countryGuideResponse
		if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got.ID != guideID.String() || got.CountryName != "Yyland" {
			t.Errorf("got = %+v, want ID=%s CountryName=Yyland", got, guideID)
		}
	})

	t.Run("returns 404 for an unknown code", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/country-guides/XX", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("GET /api/v1/country-guides/XX status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}
