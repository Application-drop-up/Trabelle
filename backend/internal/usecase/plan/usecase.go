package plan

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	domain "github.com/Application-drop-up/Travellle/internal/domain/plan"
	"github.com/google/uuid"
)

type UseCase struct {
	repo domain.Repository
}

func New(repo domain.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) CreatePlan(ctx context.Context, title string) (*domain.Plan, error) {
	token, err := generateShareToken()
	if err != nil {
		return nil, fmt.Errorf("generate share token: %w", err)
	}

	p := &domain.Plan{
		ID:         uuid.New(),
		Title:      title,
		ShareToken: token,
	}

	if err := uc.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("create plan: %w", err)
	}
	return p, nil
}

func (uc *UseCase) GetPlanByShareToken(ctx context.Context, token string) (*domain.Plan, error) {
	return uc.repo.FindByShareToken(ctx, token)
}

func generateShareToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
