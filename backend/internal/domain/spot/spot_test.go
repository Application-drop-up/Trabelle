package spot_test

import (
	"testing"

	"github.com/Application-drop-up/Travellle/internal/domain/spot"
)

func TestSpot_Fields(t *testing.T) {
	t.Parallel()

	placeID, err := spot.NewPlaceID("ChIJ5eTFBkqLGGARsV4PF3rDVAA")
	if err != nil {
		t.Fatalf("NewPlaceID() unexpected error: %v", err)
	}

	location, err := spot.NewLocation(35.6895, 139.6917)
	if err != nil {
		t.Fatalf("NewLocation() unexpected error: %v", err)
	}

	s := spot.Spot{
		PlaceID:  placeID,
		Name:     "Tokyo Tower",
		Address:  "4 Chome-2-8 Shibakoen, Minato City, Tokyo",
		Location: location,
	}

	if s.PlaceID != placeID {
		t.Errorf("Spot.PlaceID = %v, want %v", s.PlaceID, placeID)
	}
	if s.Name != "Tokyo Tower" {
		t.Errorf("Spot.Name = %q, want %q", s.Name, "Tokyo Tower")
	}
	if s.Location.Latitude != 35.6895 || s.Location.Longitude != 139.6917 {
		t.Errorf("Spot.Location = %+v, want {35.6895, 139.6917}", s.Location)
	}
}
