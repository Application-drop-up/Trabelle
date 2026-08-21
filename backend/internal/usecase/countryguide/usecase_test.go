package countryguide_test

import (
	"context"
	"errors"
	"testing"

	domain "github.com/Application-drop-up/Travellle/internal/domain/countryguide"
	countryguideUseCase "github.com/Application-drop-up/Travellle/internal/usecase/countryguide"
	"github.com/google/uuid"
)

// mockRepository は domain.Repository のテスト用実装
type mockRepository struct {
	guides        []*domain.CountryGuide
	findAllErr    error
	guideByCode   *domain.CountryGuide
	findByCodeErr error
}

func (m *mockRepository) FindAll(_ context.Context) ([]*domain.CountryGuide, error) {
	if m.findAllErr != nil {
		return nil, m.findAllErr
	}
	return m.guides, nil
}

func (m *mockRepository) FindByCode(_ context.Context, _ string) (*domain.CountryGuide, error) {
	if m.guideByCode != nil {
		return m.guideByCode, nil
	}
	if m.findByCodeErr != nil {
		return nil, m.findByCodeErr
	}
	return nil, domain.ErrNotFound
}

func newTestGuide() *domain.CountryGuide {
	return &domain.CountryGuide{
		ID:          uuid.New(),
		CountryCode: "TH",
		CountryName: "Thailand",
		Items: []domain.Item{
			{
				ID:          uuid.New(),
				Category:    domain.CategoryEntryCard,
				Title:       "TDAC",
				Description: "Apply online within 72h before arrival",
				URL:         "https://tdac.immigration.go.th",
				IsMandatory: true,
			},
			{
				ID:       uuid.New(),
				Category: domain.CategoryPackingTip,
				Title:    "Bring a plug adapter",
			},
		},
	}
}

func TestUseCase_ListCountryGuides(t *testing.T) {
	t.Parallel()

	t.Run("returns a DTO collection mapped from the repository", func(t *testing.T) {
		t.Parallel()

		guide := newTestGuide()
		useCase := countryguideUseCase.New(&mockRepository{guides: []*domain.CountryGuide{guide}})

		got, err := useCase.ListCountryGuides(context.Background())
		if err != nil {
			t.Fatalf("ListCountryGuides() unexpected error: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("ListCountryGuides() returned %d guides, want 1", len(got))
		}

		dto := got[0]
		if dto.ID != guide.ID || dto.CountryCode != guide.CountryCode || dto.CountryName != guide.CountryName {
			t.Errorf("ListCountryGuides()[0] = %+v, want fields matching %+v", dto, guide)
		}
		if len(dto.Items) != 2 {
			t.Fatalf("ListCountryGuides()[0].Items has %d items, want 2", len(dto.Items))
		}
		if dto.Items[0].Category != string(domain.CategoryEntryCard) {
			t.Errorf("Items[0].Category = %q, want %q", dto.Items[0].Category, domain.CategoryEntryCard)
		}
		if !dto.Items[0].IsMandatory {
			t.Error("Items[0].IsMandatory = false, want true")
		}
		if dto.Items[1].IsMandatory {
			t.Error("Items[1].IsMandatory = true, want false")
		}
	})

	t.Run("returns an empty collection when there are no guides", func(t *testing.T) {
		t.Parallel()

		useCase := countryguideUseCase.New(&mockRepository{guides: nil})

		got, err := useCase.ListCountryGuides(context.Background())
		if err != nil {
			t.Fatalf("ListCountryGuides() unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("ListCountryGuides() returned %d guides, want 0", len(got))
		}
	})

	t.Run("wraps a repository error", func(t *testing.T) {
		t.Parallel()

		useCase := countryguideUseCase.New(&mockRepository{findAllErr: errors.New("db down")})

		_, err := useCase.ListCountryGuides(context.Background())
		if err == nil {
			t.Fatal("ListCountryGuides() expected error, got nil")
		}
	})
}

func TestUseCase_GetCountryGuide(t *testing.T) {
	t.Parallel()

	t.Run("returns a DTO mapped from the repository", func(t *testing.T) {
		t.Parallel()

		guide := newTestGuide()
		useCase := countryguideUseCase.New(&mockRepository{guideByCode: guide})

		got, err := useCase.GetCountryGuide(context.Background(), "TH")
		if err != nil {
			t.Fatalf("GetCountryGuide() unexpected error: %v", err)
		}
		if got.ID != guide.ID || got.CountryCode != guide.CountryCode || got.CountryName != guide.CountryName {
			t.Errorf("GetCountryGuide() = %+v, want fields matching %+v", got, guide)
		}
		if len(got.Items) != len(guide.Items) {
			t.Errorf("GetCountryGuide() has %d items, want %d", len(got.Items), len(guide.Items))
		}
	})

	t.Run("returns ErrNotFound for an unknown code", func(t *testing.T) {
		t.Parallel()

		useCase := countryguideUseCase.New(&mockRepository{})

		_, err := useCase.GetCountryGuide(context.Background(), "ZZ")
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("GetCountryGuide() error = %v, want %v", err, domain.ErrNotFound)
		}
	})
}
