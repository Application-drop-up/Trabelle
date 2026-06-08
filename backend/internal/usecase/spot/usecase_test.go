package spot_test

import (
	"context"
	"errors"
	"testing"

	domain "github.com/Application-drop-up/Travellle/internal/domain/spot"
	spotuc "github.com/Application-drop-up/Travellle/internal/usecase/spot"
)

// mockSearcher は domain.Searcher のテスト用実装
type mockSearcher struct {
	spots []*domain.Spot
	err   error
}

func (m *mockSearcher) Search(_ context.Context, _ string) ([]*domain.Spot, error) {
	return m.spots, m.err
}

func TestUseCase_SearchSpots(t *testing.T) {
	t.Parallel()

	location, _ := domain.NewLocation(35.6895, 139.6917)
	placeID, _ := domain.NewPlaceID("ChIJ5eTFBkqLGGARsV4PF3rDVAA")
	dummySpot := &domain.Spot{
		PlaceID:  placeID,
		Name:     "Tokyo Tower",
		Address:  "4 Chome-2-8 Shibakoen, Minato City, Tokyo",
		Location: location,
	}

	tests := []struct {
		name       string
		query      string
		mockSpots  []*domain.Spot
		mockErr    error
		wantLen    int
		wantErr    error
	}{
		{
			name:      "returns spots on success",
			query:     "Tokyo Tower",
			mockSpots: []*domain.Spot{dummySpot},
			wantLen:   1,
		},
		{
			name:    "empty query returns ErrInvalidQuery",
			query:   "",
			wantErr: domain.ErrInvalidQuery,
		},
		{
			name:    "propagates searcher error",
			query:   "Tokyo Tower",
			mockErr: errors.New("api error"),
			wantErr: errors.New("search spots: api error"),
		},
		{
			name:      "returns empty slice when no results",
			query:     "nonexistent place xyz",
			mockSpots: []*domain.Spot{},
			wantLen:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := spotuc.New(&mockSearcher{spots: tt.mockSpots, err: tt.mockErr})
			got, err := uc.SearchSpots(context.Background(), tt.query)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != tt.wantLen {
				t.Errorf("got %d spots, want %d", len(got), tt.wantLen)
			}
		})
	}
}
