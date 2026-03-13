package conv

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ATTron/jitzu/config"
)

var headerRegex = regexp.MustCompile(`^(\w+)(\(([^)]+)\))?(!)?: (.+)$`)

type Message struct {
	Type     string
	Scope    string
	Subject  string
	Body     string
	Breaking string
	Refs     string
}

func Build(typ, scope, subject, body, breaking, refs string) string {
	var sb strings.Builder

	if scope != "" {
		fmt.Fprintf(&sb, "%s(%s): %s", typ, scope, subject)
	} else {
		fmt.Fprintf(&sb, "%s: %s", typ, subject)
	}

	if body != "" {
		fmt.Fprintf(&sb, "\n\n%s", body)
	}

	if breaking != "" {
		fmt.Fprintf(&sb, "\n\nBREAKING CHANGE: %s", breaking)
	}

	if refs != "" {
		fmt.Fprintf(&sb, "\n\nRefs: %s", refs)
	}

	return sb.String()
}

func Parse(msg string) (Message, error) {
	lines := strings.SplitN(msg, "\n", 2)
	header := strings.TrimSpace(lines[0])

	matches := headerRegex.FindStringSubmatch(header)
	if matches == nil {
		return Message{}, fmt.Errorf("invalid conventional commit header: %q", header)
	}

	m := Message{
		Type:    matches[1],
		Scope:   matches[3],
		Subject: matches[5],
	}

	// Check for ! indicator
	if matches[4] == "!" && m.Breaking == "" {
		m.Breaking = m.Subject
	}

	if len(lines) > 1 {
		parseBody(lines[1], &m)
	}

	return m, nil
}

func parseBody(raw string, m *Message) {
	// Split into paragraphs by blank lines
	sections := strings.Split(strings.TrimSpace(raw), "\n\n")
	var bodyParts []string

	for _, section := range sections {
		section = strings.TrimSpace(section)
		if section == "" {
			continue
		}
		if strings.HasPrefix(section, "BREAKING CHANGE: ") {
			m.Breaking = strings.TrimPrefix(section, "BREAKING CHANGE: ")
		} else if strings.HasPrefix(section, "Refs: ") {
			m.Refs = strings.TrimPrefix(section, "Refs: ")
		} else {
			bodyParts = append(bodyParts, section)
		}
	}

	m.Body = strings.Join(bodyParts, "\n\n")
}

func Validate(msg string, cfg config.Config) []string {
	var problems []string

	m, err := Parse(msg)
	if err != nil {
		return []string{err.Error()}
	}

	// Validate type
	validType := false
	for _, t := range cfg.Types {
		if t.Name == m.Type {
			validType = true
			break
		}
	}
	if !validType {
		problems = append(problems, fmt.Sprintf("invalid type %q", m.Type))
	}

	// Validate subject length
	if cfg.SubjectMaxLen > 0 && len(m.Subject) > cfg.SubjectMaxLen {
		problems = append(problems, fmt.Sprintf("subject exceeds max length of %d characters", cfg.SubjectMaxLen))
	}

	// Validate scope required
	if cfg.ScopeRequired && m.Scope == "" {
		problems = append(problems, "scope is required")
	}

	// Validate scope is in allowed list
	if m.Scope != "" && len(cfg.Scopes) > 0 {
		validScope := false
		for _, s := range cfg.Scopes {
			if s == m.Scope {
				validScope = true
				break
			}
		}
		if !validScope {
			problems = append(problems, fmt.Sprintf("scope %q is not in the allowed list", m.Scope))
		}
	}

	// Validate body required
	if cfg.BodyRequired && m.Body == "" {
		problems = append(problems, "body is required")
	}

	// Validate body length
	if cfg.BodyMaxLen > 0 && len(m.Body) > cfg.BodyMaxLen {
		problems = append(problems, fmt.Sprintf("body exceeds max length of %d characters", cfg.BodyMaxLen))
	}

	return problems
}
