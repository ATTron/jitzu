package jj

import (
	"strings"
	"testing"
)

func TestDescribe(t *testing.T) {
	var captured []string
	Run = func(args ...string) (string, error) {
		captured = args
		return "", nil
	}

	t.Run("without revision", func(t *testing.T) {
		err := Describe("feat: test", "")
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"describe", "-m", "feat: test"}
		if strings.Join(captured, " ") != strings.Join(want, " ") {
			t.Errorf("got args %v, want %v", captured, want)
		}
	})

	t.Run("with revision", func(t *testing.T) {
		err := Describe("fix: bug", "@-")
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"describe", "-m", "fix: bug", "-r", "@-"}
		if strings.Join(captured, " ") != strings.Join(want, " ") {
			t.Errorf("got args %v, want %v", captured, want)
		}
	})
}

func TestCommit(t *testing.T) {
	var captured []string
	Run = func(args ...string) (string, error) {
		captured = args
		return "", nil
	}

	err := Commit("feat: new feature")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"commit", "-m", "feat: new feature"}
	if strings.Join(captured, " ") != strings.Join(want, " ") {
		t.Errorf("got args %v, want %v", captured, want)
	}
}

func TestLog(t *testing.T) {
	var captured []string
	Run = func(args ...string) (string, error) {
		captured = args
		return "output", nil
	}

	t.Run("without revset", func(t *testing.T) {
		out, err := Log("description", "")
		if err != nil {
			t.Fatal(err)
		}
		if out != "output" {
			t.Errorf("got %q, want %q", out, "output")
		}
		want := []string{"log", "--no-graph", "-T", "description"}
		if strings.Join(captured, " ") != strings.Join(want, " ") {
			t.Errorf("got args %v, want %v", captured, want)
		}
	})

	t.Run("with revset", func(t *testing.T) {
		_, err := Log("description", "main..@")
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"log", "--no-graph", "-T", "description", "-r", "main..@"}
		if strings.Join(captured, " ") != strings.Join(want, " ") {
			t.Errorf("got args %v, want %v", captured, want)
		}
	})
}
