package user

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
)

var (
	ErrOTPNotFound = errors.New("otp not found")
	ErrOTPExpired  = errors.New("otp expired")
	ErrInvalidOTP  = errors.New("invalid code")
)

const (
	OTPTTL     = 10 * time.Minute
	otpCodeMax = 1000000
	otpCodeFmt = "%06d"
)

// OTP is the technical mechanism used to realize the login use case, not a
// Domain concept -- it has no business invariants of its own (see #326).
type OTP struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Code      string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type LoginOTPRepository interface {
	Create(ctx context.Context, otp *OTP) error
	FindByUserID(ctx context.Context, userID uuid.UUID) (*OTP, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

func newOTP(userID uuid.UUID) (*OTP, error) {
	code, err := generateCode()
	if err != nil {
		return nil, err
	}
	return &OTP{
		ID:        uuid.New(),
		UserID:    userID,
		Code:      code,
		ExpiresAt: time.Now().Add(OTPTTL),
	}, nil
}

func (otp *OTP) isExpired() bool {
	return time.Now().After(otp.ExpiresAt)
}

func (otp *OTP) matches(code string) bool {
	return otp.Code == code
}

func generateCode() (string, error) {
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(otpCodeMax))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(otpCodeFmt, randomNumber.Int64()), nil
}
