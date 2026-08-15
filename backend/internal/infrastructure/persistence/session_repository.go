package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	useruc "github.com/Application-drop-up/Travellle/internal/usecase/user"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (repo *SessionRepository) Create(ctx context.Context, session *useruc.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, token, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at`
	err := repo.db.QueryRowContext(ctx, query, session.ID, session.UserID, session.Token, session.ExpiresAt).
		Scan(&session.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert session: %w", err)
	}
	return nil
}

func (repo *SessionRepository) FindByToken(ctx context.Context, token string) (*useruc.Session, error) {
	query := `SELECT id, user_id, token, expires_at, created_at FROM sessions WHERE token = $1`
	session := &useruc.Session{}
	err := repo.db.QueryRowContext(ctx, query, token).
		Scan(&session.ID, &session.UserID, &session.Token, &session.ExpiresAt, &session.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, useruc.ErrSessionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find session by token: %w", err)
	}
	return session, nil
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
		return useruc.ErrSessionNotFound
	}
	return nil
}
