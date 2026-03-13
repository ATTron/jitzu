package jj

import (
	"fmt"
	"os/exec"
	"strings"
)

// Run executes jj with the given arguments and returns stdout.
// It is a variable so tests can replace it.
var Run = func(args ...string) (string, error) {
	cmd := exec.Command("jj", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("jj %s: %w\n%s", strings.Join(args, " "), err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

func Describe(message, revision string) error {
	args := []string{"describe", "-m", message}
	if revision != "" {
		args = append(args, "-r", revision)
	}
	_, err := Run(args...)
	return err
}

func Commit(message string) error {
	_, err := Run("commit", "-m", message)
	return err
}

func BookmarkCreate(name string) error {
	_, err := Run("bookmark", "create", name)
	return err
}

func BookmarkSet(name string, revision string) error {
	args := []string{"bookmark", "set", name}
	if revision != "" {
		args = append(args, "-r", revision)
	}
	_, err := Run(args...)
	return err
}

func BookmarkAdvance(name string) error {
	_, err := Run("bookmark", "advance", name)
	return err
}

func Log(template, revset string) (string, error) {
	args := []string{"log", "--no-graph", "-T", template}
	if revset != "" {
		args = append(args, "-r", revset)
	}
	return Run(args...)
}
