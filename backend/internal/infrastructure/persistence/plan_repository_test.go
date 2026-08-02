package persistence_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	domain "github.com/Application-drop-up/Travellle/internal/domain/plan"
	"github.com/Application-drop-up/Travellle/internal/infrastructure/persistence"
)

func newTestPlan(t *testing.T) *domain.Plan {
	t.Helper()
	return &domain.Plan{
		ID:         uuid.New(),
		Title:      "Test Plan",
		ShareToken: uuid.New().String(),
	}
}

func TestPlanRepository_Create(t *testing.T) {
	t.Parallel()

	conn := newTestDB(t)
	repo := persistence.NewPlanRepository(conn)

	plan := newTestPlan(t)
	t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM plans WHERE id = $1", plan.ID) })

	if err := repo.Create(context.Background(), plan); err != nil {
		t.Fatalf("Create() unexpected error: %v", err)
	}

	if plan.CreatedAt.IsZero() {
		t.Error("Create() did not populate CreatedAt")
	}
	if plan.UpdatedAt.IsZero() {
		t.Error("Create() did not populate UpdatedAt")
	}
}

func TestPlanRepository_FindByShareToken(t *testing.T) {
	t.Parallel()

	conn := newTestDB(t)
	repo := persistence.NewPlanRepository(conn)

	t.Run("returns the plan for an existing share token", func(t *testing.T) {
		t.Parallel()

		plan := newTestPlan(t)
		t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM plans WHERE id = $1", plan.ID) })

		if err := repo.Create(context.Background(), plan); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}

		got, err := repo.FindByShareToken(context.Background(), plan.ShareToken)
		if err != nil {
			t.Fatalf("FindByShareToken() unexpected error: %v", err)
		}

		if got.ID != plan.ID {
			t.Errorf("FindByShareToken().ID = %v, want %v", got.ID, plan.ID)
		}
		if got.Title != plan.Title {
			t.Errorf("FindByShareToken().Title = %q, want %q", got.Title, plan.Title)
		}
		if got.ShareToken != plan.ShareToken {
			t.Errorf("FindByShareToken().ShareToken = %q, want %q", got.ShareToken, plan.ShareToken)
		}
	})

	t.Run("returns ErrNotFound for an unknown share token", func(t *testing.T) {
		t.Parallel()

		_, err := repo.FindByShareToken(context.Background(), uuid.New().String())
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("FindByShareToken() error = %v, want %v", err, domain.ErrNotFound)
		}
	})
}
