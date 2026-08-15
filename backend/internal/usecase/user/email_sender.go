package user

import "context"

type EmailSender interface {
	SendLoginCode(ctx context.Context, to, code string) error
}
