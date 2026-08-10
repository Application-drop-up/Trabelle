package user_test

import (
	"errors"
	"testing"

	"github.com/Application-drop-up/Travellle/internal/domain/user"
)

func TestNewEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr error
	}{
		{
			name: "valid email",
			raw:  "taro@example.com",
			want: "taro@example.com",
		},
		{
			name:    "empty string",
			raw:     "",
			wantErr: user.ErrInvalidEmail,
		},
		{
			name:    "missing @",
			raw:     "taro.example.com",
			wantErr: user.ErrInvalidEmail,
		},
		{
			name:    "missing domain",
			raw:     "taro@",
			wantErr: user.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := user.NewEmail(tt.raw)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewEmail(%q) error = %v, want %v", tt.raw, err, tt.wantErr)
				return
			}
			if tt.wantErr == nil && got.String() != tt.want {
				t.Errorf("Email.String() = %q, want %q", got.String(), tt.want)
			}
		})
	}
}
