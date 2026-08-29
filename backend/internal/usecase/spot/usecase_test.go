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
	spots  []*domain.Spot
	err    error
	called bool
}

func (m *mockSearcher) Search(_ context.Context, _ string) ([]*domain.Spot, error) {
	m.called = true
	return m.spots, m.err
}

// mockRepository は domain.Repository のテスト用実装
type mockRepository struct {
	saveErr     error
	searchSpots []*domain.Spot
	searchErr   error
}

func (m *mockRepository) Save(_ context.Context, _ *domain.Spot) error {
	return m.saveErr
}

func (m *mockRepository) FindByPlaceID(_ context.Context, _ domain.PlaceID) (*domain.Spot, error) {
	return nil, domain.ErrNotFound
}

func (m *mockRepository) Search(_ context.Context, _ string) ([]*domain.Spot, error) {
	return m.searchSpots, m.searchErr
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
		name      string
		query     string
		mockSpots []*domain.Spot
		mockErr   error
		wantLen   int
		wantErr   error
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

			useCase := spotuc.New(&mockSearcher{spots: tt.mockSpots, err: tt.mockErr}, &mockRepository{})
			got, err := useCase.SearchSpots(context.Background(), tt.query)

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

	t.Run("returns local results without calling the external searcher", func(t *testing.T) {
		t.Parallel()

		searcher := &mockSearcher{spots: []*domain.Spot{dummySpot}}
		repo := &mockRepository{searchSpots: []*domain.Spot{dummySpot}}
		useCase := spotuc.New(searcher, repo)

		got, err := useCase.SearchSpots(context.Background(), "Tokyo Tower")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("got %d spots, want 1", len(got))
		}
		if searcher.called {
			t.Error("external searcher was called even though local results were found")
		}
	})

	t.Run("propagates local repository search error", func(t *testing.T) {
		t.Parallel()

		useCase := spotuc.New(&mockSearcher{}, &mockRepository{searchErr: errors.New("db error")})

		_, err := useCase.SearchSpots(context.Background(), "Tokyo Tower")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestUseCase_SaveSpot(t *testing.T) {
	t.Parallel()

	location, _ := domain.NewLocation(35.6895, 139.6917)
	placeID, _ := domain.NewPlaceID("ChIJ5eTFBkqLGGARsV4PF3rDVAA")
	dummySpot := &domain.Spot{
		PlaceID:  placeID,
		Name:     "Tokyo Tower",
		Address:  "4 Chome-2-8 Shibakoen, Minato City, Tokyo",
		Location: location,
	}

	t.Run("saves successfully", func(t *testing.T) {
		t.Parallel()

		useCase := spotuc.New(&mockSearcher{}, &mockRepository{})
		if err := useCase.SaveSpot(context.Background(), dummySpot); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		t.Parallel()

		useCase := spotuc.New(&mockSearcher{}, &mockRepository{saveErr: errors.New("db error")})

		err := useCase.SaveSpot(context.Background(), dummySpot)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
