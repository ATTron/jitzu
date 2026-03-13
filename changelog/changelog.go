package changelog

import (
	"fmt"
	"strings"

	"github.com/ATTron/jitzu/conv"
	"github.com/ATTron/jitzu/jj"
)

var typeHeaders = map[string]string{
	"feat":     "Features",
	"fix":      "Bug Fixes",
	"docs":     "Documentation",
	"style":    "Styles",
	"refactor": "Code Refactoring",
	"perf":     "Performance Improvements",
	"test":     "Tests",
	"build":    "Build System",
	"ci":       "Continuous Integration",
	"chore":    "Chores",
	"revert":   "Reverts",
}

type Entry struct {
	ChangeID string
	CommitID string
	Type     string
	Scope    string
	Subject  string
}

const logTemplate = `change_id.shortest(8) ++ "\t" ++ commit_id.shortest(8) ++ "\t" ++ description.first_line() ++ "\n"`

func Generate(revset string) (string, error) {
	out, err := jj.Log(logTemplate, revset)
	if err != nil {
		return "", fmt.Errorf("reading jj log: %w", err)
	}

	entries := parseLog(out)
	return render(entries), nil
}

func parseLog(output string) []Entry {
	var entries []Entry
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) < 3 {
			continue
		}

		msg, err := conv.Parse(parts[2])
		if err != nil {
			continue // skip non-conventional commits
		}

		entries = append(entries, Entry{
			ChangeID: parts[0],
			CommitID: parts[1],
			Type:     msg.Type,
			Scope:    msg.Scope,
			Subject:  msg.Subject,
		})
	}
	return entries
}

func render(entries []Entry) string {
	grouped := make(map[string][]Entry)
	order := []string{"feat", "fix", "docs", "style", "refactor", "perf", "test", "build", "ci", "chore", "revert"}

	for _, e := range entries {
		grouped[e.Type] = append(grouped[e.Type], e)
	}

	var sb strings.Builder
	sb.WriteString("# Changelog\n\n")

	for _, typ := range order {
		entries, ok := grouped[typ]
		if !ok {
			continue
		}
		header := typeHeaders[typ]
		if header == "" {
			header = typ
		}
		fmt.Fprintf(&sb, "## %s\n\n", header)
		for _, e := range entries {
			if e.Scope != "" {
				fmt.Fprintf(&sb, "- **%s**: %s (%s)\n", e.Scope, e.Subject, e.CommitID)
			} else {
				fmt.Fprintf(&sb, "- %s (%s)\n", e.Subject, e.CommitID)
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
