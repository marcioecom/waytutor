package discord

import (
	"strings"
	"testing"

	"github.com/marciojunior/waytutor/apps/bot/internal/api"
)

func TestEmailChoices(t *testing.T) {
	accounts := []api.Account{
		{Email: "ana@example.com", Name: "Ana"},
		{Email: "bruno@example.com"},
		{Email: "carla@outra.com", Name: "Carla"},
	}

	tests := []struct {
		name  string
		query string
		want  []string
	}{
		{"no query lists everything", "", []string{"ana@example.com (Ana)", "bruno@example.com", "carla@outra.com (Carla)"}},
		{"filters by local part", "br", []string{"bruno@example.com"}},
		{"filters by domain", "outra", []string{"carla@outra.com (Carla)"}},
		{"ignores case", "ANA@", []string{"ana@example.com (Ana)"}},
		{"no match", "zzz", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			choices := emailChoices(accounts, tt.query)
			if len(choices) != len(tt.want) {
				t.Fatalf("emailChoices() returned %d choices, want %d", len(choices), len(tt.want))
			}

			for index, choice := range choices {
				if choice.Name != tt.want[index] {
					t.Errorf("choice %d name = %q, want %q", index, choice.Name, tt.want[index])
				}
			}
		})
	}
}

func TestEmailChoicesValueIsTheEmailAndRespectsTheLimit(t *testing.T) {
	accounts := make([]api.Account, 0, 40)
	for index := range 40 {
		accounts = append(accounts, api.Account{Email: strings.Repeat("a", index+1) + "@example.com"})
	}

	choices := emailChoices(accounts, "")
	if len(choices) != maxChoices {
		t.Fatalf("emailChoices() returned %d choices, want %d", len(choices), maxChoices)
	}
	for _, choice := range choices {
		if choice.Name != choice.Value {
			t.Errorf("name %q differs from value %q for an account without a name", choice.Name, choice.Value)
		}
	}
}

func TestAutocompleteStatus(t *testing.T) {
	tests := map[string]string{
		"approve": "pending",
		"pause":   "active",
		"revoke":  "",
	}
	for subcommand, want := range tests {
		if got := autocompleteStatus(subcommand); got != want {
			t.Errorf("autocompleteStatus(%q) = %q, want %q", subcommand, got, want)
		}
	}
}

func TestTruncateKeepsRunesWhole(t *testing.T) {
	if got := truncate("ação", 2); got != "aç" {
		t.Errorf("truncate() = %q, want aç", got)
	}
	if got := truncate("abc", 10); got != "abc" {
		t.Errorf("truncate() = %q, want abc", got)
	}
}
