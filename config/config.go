package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type CommitType struct {
	Name        string `toml:"name"`
	Description string `toml:"description"`
}

type Config struct {
	Types         []CommitType `toml:"types"`
	Scopes        []string     `toml:"scopes"`
	ScopeRequired bool         `toml:"scope_required"`
	BodyRequired  bool         `toml:"body_required"`
	SubjectMaxLen int          `toml:"subject_max_len"`
	BodyMaxLen    int          `toml:"body_max_len"`
}

func Default() Config {
	return Config{
		Types: []CommitType{
			{Name: "feat", Description: "A new feature"},
			{Name: "fix", Description: "A bug fix"},
			{Name: "docs", Description: "Documentation only changes"},
			{Name: "style", Description: "Changes that do not affect the meaning of the code"},
			{Name: "refactor", Description: "A code change that neither fixes a bug nor adds a feature"},
			{Name: "perf", Description: "A code change that improves performance"},
			{Name: "test", Description: "Adding missing tests or correcting existing tests"},
			{Name: "build", Description: "Changes that affect the build system or external dependencies"},
			{Name: "ci", Description: "Changes to CI configuration files and scripts"},
			{Name: "chore", Description: "Other changes that don't modify src or test files"},
			{Name: "revert", Description: "Reverts a previous commit"},
		},
		SubjectMaxLen: 72,
	}
}

func Load() Config {
	cfg := Default()

	path := findConfig()
	if path == "" {
		return cfg
	}

	var fileCfg Config
	if _, err := toml.DecodeFile(path, &fileCfg); err != nil {
		return cfg
	}

	return merge(cfg, fileCfg)
}

func findConfig() string {
	// Search from CWD upward
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		p := filepath.Join(dir, ".jitzu.toml")
		if _, err := os.Stat(p); err == nil {
			return p
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// Check XDG config
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	p := filepath.Join(home, ".config", "jitzu", "config.toml")
	if _, err := os.Stat(p); err == nil {
		return p
	}

	return ""
}

func merge(base, file Config) Config {
	if len(file.Types) > 0 {
		base.Types = file.Types
	}
	if len(file.Scopes) > 0 {
		base.Scopes = file.Scopes
	}
	if file.ScopeRequired {
		base.ScopeRequired = true
	}
	if file.BodyRequired {
		base.BodyRequired = true
	}
	if file.SubjectMaxLen > 0 {
		base.SubjectMaxLen = file.SubjectMaxLen
	}
	if file.BodyMaxLen > 0 {
		base.BodyMaxLen = file.BodyMaxLen
	}
	return base
}
