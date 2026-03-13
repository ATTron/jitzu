package prompt

import (
	"errors"
	"fmt"

	"github.com/ATTron/jitzu/config"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

// ErrAborted is returned when the user quits the prompt.
var ErrAborted = errors.New("aborted")

type Result struct {
	Type     string
	Scope    string
	Subject  string
	Body     string
	Breaking string
	Refs     string
}

func Run(cfg config.Config) (Result, error) {
	var r Result
	var isBreaking bool

	typeOptions := make([]huh.Option[string], len(cfg.Types))
	for i, t := range cfg.Types {
		typeOptions[i] = huh.NewOption(fmt.Sprintf("%-10s %s", t.Name, t.Description), t.Name)
	}

	scopeField := huh.NewInput().
		Title("Scope").
		Description("Optional scope of the change").
		Value(&r.Scope)

	if len(cfg.Scopes) > 0 {
		scopeField.SuggestionsFunc(func() []string {
			return cfg.Scopes
		}, &r.Scope)
	}

	subjectValidate := func(s string) error {
		if s == "" {
			return fmt.Errorf("subject is required")
		}
		if cfg.SubjectMaxLen > 0 && len(s) > cfg.SubjectMaxLen {
			return fmt.Errorf("subject must be %d characters or less", cfg.SubjectMaxLen)
		}
		return nil
	}

	bodyValidate := func(s string) error {
		if cfg.BodyRequired && s == "" {
			return fmt.Errorf("body is required")
		}
		if cfg.BodyMaxLen > 0 && len(s) > cfg.BodyMaxLen {
			return fmt.Errorf("body must be %d characters or less", cfg.BodyMaxLen)
		}
		return nil
	}

	form := huh.NewForm(
		// Page 1: Type + Scope
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Commit Type").
				Options(typeOptions...).
				Value(&r.Type),
			scopeField,
		),

		// Page 2: Subject + Body
		huh.NewGroup(
			huh.NewInput().
				Title("Subject").
				Description("Short description of the change").
				Validate(subjectValidate).
				Value(&r.Subject),
			huh.NewText().
				Title("Body").
				Description("Longer description (optional)").
				Validate(bodyValidate).
				Value(&r.Body),
		),

		// Page 3: Breaking + Refs
		huh.NewGroup(
			huh.NewConfirm().
				Title("Breaking change?").
				Value(&isBreaking),
			huh.NewInput().
				Title("Breaking change description").
				Value(&r.Breaking),
			huh.NewInput().
				Title("Issue references").
				Description("e.g. #123, PROJ-456").
				Value(&r.Refs),
		),
	)

	err := form.WithKeyMap(newKeyMap()).Run()
	if err != nil {
		return Result{}, handleAbort(err)
	}

	if !isBreaking {
		r.Breaking = ""
	}

	return r, nil
}

func newKeyMap() *huh.KeyMap {
	km := huh.NewDefaultKeyMap()
	km.Quit = key.NewBinding(key.WithKeys("ctrl+c", "esc"), key.WithHelp("esc", "quit"))
	return km
}

func handleAbort(err error) error {
	if errors.Is(err, huh.ErrUserAborted) {
		return ErrAborted
	}
	return err
}

// SelectAction presents a list of choices and returns the selected one.
func SelectAction(title string, options []string) (string, error) {
	var choice string

	opts := make([]huh.Option[string], len(options))
	for i, o := range options {
		opts[i] = huh.NewOption(o, o)
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(title).
				Options(opts...).
				Value(&choice),
		),
	).WithKeyMap(newKeyMap()).Run()
	if err != nil {
		return "", handleAbort(err)
	}

	return choice, nil
}

var bookmarkTypes = []huh.Option[string]{
	huh.NewOption("feat       A new feature", "feat"),
	huh.NewOption("fix        A bug fix", "fix"),
	huh.NewOption("hotfix     An urgent production fix", "hotfix"),
	huh.NewOption("refactor   Code restructuring", "refactor"),
	huh.NewOption("chore      Maintenance / tooling", "chore"),
	huh.NewOption("docs       Documentation", "docs"),
	huh.NewOption("test       Testing", "test"),
	huh.NewOption("experiment Spike / proof of concept", "experiment"),
	huh.NewOption("release    Release preparation", "release"),
}

// BookmarkName prompts the user to build a structured bookmark name (type/short-name).
func BookmarkName(cfg config.Config) (string, error) {
	var typ, name, ticket string

	err := huh.NewForm(
		// Page 1: Category
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Bookmark category").
				Description("What kind of work is this?").
				Options(bookmarkTypes...).
				Value(&typ),
		),

		// Page 2: Name + ticket
		huh.NewGroup(
			huh.NewInput().
				Title("Short name").
				Description("Kebab-case, e.g. auth-login, parse-revsets").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name is required")
					}
					for _, c := range s {
						if c == ' ' || c == '/' || c == '_' {
							return fmt.Errorf("use kebab-case (hyphens, no spaces/slashes/underscores)")
						}
					}
					return nil
				}).
				Value(&name),
			huh.NewInput().
				Title("Ticket / issue (optional)").
				Description("e.g. PROJ-123 — will be prepended to the name").
				Value(&ticket),
		),
	).WithKeyMap(newKeyMap()).Run()
	if err != nil {
		return "", handleAbort(err)
	}

	bookmark := typ + "/" + name
	if ticket != "" {
		bookmark = typ + "/" + ticket + "-" + name
	}

	return bookmark, nil
}
