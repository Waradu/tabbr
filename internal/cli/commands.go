package cli

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/waradu/tabbr/internal/matcher"
	"github.com/waradu/tabbr/internal/shell"
	"github.com/waradu/tabbr/internal/store"
)

func (cmd *AddCommand) Run(app *App) error {
	return store.Add(app.DB, cmd.Command)
}

func (cmd *RemoveCommand) Run(app *App) error {
	return store.Remove(app.DB, cmd.Command)
}

func (cmd *ListCommand) Run(app *App) error {
	commands, err := store.List(app.DB)
	if err != nil {
		return err
	}

	for _, ranked := range matcher.Rank(commands) {
		command := ranked.Command

		fmt.Printf(
			"%s %s (Score: %g)\n",
			dim(time.Unix(command.LastUsed, 0).Format("2006-01-02 15:04")),
			highlight(command.Text),
			ranked.Score,
		)
	}

	return nil
}

func (cmd *PathCommand) Run() error {
	path, err := store.Path()
	if err != nil {
		return err
	}

	fmt.Println(path)
	return nil
}

func (cmd *PrepareCommand) Run() error {
	return nil
}

func (cmd *QueryCommand) Run(app *App) error {
	if strings.IndexFunc(cmd.Text, unicode.IsSpace) >= 0 {
		return fmt.Errorf("query must not contain spaces")
	}
	if looksLikePath(cmd.Text) {
		return fmt.Errorf("query must not be a path")
	}

	commands, err := store.List(app.DB)
	if err != nil {
		return err
	}

	results := matcher.Query(commands, cmd.Text)
	results = results[:min(len(results), 50)]

	for _, ranked := range results {
		fmt.Println(ranked.Command.Text)
	}

	return nil
}

func looksLikePath(text string) bool {
	if text == "." || text == ".." || strings.HasPrefix(text, "~") || strings.ContainsAny(text, `/\`) {
		return true
	}

	return len(text) >= 2 && text[1] == ':' &&
		(text[0] >= 'a' && text[0] <= 'z' || text[0] >= 'A' && text[0] <= 'Z')
}

func (cmd *ExcludeAddCommand) Run(app *App) error {
	return store.AddExclusion(app.DB, cmd.Pattern)
}

func (cmd *ExcludeRemoveCommand) Run(app *App) error {
	return store.RemoveExclusion(app.DB, cmd.Pattern)
}

func (cmd *ExcludeListCommand) Run(app *App) error {
	exclusions, err := store.ListExclusions(app.DB)
	if err != nil {
		return err
	}

	for _, exclusion := range exclusions {
		fmt.Printf(
			"%s %s\n",
			dim(time.Unix(exclusion.CreatedAt, 0).Format("2006-01-02 15:04")),
			highlight(exclusion.Pattern),
		)
	}

	return nil
}

func (cmd *InitCommand) Run() error {
	switch cmd.Shell {
	case "zsh":
		fmt.Print(shell.Zsh)
	case "bash":
		fmt.Print(shell.Bash)
	case "pwsh":
		fmt.Print(shell.PowerShell)
	}
	return nil
}
