package persistence_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"

	domain "github.com/Application-drop-up/Travellle/internal/domain/note"
	pindomain "github.com/Application-drop-up/Travellle/internal/domain/pin"
	plandomain "github.com/Application-drop-up/Travellle/internal/domain/plan"
	"github.com/Application-drop-up/Travellle/internal/infrastructure/persistence"
	"github.com/Application-drop-up/Travellle/internal/testutil"
)

func newTestPinForNote(t *testing.T, conn *sql.DB) uuid.UUID {
	t.Helper()

	planRepo := persistence.NewPlanRepository(conn)
	plan := &plandomain.Plan{
		ID:         uuid.New(),
		Title:      "Note Test Plan",
		ShareToken: uuid.New().String(),
	}
	if err := planRepo.Create(context.Background(), plan); err != nil {
		t.Fatalf("failed to create prerequisite plan: %v", err)
	}
	t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM plans WHERE id = $1", plan.ID) })

	pinRepo := persistence.NewPinRepository(conn)
	pin := &pindomain.Pin{
		ID:        uuid.New(),
		PlanID:    plan.ID,
		Name:      "Tokyo Tower",
		Latitude:  35.6586,
		Longitude: 139.7454,
		Category:  pindomain.CategorySightseeing,
		Colour:    "#FF5733",
	}
	if err := pinRepo.Create(context.Background(), pin); err != nil {
		t.Fatalf("failed to create prerequisite pin: %v", err)
	}
	return pin.ID
}

func newTestNote(pinID uuid.UUID) *domain.Note {
	return &domain.Note{
		ID:      uuid.New(),
		PinID:   pinID,
		Content: "Great view at sunset",
	}
}

func TestNoteRepository_Create(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	repo := persistence.NewNoteRepository(conn)
	pinID := newTestPinForNote(t, conn)

	t.Run("creates a note", func(t *testing.T) {
		t.Parallel()

		note := newTestNote(pinID)
		t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM notes WHERE id = $1", note.ID) })

		if err := repo.Create(context.Background(), note); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}
		if note.CreatedAt.IsZero() {
			t.Error("Create() did not populate CreatedAt")
		}
	})

	t.Run("returns pin.ErrNotFound for a nonexistent pin", func(t *testing.T) {
		t.Parallel()

		note := newTestNote(uuid.New())

		err := repo.Create(context.Background(), note)
		if !errors.Is(err, pindomain.ErrNotFound) {
			t.Errorf("Create() error = %v, want %v", err, pindomain.ErrNotFound)
		}
	})
}

func TestNoteRepository_FindByID(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	repo := persistence.NewNoteRepository(conn)
	pinID := newTestPinForNote(t, conn)

	t.Run("returns the note when it exists", func(t *testing.T) {
		t.Parallel()

		note := newTestNote(pinID)
		t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM notes WHERE id = $1", note.ID) })
		if err := repo.Create(context.Background(), note); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}

		got, err := repo.FindByID(context.Background(), pinID, note.ID)
		if err != nil {
			t.Fatalf("FindByID() unexpected error: %v", err)
		}
		if got.Content != note.Content {
			t.Errorf("FindByID().Content = %q, want %q", got.Content, note.Content)
		}
	})

	t.Run("returns ErrNotFound for an unknown note", func(t *testing.T) {
		t.Parallel()

		_, err := repo.FindByID(context.Background(), pinID, uuid.New())
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("FindByID() error = %v, want %v", err, domain.ErrNotFound)
		}
	})
}

func TestNoteRepository_FindByPinID(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	repo := persistence.NewNoteRepository(conn)
	pinID := newTestPinForNote(t, conn)

	noteA := newTestNote(pinID)
	noteB := newTestNote(pinID)
	for _, n := range []*domain.Note{noteA, noteB} {
		t.Cleanup(func(id uuid.UUID) func() {
			return func() { _, _ = conn.Exec("DELETE FROM notes WHERE id = $1", id) }
		}(n.ID))
		if err := repo.Create(context.Background(), n); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}
	}

	got, err := repo.FindByPinID(context.Background(), pinID)
	if err != nil {
		t.Fatalf("FindByPinID() unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("FindByPinID() returned %d notes, want 2", len(got))
	}
}

func TestNoteRepository_Update(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	repo := persistence.NewNoteRepository(conn)
	pinID := newTestPinForNote(t, conn)

	t.Run("updates content", func(t *testing.T) {
		t.Parallel()

		note := newTestNote(pinID)
		t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM notes WHERE id = $1", note.ID) })
		if err := repo.Create(context.Background(), note); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}

		note.Content = "Updated content"
		if err := repo.Update(context.Background(), note); err != nil {
			t.Fatalf("Update() unexpected error: %v", err)
		}

		got, err := repo.FindByID(context.Background(), pinID, note.ID)
		if err != nil {
			t.Fatalf("FindByID() unexpected error: %v", err)
		}
		if got.Content != "Updated content" {
			t.Errorf("after Update(), Content = %q, want %q", got.Content, "Updated content")
		}
	})

	t.Run("returns ErrNotFound for an unknown note", func(t *testing.T) {
		t.Parallel()

		note := newTestNote(pinID)
		note.ID = uuid.New()

		err := repo.Update(context.Background(), note)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("Update() error = %v, want %v", err, domain.ErrNotFound)
		}
	})
}

func TestNoteRepository_Delete(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	repo := persistence.NewNoteRepository(conn)
	pinID := newTestPinForNote(t, conn)

	t.Run("deletes an existing note", func(t *testing.T) {
		t.Parallel()

		note := newTestNote(pinID)
		if err := repo.Create(context.Background(), note); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}

		if err := repo.Delete(context.Background(), pinID, note.ID); err != nil {
			t.Fatalf("Delete() unexpected error: %v", err)
		}

		_, err := repo.FindByID(context.Background(), pinID, note.ID)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("FindByID() after Delete() error = %v, want %v", err, domain.ErrNotFound)
		}
	})

	t.Run("returns ErrNotFound for an unknown note", func(t *testing.T) {
		t.Parallel()

		err := repo.Delete(context.Background(), pinID, uuid.New())
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("Delete() error = %v, want %v", err, domain.ErrNotFound)
		}
	})
}
