package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	useruc "github.com/Application-drop-up/Travellle/internal/usecase/user"
	"github.com/google/uuid"
)

type LoginOTPRepository struct {
	db *sql.DB
}

func NewLoginOTPRepository(db *sql.DB) *LoginOTPRepository {
	return &LoginOTPRepository{db: db}
}

func (repo *LoginOTPRepository) Create(ctx context.Context, otp *useruc.OTP) error {
	query := `
		INSERT INTO login_otps (id, user_id, code, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at`
	err := repo.db.QueryRowContext(ctx, query, otp.ID, otp.UserID, otp.Code, otp.ExpiresAt).
		Scan(&otp.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert login otp: %w", err)
	}
	return nil
}

func (repo *LoginOTPRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*useruc.OTP, error) {
	query := `
		SELECT id, user_id, code, expires_at, created_at FROM login_otps
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1`
	otp := &useruc.OTP{}
	err := repo.db.QueryRowContext(ctx, query, userID).
		Scan(&otp.ID, &otp.UserID, &otp.Code, &otp.ExpiresAt, &otp.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, useruc.ErrOTPNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find login otp by user id: %w", err)
	}
	return otp, nil
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
		return useruc.ErrOTPNotFound
	}
	return nil
}
