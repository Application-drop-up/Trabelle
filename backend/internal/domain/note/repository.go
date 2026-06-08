package note

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, note *Note) error
	FindByID(ctx context.Context, pinID, noteID uuid.UUID) (*Note, error)
	FindByPinID(ctx context.Context, pinID uuid.UUID) ([]*Note, error)
	Update(ctx context.Context, note *Note) error
	Delete(ctx context.Context, pinID, noteID uuid.UUID) error
}
