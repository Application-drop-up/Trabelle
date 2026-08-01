package pin_test

import (
	"testing"

	"github.com/Application-drop-up/Travellle/internal/domain/pin"
)

func TestCategory_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		category pin.Category
		want     bool
	}{
		{"restaurant", pin.CategoryRestaurant, true},
		{"hotel", pin.CategoryHotel, true},
		{"sightseeing", pin.CategorySightseeing, true},
		{"transport", pin.CategoryTransport, true},
		{"other", pin.CategoryOther, true},
		{"empty string", pin.Category(""), false},
		{"unknown value", pin.Category("museum"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.category.IsValid(); got != tt.want {
				t.Errorf("Category(%q).IsValid() = %v, want %v", tt.category, got, tt.want)
			}
		})
	}
}
