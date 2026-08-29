package persistence_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"

	plandomain "github.com/Application-drop-up/Travellle/internal/domain/plan"
	domain "github.com/Application-drop-up/Travellle/internal/domain/spot"
	"github.com/Application-drop-up/Travellle/internal/infrastructure/persistence"
	"github.com/Application-drop-up/Travellle/internal/testutil"
)

func newTestPlanForSpot(t *testing.T, conn *sql.DB) uuid.UUID {
	t.Helper()

	planRepo := persistence.NewPlanRepository(conn)
	p := &plandomain.Plan{
		ID:         uuid.New(),
		Title:      "Spot Test Plan",
		ShareToken: uuid.New().String(),
	}
	if err := planRepo.Create(context.Background(), p); err != nil {
		t.Fatalf("failed to create prerequisite plan: %v", err)
	}
	t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM plans WHERE id = $1", p.ID) })
	return p.ID
}

func newTestSpot(firstPlanID uuid.UUID) *domain.Spot {
	return &domain.Spot{
		ID:          uuid.New(),
		PlaceID:     domain.PlaceID(uuid.New().String()),
		Name:        "Tokyo Tower",
		Address:     "4 Chome-2-8 Shibakoen, Minato City, Tokyo",
		Location:    domain.Location{Latitude: 35.6586, Longitude: 139.7454},
		FirstPlanID: firstPlanID,
	}
}

func TestSpotRepository_Save(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	repo := persistence.NewSpotRepository(conn)

	t.Run("saves a spot", func(t *testing.T) {
		t.Parallel()

		planID := newTestPlanForSpot(t, conn)
		spot := newTestSpot(planID)
		t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM spots WHERE id = $1", spot.ID) })

		if err := repo.Save(context.Background(), spot); err != nil {
			t.Fatalf("Save() unexpected error: %v", err)
		}

		got, err := repo.FindByPlaceID(context.Background(), spot.PlaceID)
		if err != nil {
			t.Fatalf("FindByPlaceID() unexpected error: %v", err)
		}
		if got.Name != spot.Name || got.Address != spot.Address {
			t.Errorf("FindByPlaceID() = %+v, want fields matching %+v", got, spot)
		}
		if got.FirstPlanID != planID {
			t.Errorf("FindByPlaceID().FirstPlanID = %v, want %v", got.FirstPlanID, planID)
		}
	})

	t.Run("returns plan.ErrNotFound for a nonexistent plan", func(t *testing.T) {
		t.Parallel()

		spot := newTestSpot(uuid.New())

		err := repo.Save(context.Background(), spot)
		if !errors.Is(err, plandomain.ErrNotFound) {
			t.Errorf("Save() error = %v, want %v", err, plandomain.ErrNotFound)
		}
	})
}
