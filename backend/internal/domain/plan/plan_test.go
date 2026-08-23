package plan_test

import (
	"testing"

	"github.com/Application-drop-up/Travellle/internal/domain/plan"
)

func TestPlan_IsViewableWithoutToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		isPublic bool
		want     bool
	}{
		{name: "public plan is viewable without token", isPublic: true, want: true},
		{name: "private plan is not viewable without token", isPublic: false, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testPlan := &plan.Plan{IsPublic: tt.isPublic}
			if got := testPlan.IsViewableWithoutToken(); got != tt.want {
				t.Errorf("IsViewableWithoutToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
