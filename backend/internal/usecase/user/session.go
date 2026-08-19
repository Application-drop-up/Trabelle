package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrSessionNotFound = errors.New("session not found")

const SessionTTL = 24 * time.Hour

// Session is the technical mechanism used to realize the login use case,
// not a Domain concept -- it has no business invariants of its own (see #326).
type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	FindByToken(ctx context.Context, token string) (*Session, error)
	Delete(ctx context.Context, token string) error
}

func (session *Session) isExpired() bool {
	return time.Now().After(session.ExpiresAt)
}

func newSession(userID uuid.UUID) (*Session, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}
	return &Session{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(SessionTTL),
	}, nil
}

func generateToken() (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(tokenBytes), nil
}
