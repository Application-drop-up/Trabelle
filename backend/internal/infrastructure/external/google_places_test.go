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

func newTestClient(t *testing.T, handler http.HandlerFunc) *external.GooglePlacesClient {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return external.NewGooglePlacesClientWithURL(srv.URL, "test-api-key")
}

func TestGooglePlacesClient_Search(t *testing.T) {
	t.Parallel()

	t.Run("returns spots on success", func(t *testing.T) {
		t.Parallel()

		handler := func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Goog-Api-Key") == "" {
				t.Error("X-Goog-Api-Key header is missing")
			}
			if r.Header.Get("X-Goog-FieldMask") == "" {
				t.Error("X-Goog-FieldMask header is missing")
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"places": []map[string]any{
					{
						"id":               "ChIJ5eTFBkqLGGARsV4PF3rDVAA",
						"displayName":      map[string]string{"text": "Tokyo Tower"},
						"formattedAddress": "4 Chome-2-8 Shibakoen, Minato City, Tokyo",
						"location":         map[string]float64{"latitude": 35.6585805, "longitude": 139.7454329},
					},
				},
			})
		}

		client := newTestClient(t, handler)
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
		if spots[0].PlaceID.String() != "ChIJ5eTFBkqLGGARsV4PF3rDVAA" {
			t.Errorf("PlaceID = %q, want ChIJ5eTFBkqLGGARsV4PF3rDVAA", spots[0].PlaceID)
		}
	})

	t.Run("empty query returns ErrInvalidQuery", func(t *testing.T) {
		t.Parallel()

		client := external.NewGooglePlacesClient("test-api-key")
		_, err := client.Search(context.Background(), "")
		if !errors.Is(err, spot.ErrInvalidQuery) {
			t.Errorf("got %v, want ErrInvalidQuery", err)
		}
	})

	t.Run("skips place with invalid location", func(t *testing.T) {
		t.Parallel()

		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"places": []map[string]any{
					{
						"id":               "valid-id",
						"displayName":      map[string]string{"text": "Valid Place"},
						"formattedAddress": "somewhere",
						"location":         map[string]float64{"latitude": 35.0, "longitude": 139.0},
					},
					{
						"id":               "invalid-id",
						"displayName":      map[string]string{"text": "Invalid Place"},
						"formattedAddress": "somewhere",
						"location":         map[string]float64{"latitude": 999.0, "longitude": 999.0},
					},
				},
			})
		}

		client := newTestClient(t, handler)
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

		client := newTestClient(t, handler)
		_, err := client.Search(context.Background(), "Tokyo")
		if err == nil {
			t.Error("expected error for non-200 response, got nil")
		}
	})
}
