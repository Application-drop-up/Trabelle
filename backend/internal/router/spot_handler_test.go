package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Application-drop-up/Travellle/internal/router"
	"github.com/Application-drop-up/Travellle/internal/testutil"
	"github.com/google/uuid"
)

type searchSpotResult struct {
	PlaceID   string  `json:"place_id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// TestSpotHandler_Search covers the local-database hit path only. The
// external-Searcher fallback path (no local match -> call TomTom) needs a
// real network call with a real API key, so it isn't exercised here -- it's
// covered by the mocked TestUseCase_SearchSpots tests in usecase/spot.
func TestSpotHandler_Search(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"}, false)

	searchReq := func(query string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/spots/search?query="+url.QueryEscape(query), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	t.Run("rejects a missing query parameter", func(t *testing.T) {
		t.Parallel()

		w := searchReq("")
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns a previously-saved Spot from the local database", func(t *testing.T) {
		t.Parallel()

		planID := createTestPlan(t, r, "Spot Search Test Plan")
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM plans WHERE id = $1", planID) })

		uniqueName := "Trabelle Test Spot " + uuid.New().String()
		placeID := uuid.New().String()

		pinResp := createPin(t, r, planID, map[string]any{
			"name": uniqueName, "latitude": 35.0, "longitude": 139.0,
			"category": "sightseeing", "colour": "#FF5733",
			"place_id": placeID, "address": "Test Address",
		})
		if pinResp.Code != http.StatusCreated {
			t.Fatalf("pin creation status = %d, want %d, body: %s", pinResp.Code, http.StatusCreated, pinResp.Body.String())
		}
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM spots WHERE place_id = $1", placeID) })

		w := searchReq(uniqueName)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d, body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		var results []searchSpotResult
		if err := json.Unmarshal(w.Body.Bytes(), &results); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(results) != 1 {
			t.Fatalf("got %d results, want 1: %+v", len(results), results)
		}
		if results[0].PlaceID != placeID {
			t.Errorf("PlaceID = %q, want %q", results[0].PlaceID, placeID)
		}
		if results[0].Name != uniqueName {
			t.Errorf("Name = %q, want %q", results[0].Name, uniqueName)
		}
	})
}
