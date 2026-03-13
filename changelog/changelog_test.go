package changelog

import (
	"strings"
	"testing"
)

func TestParseLog(t *testing.T) {
	input := `abc12345	def67890	feat(ui): add dashboard
xyz11111	uvw22222	fix: null pointer
bad line without tabs
not-conventional	commit	hello world`

	entries := parseLog(input)
	if len(entries) != 2 {
		t.Fatalf("got %d entries, want 2", len(entries))
	}

	if entries[0].Type != "feat" || entries[0].Scope != "ui" || entries[0].Subject != "add dashboard" {
		t.Errorf("entry 0 = %+v", entries[0])
	}
	if entries[1].Type != "fix" || entries[1].Subject != "null pointer" {
		t.Errorf("entry 1 = %+v", entries[1])
	}
}

func TestRender(t *testing.T) {
	entries := []Entry{
		{CommitID: "abc1", Type: "feat", Scope: "ui", Subject: "add dashboard"},
		{CommitID: "abc2", Type: "feat", Subject: "add auth"},
		{CommitID: "def1", Type: "fix", Scope: "api", Subject: "null pointer"},
	}

	out := render(entries)

	if !strings.Contains(out, "## Features") {
		t.Error("missing Features header")
	}
	if !strings.Contains(out, "## Bug Fixes") {
		t.Error("missing Bug Fixes header")
	}
	if !strings.Contains(out, "- **ui**: add dashboard (abc1)") {
		t.Error("missing scoped entry")
	}
	if !strings.Contains(out, "- add auth (abc2)") {
		t.Error("missing unscoped entry")
	}

	// Features should come before Bug Fixes
	featIdx := strings.Index(out, "## Features")
	fixIdx := strings.Index(out, "## Bug Fixes")
	if featIdx > fixIdx {
		t.Error("Features should come before Bug Fixes")
	}
}
