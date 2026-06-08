package spot_test

import (
	"errors"
	"testing"

	"github.com/Application-drop-up/Travellle/internal/domain/spot"
)

func TestNewPlaceID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{
			name: "valid place ID",
			id:   "ChIJ5eTFBkqLGGARsV4PF3rDVAA",
		},
		{
			name:    "empty place ID",
			id:      "",
			wantErr: spot.ErrEmptyPlaceID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := spot.NewPlaceID(tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewPlaceID(%q) error = %v, want %v", tt.id, err, tt.wantErr)
				return
			}
			if tt.wantErr == nil && got.String() != tt.id {
				t.Errorf("PlaceID.String() = %q, want %q", got.String(), tt.id)
			}
		})
	}
}
