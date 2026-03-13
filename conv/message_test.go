package conv

import (
	"testing"

	"github.com/ATTron/jitzu/config"
)

func TestBuild(t *testing.T) {
	tests := []struct {
		name                                      string
		typ, scope, subject, body, breaking, refs string
		want                                      string
	}{
		{
			name:    "simple",
			typ:     "feat",
			subject: "add login",
			want:    "feat: add login",
		},
		{
			name:    "with scope",
			typ:     "fix",
			scope:   "auth",
			subject: "handle expired tokens",
			want:    "fix(auth): handle expired tokens",
		},
		{
			name:    "with body",
			typ:     "feat",
			subject: "add login",
			body:    "This adds the login flow.",
			want:    "feat: add login\n\nThis adds the login flow.",
		},
		{
			name:     "with breaking change",
			typ:      "feat",
			scope:    "api",
			subject:  "change response format",
			breaking: "response is now JSON array",
			want:     "feat(api): change response format\n\nBREAKING CHANGE: response is now JSON array",
		},
		{
			name:    "with refs",
			typ:     "fix",
			subject: "null pointer",
			refs:    "#123",
			want:    "fix: null pointer\n\nRefs: #123",
		},
		{
			name:     "full message",
			typ:      "feat",
			scope:    "ui",
			subject:  "redesign dashboard",
			body:     "Complete redesign of the main dashboard.",
			breaking: "removed legacy widgets",
			refs:     "#456, #789",
			want:     "feat(ui): redesign dashboard\n\nComplete redesign of the main dashboard.\n\nBREAKING CHANGE: removed legacy widgets\n\nRefs: #456, #789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Build(tt.typ, tt.scope, tt.subject, tt.body, tt.breaking, tt.refs)
			if got != tt.want {
				t.Errorf("Build() =\n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Message
		wantErr bool
	}{
		{
			name:  "simple",
			input: "feat: add login",
			want:  Message{Type: "feat", Subject: "add login"},
		},
		{
			name:  "with scope",
			input: "fix(auth): handle expired tokens",
			want:  Message{Type: "fix", Scope: "auth", Subject: "handle expired tokens"},
		},
		{
			name:  "with body",
			input: "feat: add login\n\nThis adds the login flow.",
			want:  Message{Type: "feat", Subject: "add login", Body: "This adds the login flow."},
		},
		{
			name:  "with breaking indicator",
			input: "feat(api)!: change response format",
			want:  Message{Type: "feat", Scope: "api", Subject: "change response format", Breaking: "change response format"},
		},
		{
			name:  "with breaking footer",
			input: "feat: change format\n\nBREAKING CHANGE: new format",
			want:  Message{Type: "feat", Subject: "change format", Breaking: "new format"},
		},
		{
			name:  "with refs",
			input: "fix: bug\n\nRefs: #123",
			want:  Message{Type: "fix", Subject: "bug", Refs: "#123"},
		},
		{
			name:  "full message",
			input: "feat(ui): redesign\n\nComplete redesign.\n\nBREAKING CHANGE: removed widgets\n\nRefs: #456",
			want:  Message{Type: "feat", Scope: "ui", Subject: "redesign", Body: "Complete redesign.", Breaking: "removed widgets", Refs: "#456"},
		},
		{
			name:    "invalid header",
			input:   "not a conventional commit",
			wantErr: true,
		},
		{
			name:    "empty",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Parse() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	cfg := config.Default()

	tests := []struct {
		name     string
		msg      string
		cfg      config.Config
		problems int
	}{
		{
			name: "valid message",
			msg:  "feat: add login",
			cfg:  cfg,
		},
		{
			name:     "invalid type",
			msg:      "yolo: whatever",
			cfg:      cfg,
			problems: 1,
		},
		{
			name:     "subject too long",
			msg:      "feat: " + string(make([]byte, 100)),
			cfg:      config.Config{Types: cfg.Types, SubjectMaxLen: 72},
			problems: 1,
		},
		{
			name:     "scope required but missing",
			msg:      "feat: add login",
			cfg:      config.Config{Types: cfg.Types, ScopeRequired: true},
			problems: 1,
		},
		{
			name:     "scope not in allowed list",
			msg:      "feat(unknown): add login",
			cfg:      config.Config{Types: cfg.Types, Scopes: []string{"auth", "ui"}},
			problems: 1,
		},
		{
			name:     "body required but missing",
			msg:      "feat: add login",
			cfg:      config.Config{Types: cfg.Types, BodyRequired: true},
			problems: 1,
		},
		{
			name:     "invalid header",
			msg:      "not conventional",
			cfg:      cfg,
			problems: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Validate(tt.msg, tt.cfg)
			if len(got) != tt.problems {
				t.Errorf("Validate() returned %d problems %v, want %d", len(got), got, tt.problems)
			}
		})
	}
}
