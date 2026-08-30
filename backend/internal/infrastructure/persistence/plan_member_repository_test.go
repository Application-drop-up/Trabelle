package persistence_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"

	plandomain "github.com/Application-drop-up/Travellle/internal/domain/plan"
	domain "github.com/Application-drop-up/Travellle/internal/domain/planmember"
	userdomain "github.com/Application-drop-up/Travellle/internal/domain/user"
	"github.com/Application-drop-up/Travellle/internal/infrastructure/persistence"
	"github.com/Application-drop-up/Travellle/internal/testutil"
)

func newTestPlanForMember(t *testing.T, conn *sql.DB) uuid.UUID {
	t.Helper()

	planRepo := persistence.NewPlanRepository(conn)
	p := &plandomain.Plan{
		ID:         uuid.New(),
		Title:      "Plan Member Test Plan",
		ShareToken: uuid.New().String(),
	}
	if err := planRepo.Create(context.Background(), p); err != nil {
		t.Fatalf("failed to create prerequisite plan: %v", err)
	}
	t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM plans WHERE id = $1", p.ID) })
	return p.ID
}

func newTestUserForMember(t *testing.T, conn *sql.DB) uuid.UUID {
	t.Helper()

	userRepo := persistence.NewUserRepository(conn)
	u := &userdomain.User{
		ID:           uuid.New(),
		Email:        uuid.New().String() + "@example.com",
		PasswordHash: "hash",
		Name:         "Plan Member Test User",
	}
	if err := userRepo.Create(context.Background(), u); err != nil {
		t.Fatalf("failed to create prerequisite user: %v", err)
	}
	t.Cleanup(func() { _, _ = conn.Exec("DELETE FROM users WHERE id = $1", u.ID) })
	return u.ID
}

func TestPlanMemberRepository_Create(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	repo := persistence.NewPlanMemberRepository(conn)

	t.Run("adds a member", func(t *testing.T) {
		t.Parallel()

		planID := newTestPlanForMember(t, conn)
		userID := newTestUserForMember(t, conn)
		member := &domain.PlanMember{ID: uuid.New(), PlanID: planID, UserID: userID}

		if err := repo.Create(context.Background(), member); err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}
		if member.CreatedAt.IsZero() {
			t.Error("Create() did not populate CreatedAt")
		}
	})

	t.Run("returns ErrAlreadyMember for a duplicate membership", func(t *testing.T) {
		t.Parallel()

		planID := newTestPlanForMember(t, conn)
		userID := newTestUserForMember(t, conn)
		first := &domain.PlanMember{ID: uuid.New(), PlanID: planID, UserID: userID}
		if err := repo.Create(context.Background(), first); err != nil {
			t.Fatalf("first Create() unexpected error: %v", err)
		}

		second := &domain.PlanMember{ID: uuid.New(), PlanID: planID, UserID: userID}
		err := repo.Create(context.Background(), second)
		if !errors.Is(err, domain.ErrAlreadyMember) {
			t.Errorf("Create() error = %v, want %v", err, domain.ErrAlreadyMember)
		}
	})

	t.Run("returns plan.ErrNotFound for a nonexistent plan", func(t *testing.T) {
		t.Parallel()

		userID := newTestUserForMember(t, conn)
		member := &domain.PlanMember{ID: uuid.New(), PlanID: uuid.New(), UserID: userID}

		err := repo.Create(context.Background(), member)
		if !errors.Is(err, plandomain.ErrNotFound) {
			t.Errorf("Create() error = %v, want %v", err, plandomain.ErrNotFound)
		}
	})

	t.Run("returns user.ErrNotFound for a nonexistent user", func(t *testing.T) {
		t.Parallel()

		planID := newTestPlanForMember(t, conn)
		member := &domain.PlanMember{ID: uuid.New(), PlanID: planID, UserID: uuid.New()}

		err := repo.Create(context.Background(), member)
		if !errors.Is(err, userdomain.ErrNotFound) {
			t.Errorf("Create() error = %v, want %v", err, userdomain.ErrNotFound)
		}
	})
}

func TestPlanMemberRepository_FindByPlanID(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	repo := persistence.NewPlanMemberRepository(conn)

	planID := newTestPlanForMember(t, conn)
	userA := newTestUserForMember(t, conn)
	userB := newTestUserForMember(t, conn)

	for _, userID := range []uuid.UUID{userA, userB} {
		member := &domain.PlanMember{ID: uuid.New(), PlanID: planID, UserID: userID}
		if err := repo.Create(context.Background(), member); err != nil {
			t.Fatalf("failed to seed member: %v", err)
		}
	}

	members, err := repo.FindByPlanID(context.Background(), planID)
	if err != nil {
		t.Fatalf("FindByPlanID() unexpected error: %v", err)
	}
	if len(members) != 2 {
		t.Fatalf("FindByPlanID() returned %d members, want 2", len(members))
	}
}

func TestPlanMemberRepository_Delete(t *testing.T) {
	t.Parallel()

	conn := testutil.NewTestDB(t)
	repo := persistence.NewPlanMemberRepository(conn)

	t.Run("removes a member", func(t *testing.T) {
		t.Parallel()

		planID := newTestPlanForMember(t, conn)
		userID := newTestUserForMember(t, conn)
		member := &domain.PlanMember{ID: uuid.New(), PlanID: planID, UserID: userID}
		if err := repo.Create(context.Background(), member); err != nil {
			t.Fatalf("failed to seed member: %v", err)
		}

		if err := repo.Delete(context.Background(), planID, userID); err != nil {
			t.Fatalf("Delete() unexpected error: %v", err)
		}

		members, err := repo.FindByPlanID(context.Background(), planID)
		if err != nil {
			t.Fatalf("FindByPlanID() unexpected error: %v", err)
		}
		if len(members) != 0 {
			t.Errorf("FindByPlanID() after delete = %d members, want 0", len(members))
		}
	})

	t.Run("returns ErrNotFound when the membership doesn't exist", func(t *testing.T) {
		t.Parallel()

		planID := newTestPlanForMember(t, conn)
		userID := newTestUserForMember(t, conn)

		err := repo.Delete(context.Background(), planID, userID)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("Delete() error = %v, want %v", err, domain.ErrNotFound)
		}
	})
}
