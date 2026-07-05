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

func (repo *NoteRepository) Create(ctx context.Context, note *domain.Note) error {
	query := `
		INSERT INTO notes (id, pin_id, content)
		VALUES ($1, $2, $3)
		RETURNING created_at, updated_at`
	err := repo.db.QueryRowContext(ctx, query, note.ID, note.PinID, note.Content).
		Scan(&note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgFKViolation {
			return pindomain.ErrNotFound
		}
		return fmt.Errorf("insert note: %w", err)
	}
	return nil
}

func (repo *NoteRepository) FindByID(ctx context.Context, pinID, noteID uuid.UUID) (*domain.Note, error) {
	query := `
		SELECT id, pin_id, content, created_at, updated_at
		FROM notes WHERE id = $1 AND pin_id = $2`
	note := &domain.Note{}
	err := repo.db.QueryRowContext(ctx, query, noteID, pinID).
		Scan(&note.ID, &note.PinID, &note.Content, &note.CreatedAt, &note.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find note by id: %w", err)
	}
	return note, nil
}

func (repo *NoteRepository) FindByPinID(ctx context.Context, pinID uuid.UUID) ([]*domain.Note, error) {
	query := `
		SELECT id, pin_id, content, created_at, updated_at
		FROM notes WHERE pin_id = $1 ORDER BY created_at ASC`
	rows, err := repo.db.QueryContext(ctx, query, pinID)
	if err != nil {
		return nil, fmt.Errorf("find notes by pin id: %w", err)
	}
	defer rows.Close()

	var notes []*domain.Note
	for rows.Next() {
		note := &domain.Note{}
		if err := rows.Scan(&note.ID, &note.PinID, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan note: %w", err)
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("find notes by pin id rows: %w", err)
	}
	return notes, nil
}

func (repo *NoteRepository) Update(ctx context.Context, note *domain.Note) error {
	query := `
		UPDATE notes SET content = $1, updated_at = NOW()
		WHERE id = $2 AND pin_id = $3
		RETURNING updated_at`
	err := repo.db.QueryRowContext(ctx, query, note.Content, note.ID, note.PinID).Scan(&note.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("update note: %w", err)
	}
	return nil
}

func (repo *NoteRepository) Delete(ctx context.Context, pinID, noteID uuid.UUID) error {
	query := `DELETE FROM notes WHERE id = $1 AND pin_id = $2`
	result, err := repo.db.ExecContext(ctx, query, noteID, pinID)
	if err != nil {
		return fmt.Errorf("delete note: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete note rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
