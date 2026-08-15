package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrSessionNotFound = errors.New("session not found")

// SessionRecord is this layer's own representation of a session row,
// kept independent of the Application layer's Session type.
type SessionRecord struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (repo *SessionRepository) Create(ctx context.Context, record *SessionRecord) error {
	query := `
		INSERT INTO sessions (id, user_id, token, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at`
	err := repo.db.QueryRowContext(ctx, query, record.ID, record.UserID, record.Token, record.ExpiresAt).
		Scan(&record.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert session: %w", err)
	}
	return nil
}

func (repo *SessionRepository) FindByToken(ctx context.Context, token string) (*SessionRecord, error) {
	query := `SELECT id, user_id, token, expires_at, created_at FROM sessions WHERE token = $1`
	record := &SessionRecord{}
	err := repo.db.QueryRowContext(ctx, query, token).
		Scan(&record.ID, &record.UserID, &record.Token, &record.ExpiresAt, &record.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrSessionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find session by token: %w", err)
	}
	return record, nil
}

func (repo *SessionRepository) Delete(ctx context.Context, token string) error {
	query := `DELETE FROM sessions WHERE token = $1`
	result, err := repo.db.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete session rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrSessionNotFound
	}
	return nil
}
