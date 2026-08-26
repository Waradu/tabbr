package cli

import (
	"database/sql"

	"github.com/waradu/tabbr/internal/store"

	"github.com/alecthomas/kong"
)

type AddCommand struct {
	Command string `arg:"" help:"The command to add"`
}

type RemoveCommand struct {
	Command string `arg:"" help:"The command to remove"`
}

type ListCommand struct{}

type PathCommand struct{}

type PrepareCommand struct{}

type QueryCommand struct {
	Text string `arg:"" help:"The text to match"`
}

type InitCommand struct {
	Shell string `arg:"" enum:"zsh,bash,pwsh" help:"The shell to initialize"`
}

type ExcludeAddCommand struct {
	Pattern string `arg:"" help:"The pattern to exclude"`
}

type ExcludeRemoveCommand struct {
	Pattern string `arg:"" help:"The exclusion pattern to remove"`
}

type ExcludeListCommand struct{}

type ExcludeCommand struct {
	Add    ExcludeAddCommand    `cmd:"" help:"Add an exclusion pattern"`
	Remove ExcludeRemoveCommand `cmd:"" help:"Remove an exclusion pattern"`
	List   ExcludeListCommand   `cmd:"" help:"List exclusion patterns"`
}

type CLI struct {
	Init    InitCommand    `cmd:"" db:"none" help:"Print shell integration"`
	Add     AddCommand     `cmd:"" help:"Add a new command to the store"`
	Remove  RemoveCommand  `cmd:"" help:"Remove a command from the store"`
	List    ListCommand    `cmd:"" help:"List all commands in the store"`
	Path    PathCommand    `cmd:"" db:"none" help:"Print the database path"`
	Prepare PrepareCommand `cmd:"" help:"Prepare the database"`
	Query   QueryCommand   `cmd:"" help:"Find matching commands"`
	Exclude ExcludeCommand `cmd:"" help:"Manage exclusions"`
}

type App struct {
	DB *sql.DB
}

func Init() error {
	var cli CLI
	ctx := kong.Parse(&cli)

	selected := ctx.Selected()
	if selected != nil && selected.Tag.Get("db") == "none" {
		return ctx.Run()
	}

	db, err := store.Open()
	if err != nil {
		return err
	}

	defer db.Close()

	return ctx.Run(&App{DB: db})
}
