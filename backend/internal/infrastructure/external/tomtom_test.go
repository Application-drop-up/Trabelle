package external_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Application-drop-up/Travellle/internal/domain/spot"
	"github.com/Application-drop-up/Travellle/internal/infrastructure/external"
)

func newTestTomTomClient(t *testing.T, handler http.HandlerFunc) *external.TomTomClient {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return external.NewTomTomClientWithURL(srv.URL, "test-api-key")
}

func TestTomTomClient_Search(t *testing.T) {
	t.Parallel()

	t.Run("returns spots on success", func(t *testing.T) {
		t.Parallel()

		handler := func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("key") == "" {
				t.Error("key query parameter is missing")
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{
					{
						"id":       "826003003559627",
						"poi":      map[string]string{"name": "Tokyo Tower"},
						"address":  map[string]string{"freeformAddress": "4 Chome-2-8 Shibakoen, Minato City, Tokyo"},
						"position": map[string]float64{"lat": 35.6585805, "lon": 139.7454329},
					},
				},
			})
		}

		client := newTestTomTomClient(t, handler)
		spots, err := client.Search(context.Background(), "Tokyo Tower")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(spots) != 1 {
			t.Fatalf("got %d spots, want 1", len(spots))
		}
		if spots[0].Name != "Tokyo Tower" {
			t.Errorf("Name = %q, want %q", spots[0].Name, "Tokyo Tower")
		}
		if spots[0].PlaceID.String() != "826003003559627" {
			t.Errorf("PlaceID = %q, want 826003003559627", spots[0].PlaceID)
		}
	})

	t.Run("empty query returns ErrInvalidQuery", func(t *testing.T) {
		t.Parallel()

		client := external.NewTomTomClient("test-api-key")
		_, err := client.Search(context.Background(), "")
		if !errors.Is(err, spot.ErrInvalidQuery) {
			t.Errorf("got %v, want ErrInvalidQuery", err)
		}
	})

	t.Run("skips result with invalid location", func(t *testing.T) {
		t.Parallel()

		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"results": []map[string]any{
					{
						"id":       "valid-id",
						"poi":      map[string]string{"name": "Valid Place"},
						"address":  map[string]string{"freeformAddress": "somewhere"},
						"position": map[string]float64{"lat": 35.0, "lon": 139.0},
					},
					{
						"id":       "invalid-id",
						"poi":      map[string]string{"name": "Invalid Place"},
						"address":  map[string]string{"freeformAddress": "somewhere"},
						"position": map[string]float64{"lat": 999.0, "lon": 999.0},
					},
				},
			})
		}

		client := newTestTomTomClient(t, handler)
		spots, err := client.Search(context.Background(), "test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(spots) != 1 {
			t.Errorf("got %d spots, want 1 (invalid location should be skipped)", len(spots))
		}
	})

	t.Run("non-200 response returns error", func(t *testing.T) {
		t.Parallel()

		handler := func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
		}

		client := newTestTomTomClient(t, handler)
		_, err := client.Search(context.Background(), "Tokyo")
		if err == nil {
			t.Error("expected error for non-200 response, got nil")
		}
	})
}
