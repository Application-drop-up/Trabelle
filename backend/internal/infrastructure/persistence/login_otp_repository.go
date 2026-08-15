package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrLoginOTPNotFound = errors.New("login otp not found")

// LoginOTPRecord is this layer's own representation of a login_otps row,
// kept independent of the Application layer's OTP type.
type LoginOTPRecord struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Code      string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type LoginOTPRepository struct {
	db *sql.DB
}

func NewLoginOTPRepository(db *sql.DB) *LoginOTPRepository {
	return &LoginOTPRepository{db: db}
}

func (repo *LoginOTPRepository) Create(ctx context.Context, record *LoginOTPRecord) error {
	query := `
		INSERT INTO login_otps (id, user_id, code, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at`
	err := repo.db.QueryRowContext(ctx, query, record.ID, record.UserID, record.Code, record.ExpiresAt).
		Scan(&record.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert login otp: %w", err)
	}
	return nil
}

func (repo *LoginOTPRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*LoginOTPRecord, error) {
	query := `
		SELECT id, user_id, code, expires_at, created_at FROM login_otps
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1`
	record := &LoginOTPRecord{}
	err := repo.db.QueryRowContext(ctx, query, userID).
		Scan(&record.ID, &record.UserID, &record.Code, &record.ExpiresAt, &record.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrLoginOTPNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find login otp by user id: %w", err)
	}
	return record, nil
}

func (repo *LoginOTPRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM login_otps WHERE id = $1`
	result, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete login otp: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete login otp rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrLoginOTPNotFound
	}
	return nil
}
