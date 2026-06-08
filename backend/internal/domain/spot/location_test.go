package spot_test

import (
	"errors"
	"testing"

	"github.com/Application-drop-up/Travellle/internal/domain/spot"
)

func TestNewLocation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		lat     float64
		lng     float64
		wantErr error
	}{
		{
			name: "valid location",
			lat:  35.6895,
			lng:  139.6917,
		},
		{
			name: "boundary: lat=90 lng=180",
			lat:  90,
			lng:  180,
		},
		{
			name: "boundary: lat=-90 lng=-180",
			lat:  -90,
			lng:  -180,
		},
		{
			name:    "latitude too high",
			lat:     90.0001,
			lng:     0,
			wantErr: spot.ErrInvalidLocation,
		},
		{
			name:    "latitude too low",
			lat:     -90.0001,
			lng:     0,
			wantErr: spot.ErrInvalidLocation,
		},
		{
			name:    "longitude too high",
			lat:     0,
			lng:     180.0001,
			wantErr: spot.ErrInvalidLocation,
		},
		{
			name:    "longitude too low",
			lat:     0,
			lng:     -180.0001,
			wantErr: spot.ErrInvalidLocation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := spot.NewLocation(tt.lat, tt.lng)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewLocation(%v, %v) error = %v, want %v", tt.lat, tt.lng, err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				if got.Latitude != tt.lat || got.Longitude != tt.lng {
					t.Errorf("NewLocation() = {%v, %v}, want {%v, %v}", got.Latitude, got.Longitude, tt.lat, tt.lng)
				}
			}
		})
	}
}
