package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Application-drop-up/Travellle/internal/router"
	"github.com/Application-drop-up/Travellle/internal/testutil"
	"github.com/google/uuid"
)

type pinResponse struct {
	ID       string `json:"id"`
	PlanID   string `json:"plan_id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Colour   string `json:"colour"`
}

// createTestPlan creates a plan via the router and returns its ID. The
// response type is declared locally to avoid colliding with the
// planResponse type used by Plan handler tests.
func createTestPlan(t *testing.T, r http.Handler, title string) string {
	t.Helper()

	body, _ := json.Marshal(map[string]string{"title": title})
	req := httptest.NewRequest(http.MethodPost, "/plans", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("POST /plans status = %d, want %d, body: %s", w.Code, http.StatusCreated, w.Body.String())
	}

	var plan struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &plan); err != nil {
		t.Fatalf("failed to decode plan response: %v", err)
	}
	return plan.ID
}

func createTestPin(t *testing.T, r http.Handler, planID string) pinResponse {
	t.Helper()

	body, _ := json.Marshal(map[string]any{
		"name":      "Tokyo Tower",
		"latitude":  35.6586,
		"longitude": 139.7454,
		"category":  "sightseeing",
		"colour":    "#FF5733",
	})
	req := httptest.NewRequest(http.MethodPost, "/plans/"+planID+"/pins", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("POST /plans/{plan_id}/pins status = %d, want %d, body: %s", w.Code, http.StatusCreated, w.Body.String())
	}

	var pin pinResponse
	if err := json.Unmarshal(w.Body.Bytes(), &pin); err != nil {
		t.Fatalf("failed to decode pin response: %v", err)
	}
	return pin
}

func TestPinHandler_CreateListUpdateDelete(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"}, false)

	planID := createTestPlan(t, r, "Pin Handler Test Plan")
	t.Cleanup(func() { _, _ = db.Exec("DELETE FROM plans WHERE id = $1", planID) })

	t.Run("rejects an invalid category", func(t *testing.T) {
		t.Parallel()

		body, _ := json.Marshal(map[string]any{
			"name": "Bad Pin", "latitude": 1.0, "longitude": 1.0,
			"category": "museum", "colour": "#000000",
		})
		req := httptest.NewRequest(http.MethodPost, "/plans/"+planID+"/pins", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns 404 for a nonexistent plan", func(t *testing.T) {
		t.Parallel()

		body, _ := json.Marshal(map[string]any{
			"name": "Pin", "latitude": 1.0, "longitude": 1.0,
			"category": "other", "colour": "#000000",
		})
		req := httptest.NewRequest(http.MethodPost, "/plans/00000000-0000-0000-0000-000000000000/pins", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d, body: %s", w.Code, http.StatusNotFound, w.Body.String())
		}
	})

	t.Run("creates, lists, updates, and deletes a pin", func(t *testing.T) {
		pin := createTestPin(t, r, planID)

		listReq := httptest.NewRequest(http.MethodGet, "/plans/"+planID+"/pins", nil)
		listW := httptest.NewRecorder()
		r.ServeHTTP(listW, listReq)
		if listW.Code != http.StatusOK {
			t.Fatalf("GET pins status = %d, want %d", listW.Code, http.StatusOK)
		}
		var pins []pinResponse
		if err := json.Unmarshal(listW.Body.Bytes(), &pins); err != nil {
			t.Fatalf("failed to decode list response: %v", err)
		}
		if len(pins) != 1 || pins[0].ID != pin.ID {
			t.Errorf("List() = %+v, want a single pin with ID %q", pins, pin.ID)
		}

		updateBody, _ := json.Marshal(map[string]string{"category": "restaurant", "colour": "#00FF00"})
		updateReq := httptest.NewRequest(http.MethodPatch, "/plans/"+planID+"/pins/"+pin.ID, bytes.NewReader(updateBody))
		updateReq.Header.Set("Content-Type", "application/json")
		updateW := httptest.NewRecorder()
		r.ServeHTTP(updateW, updateReq)
		if updateW.Code != http.StatusOK {
			t.Fatalf("PATCH pin status = %d, want %d, body: %s", updateW.Code, http.StatusOK, updateW.Body.String())
		}
		var updated pinResponse
		if err := json.Unmarshal(updateW.Body.Bytes(), &updated); err != nil {
			t.Fatalf("failed to decode update response: %v", err)
		}
		if updated.Category != "restaurant" || updated.Colour != "#00FF00" {
			t.Errorf("after update, pin = %+v, want Category=restaurant Colour=#00FF00", updated)
		}

		deleteReq := httptest.NewRequest(http.MethodDelete, "/plans/"+planID+"/pins/"+pin.ID, nil)
		deleteW := httptest.NewRecorder()
		r.ServeHTTP(deleteW, deleteReq)
		if deleteW.Code != http.StatusNoContent {
			t.Fatalf("DELETE pin status = %d, want %d", deleteW.Code, http.StatusNoContent)
		}

		deleteAgainReq := httptest.NewRequest(http.MethodDelete, "/plans/"+planID+"/pins/"+pin.ID, nil)
		deleteAgainW := httptest.NewRecorder()
		r.ServeHTTP(deleteAgainW, deleteAgainReq)
		if deleteAgainW.Code != http.StatusNotFound {
			t.Errorf("DELETE (already deleted) status = %d, want %d", deleteAgainW.Code, http.StatusNotFound)
		}
	})
}

func createPin(t *testing.T, r http.Handler, planID string, body map[string]any) *httptest.ResponseRecorder {
	t.Helper()

	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/plans/"+planID+"/pins", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestPinHandler_Create_SavesSpot(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"}, false)

	planID := createTestPlan(t, r, "Pin Spot Persistence Test Plan")
	t.Cleanup(func() { _, _ = db.Exec("DELETE FROM plans WHERE id = $1", planID) })

	t.Run("saves the Spot when place_id is provided", func(t *testing.T) {
		t.Parallel()

		placeID := uuid.New().String()
		w := createPin(t, r, planID, map[string]any{
			"name": "Tokyo Tower", "latitude": 35.6586, "longitude": 139.7454,
			"category": "sightseeing", "colour": "#FF5733",
			"place_id": placeID, "address": "4 Chome-2-8 Shibakoen, Minato City, Tokyo",
		})

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d, body: %s", w.Code, http.StatusCreated, w.Body.String())
		}
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM spots WHERE place_id = $1", placeID) })

		var name, address string
		err := db.QueryRow("SELECT name, address FROM spots WHERE place_id = $1", placeID).Scan(&name, &address)
		if err != nil {
			t.Fatalf("expected a Spot to be saved for place_id %q: %v", placeID, err)
		}
		if name != "Tokyo Tower" {
			t.Errorf("Name = %q, want %q", name, "Tokyo Tower")
		}
		if address != "4 Chome-2-8 Shibakoen, Minato City, Tokyo" {
			t.Errorf("Address = %q, want %q", address, "4 Chome-2-8 Shibakoen, Minato City, Tokyo")
		}
	})

	t.Run("does not save a Spot when place_id is absent", func(t *testing.T) {
		t.Parallel()

		// Uses its own Plan so the spots count below isn't affected by
		// sibling subtests inserting spots under the shared planID above.
		isolatedPlanID := createTestPlan(t, r, "Pin Spot Persistence Isolated Plan")
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM plans WHERE id = $1", isolatedPlanID) })

		w := createPin(t, r, isolatedPlanID, map[string]any{
			"name": "No Place Pin", "latitude": 1.0, "longitude": 1.0,
			"category": "other", "colour": "#000000",
		})

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d, body: %s", w.Code, http.StatusCreated, w.Body.String())
		}

		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM spots WHERE first_plan_id = $1", isolatedPlanID).Scan(&count); err != nil {
			t.Fatalf("failed to count spots: %v", err)
		}
		if count != 0 {
			t.Errorf("spots count = %d, want 0", count)
		}
	})

	t.Run("creates the Pin even when the Spot already exists", func(t *testing.T) {
		t.Parallel()

		placeID := uuid.New().String()
		body := map[string]any{
			"name": "Duplicate Spot", "latitude": 2.0, "longitude": 2.0,
			"category": "restaurant", "colour": "#123456",
			"place_id": placeID, "address": "Some Address",
		}

		first := createPin(t, r, planID, body)
		if first.Code != http.StatusCreated {
			t.Fatalf("first pin status = %d, want %d, body: %s", first.Code, http.StatusCreated, first.Body.String())
		}
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM spots WHERE place_id = $1", placeID) })

		second := createPin(t, r, planID, body)
		if second.Code != http.StatusCreated {
			t.Errorf("second pin status = %d, want %d, body: %s", second.Code, http.StatusCreated, second.Body.String())
		}
	})
}
