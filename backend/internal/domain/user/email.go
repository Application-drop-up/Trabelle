package user

import (
	"errors"
	"net/mail"
)

var ErrInvalidEmail = errors.New("invalid email address")

type Email struct {
	value string
}

func NewEmail(raw string) (Email, error) {
	addr, err := mail.ParseAddress(raw)
	if err != nil {
		return Email{}, ErrInvalidEmail
	}
	return Email{value: addr.Address}, nil
}

func (e Email) String() string {
	return e.value
}
