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
	IsPublic   bool          `json:"is_public"`
	Pins       []interface{} `json:"pins"`
}

func TestPlanHandler_CreateAndGet(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"}, false)

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

func TestPlanHandler_Publish(t *testing.T) {
	t.Parallel()

	db := testutil.NewTestDB(t)
	r := router.New(db, "test-api-key", []string{"http://localhost:3000"}, false)

	t.Run("makes the plan public", func(t *testing.T) {
		t.Parallel()

		body, _ := json.Marshal(map[string]string{"title": "Trip to Osaka"})
		createReq := httptest.NewRequest(http.MethodPost, "/plans", bytes.NewReader(body))
		createReq.Header.Set("Content-Type", "application/json")
		createW := httptest.NewRecorder()
		r.ServeHTTP(createW, createReq)

		var created planResponse
		if err := json.Unmarshal(createW.Body.Bytes(), &created); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		t.Cleanup(func() { _, _ = db.Exec("DELETE FROM plans WHERE id = $1", created.ID) })

		if created.IsPublic {
			t.Fatal("newly created plan IsPublic = true, want false")
		}

		publishReq := httptest.NewRequest(http.MethodPost, "/api/v1/plans/"+created.ShareToken+"/publish", nil)
		publishW := httptest.NewRecorder()
		r.ServeHTTP(publishW, publishReq)

		if publishW.Code != http.StatusOK {
			t.Fatalf("POST /api/v1/plans/{share_token}/publish status = %d, want %d, body: %s", publishW.Code, http.StatusOK, publishW.Body.String())
		}

		var published planResponse
		if err := json.Unmarshal(publishW.Body.Bytes(), &published); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if !published.IsPublic {
			t.Error("after Publish(), IsPublic = false, want true")
		}

		getReq := httptest.NewRequest(http.MethodGet, "/plans/"+created.ShareToken, nil)
		getW := httptest.NewRecorder()
		r.ServeHTTP(getW, getReq)

		var fetched planResponse
		if err := json.Unmarshal(getW.Body.Bytes(), &fetched); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if !fetched.IsPublic {
			t.Error("after Publish(), GET /plans/{share_token} IsPublic = false, want true")
		}
	})

	t.Run("returns 404 for an unknown share token", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/plans/does-not-exist/publish", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("POST /api/v1/plans/{share_token}/publish (unknown) status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}
