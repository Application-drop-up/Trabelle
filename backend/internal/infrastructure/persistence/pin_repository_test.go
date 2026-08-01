package persistence_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"

	domain "github.com/Application-drop-up/Travellle/internal/domain/pin"
	plandomain "github.com/Application-drop-up/Travellle/internal/domain/plan"
	"github.com/Application-drop-up/Travellle/internal/infrastructure/persistence"
	"github.com/Application-drop-up/Travellle/internal/testutil"
)

func newTestPlanForPin(t *testing.T, conn *sql.DB) uuid.UUID {
	t.Helper()

	planRepo := persistence.NewPlanRepository(conn)
	p := &plandomain.Plan{
		ID:         uuid.New(),
		Title:      "Pin Test Plan",
		ShareToken: uuid.New().String(),
	}
	if err := planRepo.Create(context.Background(), p); err != nil {
		t.Fatalf("failed to create prerequisite plan: %v", err)
	}
	t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM plans WHERE id = $1", p.ID) })
	return p.ID
}

func newTestPin(planID uuid.UUID) *domain.Pin {
	return &domain.Pin{
		ID:        uuid.New(),
		PlanID:    planID,
		Name:      "Tokyo Tower",
		Latitude:  35.6586,
		Longitude: 139.7454,
		Category:  domain.CategorySightseeing,
		Colour:    "#FF5733",
	}
}

func TestPinRepository_Create(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	repo := persistence.NewPinRepository(conn)
	planID := newTestPlanForPin(t, conn)

	t.Run("creates a pin", func(t *testing.T) {
		t.Parallel()

		pin := newTestPin(planID)
		t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM pins WHERE id = $1", pin.ID) })

		if err := repo.Create(context.Background(), pin); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}
		if pin.CreatedAt.IsZero() {
			t.Error("Create() did not populate CreatedAt")
		}
	})

	t.Run("returns plan.ErrNotFound for a nonexistent plan", func(t *testing.T) {
		t.Parallel()

		pin := newTestPin(uuid.New())

		err := repo.Create(context.Background(), pin)
		if !errors.Is(err, plandomain.ErrNotFound) {
			t.Errorf("Create() error = %v, want %v", err, plandomain.ErrNotFound)
		}
	})
}

func TestPinRepository_FindByID(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	repo := persistence.NewPinRepository(conn)
	planID := newTestPlanForPin(t, conn)

	t.Run("returns the pin when it exists", func(t *testing.T) {
		t.Parallel()

		pin := newTestPin(planID)
		t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM pins WHERE id = $1", pin.ID) })
		if err := repo.Create(context.Background(), pin); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}

		got, err := repo.FindByID(context.Background(), planID, pin.ID)
		if err != nil {
			t.Fatalf("FindByID() unexpected error: %v", err)
		}
		if got.Name != pin.Name || got.Category != pin.Category || got.Colour != pin.Colour {
			t.Errorf("FindByID() = %+v, want fields matching %+v", got, pin)
		}
	})

	t.Run("returns ErrNotFound for an unknown pin", func(t *testing.T) {
		t.Parallel()

		_, err := repo.FindByID(context.Background(), planID, uuid.New())
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("FindByID() error = %v, want %v", err, domain.ErrNotFound)
		}
	})
}

func TestPinRepository_FindByPlanID(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	repo := persistence.NewPinRepository(conn)
	planID := newTestPlanForPin(t, conn)

	pinA := newTestPin(planID)
	pinB := newTestPin(planID)
	for _, p := range []*domain.Pin{pinA, pinB} {
		t.Cleanup(func(id uuid.UUID) func() {
			return func() { _, _ = conn.Exec("DELETE FROM pins WHERE id = $1", id) }
		}(p.ID))
		if err := repo.Create(context.Background(), p); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}
	}

	got, err := repo.FindByPlanID(context.Background(), planID)
	if err != nil {
		t.Fatalf("FindByPlanID() unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("FindByPlanID() returned %d pins, want 2", len(got))
	}
}

func TestPinRepository_Update(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	repo := persistence.NewPinRepository(conn)
	planID := newTestPlanForPin(t, conn)

	t.Run("updates category and colour", func(t *testing.T) {
		t.Parallel()

		pin := newTestPin(planID)
		t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM pins WHERE id = $1", pin.ID) })
		if err := repo.Create(context.Background(), pin); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}

		pin.Category = domain.CategoryRestaurant
		pin.Colour = "#00FF00"
		if err := repo.Update(context.Background(), pin); err != nil {
			t.Fatalf("Update() unexpected error: %v", err)
		}

		got, err := repo.FindByID(context.Background(), planID, pin.ID)
		if err != nil {
			t.Fatalf("FindByID() unexpected error: %v", err)
		}
		if got.Category != domain.CategoryRestaurant || got.Colour != "#00FF00" {
			t.Errorf("after Update(), FindByID() = %+v, want Category=restaurant Colour=#00FF00", got)
		}
	})

	t.Run("returns ErrNotFound for an unknown pin", func(t *testing.T) {
		t.Parallel()

		pin := newTestPin(planID)
		pin.ID = uuid.New()

		err := repo.Update(context.Background(), pin)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("Update() error = %v, want %v", err, domain.ErrNotFound)
		}
	})
}

func TestPinRepository_Delete(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	repo := persistence.NewPinRepository(conn)
	planID := newTestPlanForPin(t, conn)

	t.Run("deletes an existing pin", func(t *testing.T) {
		t.Parallel()

		pin := newTestPin(planID)
		if err := repo.Create(context.Background(), pin); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}

		if err := repo.Delete(context.Background(), planID, pin.ID); err != nil {
			t.Fatalf("Delete() unexpected error: %v", err)
		}

		_, err := repo.FindByID(context.Background(), planID, pin.ID)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("FindByID() after Delete() error = %v, want %v", err, domain.ErrNotFound)
		}
	})

	t.Run("returns ErrNotFound for an unknown pin", func(t *testing.T) {
		t.Parallel()

		err := repo.Delete(context.Background(), planID, uuid.New())
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("Delete() error = %v, want %v", err, domain.ErrNotFound)
		}
	})
}
