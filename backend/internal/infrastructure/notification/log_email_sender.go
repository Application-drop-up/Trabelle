package notification

import (
	"context"
	"log"
)

// LogEmailSender logs the login code instead of sending a real email.
// Swap for a real provider (e.g. SendGrid, SES) when one is available.
type LogEmailSender struct{}

func NewLogEmailSender() *LogEmailSender {
	return &LogEmailSender{}
}

func (s *LogEmailSender) SendLoginCode(_ context.Context, to, code string) error {
	log.Printf("[email] login code for %s: %s", to, code)
	return nil
}
