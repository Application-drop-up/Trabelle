package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Application-drop-up/Travellle/internal/router"
	"github.com/Application-drop-up/Travellle/internal/testutil"
)

type planResponse struct {
	ID         string        `json:"id"`
	ShareToken string        `json:"share_token"`
	Title      string        `json:"title"`
	Pins       []interface{} `json:"pins"`
}

func TestPlanHandler_CreateAndGet(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"})

	t.Run("creates a plan and returns it", func(t *testing.T) {
		t.Parallel()

		body, _ := json.Marshal(map[string]string{"title": "Trip to Kyoto"})
		req := httptest.NewRequest(http.MethodPost, "/plans", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("POST /plans status = %d, want %d, body: %s", w.Code, http.StatusCreated, w.Body.String())
		}

		var created planResponse
		if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM plans WHERE id = $1", created.ID) })

		if created.Title != "Trip to Kyoto" {
			t.Errorf("Title = %q, want %q", created.Title, "Trip to Kyoto")
		}
		if created.ShareToken == "" {
			t.Error("ShareToken is empty")
		}
		if len(created.Pins) != 0 {
			t.Errorf("Pins = %v, want empty", created.Pins)
		}

		t.Run("fetches the plan by share token", func(t *testing.T) {
			getReq := httptest.NewRequest(http.MethodGet, "/plans/"+created.ShareToken, nil)
			getW := httptest.NewRecorder()

			r.ServeHTTP(getW, getReq)

			if getW.Code != http.StatusOK {
				t.Fatalf("GET /plans/{share_token} status = %d, want %d, body: %s", getW.Code, http.StatusOK, getW.Body.String())
			}

			var fetched planResponse
			if err := json.Unmarshal(getW.Body.Bytes(), &fetched); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if fetched.ID != created.ID {
				t.Errorf("ID = %q, want %q", fetched.ID, created.ID)
			}
			if fetched.Title != created.Title {
				t.Errorf("Title = %q, want %q", fetched.Title, created.Title)
			}
		})
	})

	t.Run("rejects an empty title", func(t *testing.T) {
		t.Parallel()

		body, _ := json.Marshal(map[string]string{"title": ""})
		req := httptest.NewRequest(http.MethodPost, "/plans", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("POST /plans (empty title) status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("returns 404 for an unknown share token", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/plans/does-not-exist", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("GET /plans/{share_token} (unknown) status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}
