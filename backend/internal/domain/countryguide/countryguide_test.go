package countryguide_test

import (
	"errors"
	"testing"

	"github.com/Application-drop-up/Travellle/internal/domain/countryguide"
	"github.com/google/uuid"
)

func TestNewCountryGuide(t *testing.T) {
	t.Parallel()

	validItems := []countryguide.Item{
		{Category: countryguide.CategoryEntryCard, Title: "Digital Arrival Card", URL: "https://example.com"},
		{Category: countryguide.CategoryPackingTip, Title: "Bring a plug adapter"},
	}

	tests := []struct {
		name        string
		countryCode string
		countryName string
		items       []countryguide.Item
		wantErr     error
	}{
		{
			name:        "valid guide",
			countryCode: "TH",
			countryName: "Thailand",
			items:       validItems,
		},
		{
			name:        "valid guide with no items",
			countryCode: "TH",
			countryName: "Thailand",
			items:       nil,
		},
		{
			name:        "country code too short",
			countryCode: "T",
			countryName: "Thailand",
			items:       validItems,
			wantErr:     countryguide.ErrInvalidCountryCode,
		},
		{
			name:        "country code too long",
			countryCode: "THA",
			countryName: "Thailand",
			items:       validItems,
			wantErr:     countryguide.ErrInvalidCountryCode,
		},
		{
			name:        "country code lowercase",
			countryCode: "th",
			countryName: "Thailand",
			items:       validItems,
			wantErr:     countryguide.ErrInvalidCountryCode,
		},
		{
			name:        "country code digits",
			countryCode: "12",
			countryName: "Thailand",
			items:       validItems,
			wantErr:     countryguide.ErrInvalidCountryCode,
		},
		{
			name:        "empty country name",
			countryCode: "TH",
			countryName: "",
			items:       validItems,
			wantErr:     countryguide.ErrEmptyCountryName,
		},
		{
			name:        "item with empty title",
			countryCode: "TH",
			countryName: "Thailand",
			items: []countryguide.Item{
				{Category: countryguide.CategoryEntryCard, Title: ""},
			},
			wantErr: countryguide.ErrEmptyItemTitle,
		},
		{
			name:        "item with invalid category",
			countryCode: "TH",
			countryName: "Thailand",
			items: []countryguide.Item{
				{Category: "not_a_real_category", Title: "Something"},
			},
			wantErr: countryguide.ErrInvalidItemCategory,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := countryguide.NewCountryGuide(tt.countryCode, tt.countryName, tt.items)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewCountryGuide() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}

			if got.CountryCode != tt.countryCode {
				t.Errorf("CountryCode = %q, want %q", got.CountryCode, tt.countryCode)
			}
			if got.CountryName != tt.countryName {
				t.Errorf("CountryName = %q, want %q", got.CountryName, tt.countryName)
			}
			if len(got.Items) != len(tt.items) {
				t.Errorf("len(Items) = %d, want %d", len(got.Items), len(tt.items))
			}
			if got.ID == uuid.Nil {
				t.Error("ID was not generated")
			}
			for i, item := range got.Items {
				if item.ID == uuid.Nil {
					t.Errorf("Items[%d].ID was not generated", i)
				}
			}
		})
	}
}
