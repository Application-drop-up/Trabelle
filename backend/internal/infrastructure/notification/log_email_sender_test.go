package notification_test

import (
	"context"
	"testing"

	"github.com/Application-drop-up/Travellle/internal/infrastructure/notification"
)

func TestLogEmailSender_SendLoginCode(t *testing.T) {
	t.Parallel()

	sender := notification.NewLogEmailSender()

	if err := sender.SendLoginCode(context.Background(), "taro@example.com", "123456"); err != nil {
		t.Fatalf("SendLoginCode() unexpected error: %v", err)
	}
}
