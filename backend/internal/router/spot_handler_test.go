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

type saveSpotResponse struct {
	PlaceID   string  `json:"place_id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func TestSpotHandler_Save(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"}, false)

	planID := createTestPlan(t, r, "Spot Handler Test Plan")
	t.Cleanup(func() { _, _ = db.Exec("DELETE FROM plans WHERE id = $1", planID) })

	saveReq := func(body any) *httptest.ResponseRecorder {
		raw, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/"+uuid.New().String()+"/spot/share", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	t.Run("saves a spot shared to a plan", func(t *testing.T) {
		t.Parallel()

		placeID := uuid.New().String()
		w := saveReq(map[string]any{
			"place_id":  placeID,
			"name":      "Tokyo Tower",
			"address":   "4 Chome-2-8 Shibakoen, Minato City, Tokyo",
			"latitude":  35.6586,
			"longitude": 139.7454,
			"plan_id":   planID,
		})

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d, body: %s", w.Code, http.StatusCreated, w.Body.String())
		}

		var saved saveSpotResponse
		if err := json.Unmarshal(w.Body.Bytes(), &saved); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM spots WHERE place_id = $1", saved.PlaceID) })

		if saved.PlaceID != placeID {
			t.Errorf("PlaceID = %q, want %q", saved.PlaceID, placeID)
		}
		if saved.Name != "Tokyo Tower" {
			t.Errorf("Name = %q, want %q", saved.Name, "Tokyo Tower")
		}
		if saved.Address != "4 Chome-2-8 Shibakoen, Minato City, Tokyo" {
			t.Errorf("Address = %q, want %q", saved.Address, "4 Chome-2-8 Shibakoen, Minato City, Tokyo")
		}
		if saved.Latitude != 35.6586 || saved.Longitude != 139.7454 {
			t.Errorf("Location = (%v, %v), want (35.6586, 139.7454)", saved.Latitude, saved.Longitude)
		}
	})

	t.Run("rejects an invalid user id", func(t *testing.T) {
		t.Parallel()

		raw, _ := json.Marshal(map[string]any{
			"place_id": uuid.New().String(), "name": "Spot", "address": "Addr",
			"latitude": 1.0, "longitude": 1.0, "plan_id": planID,
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/not-a-uuid/spot/share", bytes.NewReader(raw))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("rejects an invalid request body", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/user/"+uuid.New().String()+"/spot/share", bytes.NewReader([]byte("not json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("rejects a missing name", func(t *testing.T) {
		t.Parallel()

		w := saveReq(map[string]any{
			"place_id": uuid.New().String(), "name": "", "address": "Addr",
			"latitude": 1.0, "longitude": 1.0, "plan_id": planID,
		})

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("rejects a missing address", func(t *testing.T) {
		t.Parallel()

		w := saveReq(map[string]any{
			"place_id": uuid.New().String(), "name": "Spot", "address": "",
			"latitude": 1.0, "longitude": 1.0, "plan_id": planID,
		})

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("rejects an invalid plan_id", func(t *testing.T) {
		t.Parallel()

		w := saveReq(map[string]any{
			"place_id": uuid.New().String(), "name": "Spot", "address": "Addr",
			"latitude": 1.0, "longitude": 1.0, "plan_id": "not-a-uuid",
		})

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("rejects an empty place_id", func(t *testing.T) {
		t.Parallel()

		w := saveReq(map[string]any{
			"place_id": "", "name": "Spot", "address": "Addr",
			"latitude": 1.0, "longitude": 1.0, "plan_id": planID,
		})

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("rejects an out-of-range location", func(t *testing.T) {
		t.Parallel()

		w := saveReq(map[string]any{
			"place_id": uuid.New().String(), "name": "Spot", "address": "Addr",
			"latitude": 999.0, "longitude": 1.0, "plan_id": planID,
		})

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns 404 when the plan does not exist", func(t *testing.T) {
		t.Parallel()

		w := saveReq(map[string]any{
			"place_id": uuid.New().String(), "name": "Spot", "address": "Addr",
			"latitude": 1.0, "longitude": 1.0, "plan_id": uuid.New().String(),
		})

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d, body: %s", w.Code, http.StatusNotFound, w.Body.String())
		}
	})
}
