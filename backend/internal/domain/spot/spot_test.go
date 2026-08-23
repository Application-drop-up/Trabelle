package spot_test

import (
	"testing"

	"github.com/Application-drop-up/Travellle/internal/domain/spot"
	"github.com/google/uuid"
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

	spotID := uuid.New()
	firstPlanID := uuid.New()

	got := spot.Spot{
		ID:          spotID,
		PlaceID:     placeID,
		Name:        "Tokyo Tower",
		Address:     "4 Chome-2-8 Shibakoen, Minato City, Tokyo",
		Location:    location,
		FirstPlanID: firstPlanID,
	}

	if got.ID != spotID {
		t.Errorf("Spot.ID = %v, want %v", got.ID, spotID)
	}
	if got.PlaceID != placeID {
		t.Errorf("Spot.PlaceID = %v, want %v", got.PlaceID, placeID)
	}
	if got.Name != "Tokyo Tower" {
		t.Errorf("Spot.Name = %q, want %q", got.Name, "Tokyo Tower")
	}
	if got.Location.Latitude != 35.6895 || got.Location.Longitude != 139.6917 {
		t.Errorf("Spot.Location = %+v, want {35.6895, 139.6917}", got.Location)
	}
	if got.FirstPlanID != firstPlanID {
		t.Errorf("Spot.FirstPlanID = %v, want %v", got.FirstPlanID, firstPlanID)
	}
}
