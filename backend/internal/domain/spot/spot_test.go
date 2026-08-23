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
		ID:                spotID,
		PlaceID:           placeID,
		Name:              "Tokyo Tower",
		Address:           "4 Chome-2-8 Shibakoen, Minato City, Tokyo",
		Location:          location,
		FirstPlanID:       firstPlanID,
		FirstPlanIsPublic: true,
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
	if !got.FirstPlanIsPublic {
		t.Error("Spot.FirstPlanIsPublic = false, want true")
	}
}

func TestSpot_IsAttributionVisible(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		firstPlanIsPublic bool
		want              bool
	}{
		{name: "visible when first plan is public", firstPlanIsPublic: true, want: true},
		{name: "not visible when first plan is private", firstPlanIsPublic: false, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testSpot := &spot.Spot{FirstPlanIsPublic: tt.firstPlanIsPublic}
			if got := testSpot.IsAttributionVisible(); got != tt.want {
				t.Errorf("IsAttributionVisible() = %v, want %v", got, tt.want)
			}
		})
	}
}
