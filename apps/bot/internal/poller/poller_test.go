package poller

import (
	"strings"
	"testing"

	"github.com/marciojunior/waytutor/apps/bot/internal/api"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name     string
		accounts []api.Account
		want     []string
	}{
		{
			name:     "one account",
			accounts: []api.Account{{Email: "a@example.com", Name: "Ana"}},
			want:     []string{"Cadastro aguardando liberação:", "• a@example.com (Ana)"},
		},
		{
			name: "two accounts, one without a name",
			accounts: []api.Account{
				{Email: "a@example.com", Name: "Ana"},
				{Email: "b@example.com"},
			},
			want: []string{"2 cadastros aguardando liberação:", "• a@example.com (Ana)", "• b@example.com (sem nome)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := format(tt.accounts)

			for _, fragment := range tt.want {
				if !strings.Contains(got, fragment) {
					t.Errorf("format() = %q, missing %q", got, fragment)
				}
			}
			if !strings.Contains(got, "/tutor approve") {
				t.Errorf("format() = %q, missing the approve hint", got)
			}
		})
	}
}
