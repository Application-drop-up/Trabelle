package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	pindomain "github.com/Application-drop-up/Travellle/internal/domain/pin"
	domain "github.com/Application-drop-up/Travellle/internal/domain/note"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type NoteRepository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) Create(ctx context.Context, n *domain.Note) error {
	query := `
		INSERT INTO notes (id, pin_id, content)
		VALUES ($1, $2, $3)
		RETURNING created_at, updated_at`
	err := r.db.QueryRowContext(ctx, query, n.ID, n.PinID, n.Content).
		Scan(&n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgFKViolation {
			return pindomain.ErrNotFound
		}
		return fmt.Errorf("insert note: %w", err)
	}
	return nil
}

func (r *NoteRepository) FindByID(ctx context.Context, pinID, noteID uuid.UUID) (*domain.Note, error) {
	query := `
		SELECT id, pin_id, content, created_at, updated_at
		FROM notes WHERE id = $1 AND pin_id = $2`
	n := &domain.Note{}
	err := r.db.QueryRowContext(ctx, query, noteID, pinID).
		Scan(&n.ID, &n.PinID, &n.Content, &n.CreatedAt, &n.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find note by id: %w", err)
	}
	return n, nil
}

func (r *NoteRepository) FindByPinID(ctx context.Context, pinID uuid.UUID) ([]*domain.Note, error) {
	query := `
		SELECT id, pin_id, content, created_at, updated_at
		FROM notes WHERE pin_id = $1 ORDER BY created_at ASC`
	rows, err := r.db.QueryContext(ctx, query, pinID)
	if err != nil {
		return nil, fmt.Errorf("find notes by pin id: %w", err)
	}
	defer rows.Close()

	var notes []*domain.Note
	for rows.Next() {
		n := &domain.Note{}
		if err := rows.Scan(&n.ID, &n.PinID, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan note: %w", err)
		}
		notes = append(notes, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("find notes by pin id rows: %w", err)
	}
	return notes, nil
}

func (r *NoteRepository) Update(ctx context.Context, n *domain.Note) error {
	query := `
		UPDATE notes SET content = $1, updated_at = NOW()
		WHERE id = $2 AND pin_id = $3
		RETURNING updated_at`
	err := r.db.QueryRowContext(ctx, query, n.Content, n.ID, n.PinID).Scan(&n.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("update note: %w", err)
	}
	return nil
}

func (r *NoteRepository) Delete(ctx context.Context, pinID, noteID uuid.UUID) error {
	query := `DELETE FROM notes WHERE id = $1 AND pin_id = $2`
	result, err := r.db.ExecContext(ctx, query, noteID, pinID)
	if err != nil {
		return fmt.Errorf("delete note: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete note rows affected: %w", err)
	}
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
