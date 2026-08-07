package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordTooShort = errors.New("password must be at least 8 characters")

type Password struct {
	hash string
}

func NewPassword(plaintext string) (Password, error) {
	if len(plaintext) < 8 {
		return Password{}, ErrPasswordTooShort
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}
	return Password{hash: string(hashed)}, nil
}

func (p Password) Hash() string {
	return p.hash
}

func (p Password) Matches(plaintext string) bool {
	return bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plaintext)) == nil
}
