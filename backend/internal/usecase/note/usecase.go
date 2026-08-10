package note

import (
	"context"
	"fmt"

	domain "github.com/Application-drop-up/Travellle/internal/domain/note"
	"github.com/google/uuid"
)

type UseCase struct {
	repo domain.Repository
}

func New(repo domain.Repository) *UseCase {
	return &UseCase{repo: repo}
}

type CreateInput struct {
	PinID   uuid.UUID
	Content string
}

type UpdateInput struct {
	Content string
}

func (useCase *UseCase) CreateNote(ctx context.Context, input CreateInput) (*domain.Note, error) {
	note := &domain.Note{
		ID:      uuid.New(),
		PinID:   input.PinID,
		Content: input.Content,
	}
	if err := useCase.repo.Create(ctx, note); err != nil {
		return nil, fmt.Errorf("create note: %w", err)
	}
	return note, nil
}

func (useCase *UseCase) UpdateNote(ctx context.Context, pinID, noteID uuid.UUID, input UpdateInput) (*domain.Note, error) {
	note, err := useCase.repo.FindByID(ctx, pinID, noteID)
	if err != nil {
		return nil, err
	}
	note.Content = input.Content
	if err := useCase.repo.Update(ctx, note); err != nil {
		return nil, fmt.Errorf("update note: %w", err)
	}
	return note, nil
}

func (useCase *UseCase) DeleteNote(ctx context.Context, pinID, noteID uuid.UUID) error {
	return useCase.repo.Delete(ctx, pinID, noteID)
}

func (useCase *UseCase) ListNotes(ctx context.Context, pinID uuid.UUID) ([]*domain.Note, error) {
	return useCase.repo.FindByPinID(ctx, pinID)
}
