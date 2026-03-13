package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if len(cfg.Types) != 11 {
		t.Errorf("expected 11 default types, got %d", len(cfg.Types))
	}
	if cfg.SubjectMaxLen != 72 {
		t.Errorf("expected default SubjectMaxLen 72, got %d", cfg.SubjectMaxLen)
	}
	if cfg.ScopeRequired {
		t.Error("expected ScopeRequired to be false by default")
	}
	if cfg.BodyRequired {
		t.Error("expected BodyRequired to be false by default")
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".jitzu.toml")

	content := `
scopes = ["api", "ui"]
scope_required = true
subject_max_len = 50

[[types]]
name = "feat"
description = "A new feature"

[[types]]
name = "fix"
description = "A bug fix"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Change to temp dir so config is found
	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	cfg := Load()

	if len(cfg.Types) != 2 {
		t.Errorf("expected 2 types from file, got %d", len(cfg.Types))
	}
	if len(cfg.Scopes) != 2 {
		t.Errorf("expected 2 scopes, got %d", len(cfg.Scopes))
	}
	if !cfg.ScopeRequired {
		t.Error("expected ScopeRequired to be true")
	}
	if cfg.SubjectMaxLen != 50 {
		t.Errorf("expected SubjectMaxLen 50, got %d", cfg.SubjectMaxLen)
	}
}

func TestMerge(t *testing.T) {
	base := Default()
	file := Config{
		Scopes:        []string{"core"},
		ScopeRequired: true,
		SubjectMaxLen: 50,
	}

	merged := merge(base, file)

	// Types should remain from base since file has none
	if len(merged.Types) != 11 {
		t.Errorf("expected 11 types from base, got %d", len(merged.Types))
	}
	// Scopes should come from file
	if len(merged.Scopes) != 1 || merged.Scopes[0] != "core" {
		t.Errorf("expected scopes [core], got %v", merged.Scopes)
	}
	if !merged.ScopeRequired {
		t.Error("expected ScopeRequired true")
	}
	if merged.SubjectMaxLen != 50 {
		t.Errorf("expected SubjectMaxLen 50, got %d", merged.SubjectMaxLen)
	}
}
